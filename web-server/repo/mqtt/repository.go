package mqtt

import (
	"context"
	"fmt"

	mochi "github.com/mochi-mqtt/server/v2"
)

type response[T any] struct {
	val T
	err error
}

type pushConfigReq struct {
	ctx   context.Context
	input PushConfigInput
	resp  chan response[struct{}]
}

type publishReq struct {
	ctx   context.Context
	input PublishInput
	resp  chan response[struct{}]
}

type shutdownReq struct {
	resp chan struct{}
}

type mqttRepo struct {
	server      *mochi.Server
	pushCfgCh   chan pushConfigReq
	publishCh   chan publishReq
	shutdownCh  chan shutdownReq
}

// NewRepository creates an MQTTRepo wrapping the embedded Mochi server's inline client.
// The Mochi server must have been created with InlineClient: true.
func NewRepository(server *mochi.Server) Repository {
	r := &mqttRepo{
		server:     server,
		pushCfgCh:  make(chan pushConfigReq),
		publishCh:  make(chan publishReq),
		shutdownCh: make(chan shutdownReq),
	}
	go r.manage()
	return r
}

func (r *mqttRepo) manage() {
	for {
		select {
		case req := <-r.pushCfgCh:
			err := r.doPushConfig(req.input)
			req.resp <- response[struct{}]{err: err}
		case req := <-r.publishCh:
			err := r.doPublish(req.input)
			req.resp <- response[struct{}]{err: err}
		case req := <-r.shutdownCh:
			close(req.resp)
			return
		}
	}
}

func (r *mqttRepo) doPushConfig(input PushConfigInput) error {
	topic := fmt.Sprintf("rootstock/%s/config", input.DeviceID)
	return r.server.Publish(topic, input.Payload, true, 1)
}

func (r *mqttRepo) doPublish(input PublishInput) error {
	return r.server.Publish(input.Topic, input.Payload, false, input.QoS)
}

func (r *mqttRepo) PushDeviceConfig(_ context.Context, input PushConfigInput) error {
	resp := make(chan response[struct{}], 1)
	r.pushCfgCh <- pushConfigReq{input: input, resp: resp}
	res := <-resp
	return res.err
}

func (r *mqttRepo) PublishToDevice(_ context.Context, input PublishInput) error {
	resp := make(chan response[struct{}], 1)
	r.publishCh <- publishReq{input: input, resp: resp}
	res := <-resp
	return res.err
}

func (r *mqttRepo) Shutdown() {
	resp := make(chan struct{}, 1)
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}
