package natsconnection

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// NC is a connection to the nats Server
var NC *nats.Conn

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

// Connect connects to nats server and returns instance
// On shutdown, you'll need to call natsconnection.NC.Drain()
func Connect(url string, username string, password string) {
	opts := []nats.Option{}
	hostname, err := os.Hostname()
	if err == nil {
		opts = append(opts, nats.Name(hostname))
	}

	opts = append(opts, nats.UserInfo(username, password))
	opts = setupConnOptions(opts)

	NC, err = nats.Connect(url, opts...)
	if err != nil {
		log.Fatal(err)
	}
}

// SubscribeToQueue subscribes on a channel
// subj is the channel to subscribe on
// queueName, if specified joins the queue so each message is only received by one running container
// all containers subscribed to a specified subj receive the message unless in a queueName, then
// only one in each queueName receives the message
// nats.Msg has the subj in it, so we can do routing from there
func SubscribeToQueue(subjBase string, queueName string, handlerRouter func(*nats.Msg)) {
	NC.QueueSubscribe(subjBase+".>", queueName, handlerRouter)
	NC.Flush()
}

// Request sends a message and expects a response back
func Request(ctx context.Context, subj string, data []byte, timeOut time.Duration) (*nats.Msg, error) {
	if ctx != nil {
		return NC.RequestWithContext(ctx, subj, data)
	}
	return NC.Request(subj, data, timeOut)
}
