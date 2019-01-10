package redis

import "context"

type Subscription struct {
	ctx      context.Context
	cancel   context.CancelFunc // unsubscribe
	channels map[string]chan Value
}

func newSubscription(ctx context.Context, length int) *Subscription {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	return &Subscription{
		channels: make(map[string]chan Value, length),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (s *Subscription) Handle(handler func(string, Value)) {
	for channel, c := range s.channels {
		go func(channel string, c <-chan Value) {
			for {
				select {
				case <-s.ctx.Done():
					return
				case v := <-c:
					handler(channel, v)
				}
			}
		}(channel, c)
	}
}

func (s *Subscription) Unsubscribe() {
	s.cancel()
	for _, c := range s.channels {
		close(c)
	}
}
