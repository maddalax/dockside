package app

import (
	"context"
	"dockman/app/subject"
	"github.com/nats-io/nats.go"
	"io"
)

type WriterSubscriber struct {
	Writer     io.WriteCloser
	Subscriber chan *nats.Msg
}

type NatsWriterCreateOptions struct {
	BeforeWrite func(data string) bool
}

func (c *KvClient) CreateEphemeralWriterSubscriber(ctx context.Context, subject subject.Subject, opts NatsWriterCreateOptions) *WriterSubscriber {
	ch := make(chan *nats.Msg, 100)
	_, err := c.SubscribeSubject(ctx, subject, func(msg *nats.Msg) {
		ch <- msg
	})
	if err != nil {
		return &WriterSubscriber{
			Writer:     &EmptyWriter{},
			Subscriber: ch,
		}
	}

	w := c.NewEphemeralNatsWriter(subject)
	if opts.BeforeWrite != nil {
		w.SetBeforeWrite(opts.BeforeWrite)
	}

	return &WriterSubscriber{
		Writer:     w,
		Subscriber: ch,
	}
}
