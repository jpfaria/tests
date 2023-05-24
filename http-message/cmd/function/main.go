package main

import (
	"context"
	"github.com/americanas-go/config"
	ilog "github.com/americanas-go/ignite/americanas-go/log.v1"
	ignats "github.com/americanas-go/ignite/nats-io/nats.go.v1"
	"github.com/americanas-go/log"
	"github.com/jpfaria/tests/http-message/internal/pkg/settings"
	"github.com/nats-io/nats.go"
	"strings"
)

func main() {

	config.Load()
	ilog.New()

	conn, err := ignats.NewConn(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	subject := settings.Subject

	ch := make(chan *nats.Msg)

	conn.QueueSubscribeSyncWithChan(subject, "function", ch)

	for msg := range ch {

		msg := msg

		go func() {
			id := msg.Header.Get(settings.UUIDHeader)
			instance := msg.Header.Get(settings.InstanceHeader)

			log.Infof("received from subject [%s], ID [%s], Instance [%s]", subject, id, instance)

			respSubject := strings.Join([]string{instance, subject, id}, "-")

			msgResp := &nats.Msg{
				Subject: respSubject,
				// Data:    b,
			}

			// time.Sleep(10 * time.Millisecond)

			if err := conn.PublishMsg(msgResp); err != nil {
				log.Fatal(err)
			}

			log.Infof("published group message on subject [%s]", respSubject)
		}()

	}

}
