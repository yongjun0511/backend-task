package notifier

import (
	"fmt"
	"log"
	"sync"
	"time"

	"banksalad-backend-task/internal/domain"
	"banksalad-backend-task/internal/handler/notifier/channelhandler"
)

const (
	startBackoff = 50 * time.Millisecond
	maxBackoff   = time.Second
)

type Notifier struct {
	handlers map[domain.FieldType]channelhandler.ChannelHandler
}

func NewNotifier(list []channelhandler.ChannelHandler) *Notifier {
	m := make(map[domain.FieldType]channelhandler.ChannelHandler, len(list))
	for _, h := range list {
		m[h.TargetField()] = h
	}
	return &Notifier{handlers: m}
}
func (n *Notifier) NotifyAll(data map[domain.ChannelDTO]map[string]struct{}) error {
	var (
		wg    sync.WaitGroup
		once  sync.Once
		first error
		bkts  = groupByFieldType(n.handlers, data)
	)

	if first = bkts.err; first != nil {
		return first
	}

	for _, b := range bkts.good {
		for _, v := range b.values {
			wg.Add(1)
			go func(h channelhandler.ChannelHandler, val string) {
				defer wg.Done()
				if err := sendUntilSuccess(h, val); err != nil {
					once.Do(func() { first = err })
				}
			}(b.handler, v)
		}
	}

	wg.Wait()
	return first
}

type bucket struct {
	handler channelhandler.ChannelHandler
	values  []string
}

type bucketResult struct {
	good []bucket
	err  error
}

func groupByFieldType(
	all map[domain.FieldType]channelhandler.ChannelHandler,
	data map[domain.ChannelDTO]map[string]struct{},
) bucketResult {

	res := bucketResult{}
	bk := map[domain.FieldType]*bucket{}

	for dto, set := range data {
		h, ok := all[dto.FieldType]
		if !ok {
			res.err = fmt.Errorf("unsupported field type: %s", dto.FieldType)
			return res
		}
		b, ok := bk[dto.FieldType]
		if !ok {
			b = &bucket{handler: h}
			bk[dto.FieldType] = b
		}
		for v := range set {
			b.values = append(b.values, v)
		}
	}
	for _, b := range bk {
		res.good = append(res.good, *b)
	}
	return res
}

func sendUntilSuccess(h channelhandler.ChannelHandler, v string) error {
	backoff := startBackoff
	for {
		if err := h.Send(v); err == nil {
			return nil
		}
		log.Printf("[WARN] send failed, retrying in %v for %s", backoff, v)
		time.Sleep(backoff)
		if backoff < maxBackoff {
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
}
