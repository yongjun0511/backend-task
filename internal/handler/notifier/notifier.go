package notifier

import (
	"fmt"
	"sync"

	"banksalad-backend-task/internal/domain"
	"banksalad-backend-task/internal/handler/notifier/channelhandler"
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
		wg.Add(1)
		go func(b bucket) {
			defer wg.Done()
			if err := b.handler.SendBatch(b.values); err != nil {
				once.Do(func() { first = err })
			}
		}(b)
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
