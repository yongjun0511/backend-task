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

func NewNotifier(h map[domain.FieldType]channelhandler.ChannelHandler) *Notifier {
	return &Notifier{handlers: h}
}

func (n *Notifier) NotifyAll(data map[domain.FieldType]map[string]struct{}) error {
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
	data map[domain.FieldType]map[string]struct{},
) bucketResult {
	res := bucketResult{}

	for ft, set := range data {
		h, ok := all[ft]
		if !ok {
			res.err = fmt.Errorf("unsupported field type: %s", ft)
			return res
		}

		b := &bucket{handler: h}
		for v := range set {
			b.values = append(b.values, v)
		}
		res.good = append(res.good, *b)
	}
	return res
}
