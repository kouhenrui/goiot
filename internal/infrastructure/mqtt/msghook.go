package mqtt

import (
	"context"

	"github.com/gogf/gf/v2/os/glog"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

type MessageHook struct {
	mqtt.HookBase
	ctx context.Context
}

func NewMessageHook(ctx context.Context) *MessageHook {
	return &MessageHook{ctx: ctx}
}

// 声明监听Publish事件
func (h *MessageHook) Provides(
	event byte,
) bool {

	return event == mqtt.OnPublish
}
func (h *MessageHook) ID() string {
	return "message-hook"
}

// MQTT设备发布消息时触发
func (h *MessageHook) OnPublish(
	client *mqtt.Client,
	packet packets.Packet,
) (packets.Packet, error) {
	deviceKey := client.ID

	topic := packet.TopicName

	payload := packet.Payload

	glog.Info(
		h.ctx,
		"mqtt message received",
		"device",
		deviceKey,
		"topic",
		topic,
		"payload",
		string(payload),
	)

	// 其他Topic以后扩展
	// switch topic {
	// }
	return packet, nil
}
