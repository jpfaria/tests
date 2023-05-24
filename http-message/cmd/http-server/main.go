package main

import (
	"context"
	"github.com/americanas-go/config"
	ilog "github.com/americanas-go/ignite/americanas-go/log.v1"
	"github.com/americanas-go/ignite/labstack/echo.v4"
	ignats "github.com/americanas-go/ignite/nats-io/nats.go.v1"
	"github.com/americanas-go/log"
	"github.com/google/uuid"
	"github.com/jpfaria/tests/http-message/internal/pkg/settings"
	e "github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"net/http"
	"os"
	"strings"
	"time"
)

type Handler struct {
	conn     *nats.Conn
	subject  string
	instance string
}

func NewHandler(conn *nats.Conn, instance string, subject string) *Handler {
	return &Handler{conn: conn, instance: instance, subject: subject}
}

func (h *Handler) Get(c e.Context) (err error) {

	id := uuid.New()

	log.Infof("received HTTP request with ID [%s]", id.String())

	respSubject := strings.Join([]string{h.instance, h.subject, id.String()}, "-")

	ch := make(chan *nats.Msg)
	subs, err := h.conn.ChanQueueSubscribe(respSubject, "http-server", ch)
	if err != nil {
		return err
	}
	if err := subs.AutoUnsubscribe(1); err != nil {
		return err
	}

	headers := nats.Header{}
	headers.Set(settings.UUIDHeader, id.String())
	headers.Set(settings.InstanceHeader, h.instance)

	msg := &nats.Msg{
		Header:  headers,
		Subject: h.subject,
		// Data:    b,
	}

	if err := h.conn.PublishMsg(msg); err != nil {
		log.Fatal(err)
	}

	log.Infof("published group message on subject [%s]", h.subject)

	t := time.Now()

	select {
	case msg := <-ch:
		log.Infof("returned on subject %s in %v", msg.Subject, time.Since(t).Milliseconds())
		close(ch)
		return c.JSON(http.StatusOK, respSubject)
	case <-time.After(20 * time.Millisecond):
		log.Infof("timeout on subject %s in %v", msg.Subject, time.Since(t).Milliseconds())
		return c.JSON(http.StatusRequestTimeout, respSubject)
	}

}

func main() {

	config.Load()
	ilog.New()

	conn, err := ignats.NewConn(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	instance := os.Getenv("INSTANCE")

	handler := NewHandler(conn, instance, settings.Subject)

	ctx := context.Background()

	srv := echo.NewServer(ctx)

	srv.GET("/test", handler.Get)

	srv.Serve(ctx)
}
