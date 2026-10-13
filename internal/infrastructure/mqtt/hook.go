package mqtt

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/glog"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	app "goiot/internal/app/device"
)

type MqTopicType string

const (
	MqTopicTelemetry MqTopicType = "telemetry"
	MqTopicCommands  MqTopicType = "commands"
	MqTopicState     MqTopicType = "state"
	MqTopicConfig    MqTopicType = "config"
	MqTopicEvents    MqTopicType = "events"
)

type AuthHook struct {
	ctx context.Context
	mqtt.HookBase
	auth *app.AuthService
}

func NewAuthHook(ctx context.Context, auth *app.AuthService) *AuthHook {
	return &AuthHook{ctx: ctx, auth: auth}
}

func (h *AuthHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnectAuthenticate,
		mqtt.OnACLCheck,
	}, []byte{b})
}

// 连接认证：用户名和 ClientID 都必须等于 DeviceKey。
func (h *AuthHook) OnConnectAuthenticate(
	cl *mqtt.Client,
	pk packets.Packet,
) bool {
	username := string(pk.Connect.Username)
	password := string(pk.Connect.Password)
	glog.Info(h.ctx, "mqtt auth attempt",
		"client_id", cl.ID,
		"username", username,
		"password_present", password != "",
	)
	ctx, cancel := context.WithTimeout(
		h.ctx,
		3*time.Second,
	)
	defer cancel()
	return h.auth.AuthenticateMQTT(ctx, username, cl.ID, password)

}

// Topic 授权：设备只能发布遥测、订阅命令。
func (h *AuthHook) OnACLCheck(
	cl *mqtt.Client,
	topic string,
	write bool,
) bool {
	deviceKey := cl.ID
	if deviceKey == "" || strings.Contains(deviceKey, "/") {
		return false
	}

	telemetryTopic := "devices/" + deviceKey + "/telemetry"
	commandTopic := "devices/" + deviceKey + "/commands"

	if write {
		// 设备发布消息：只允许自己的遥测 Topic。
		return topic == telemetryTopic
	}

	// 设备订阅消息：只允许自己的命令 Topic。
	return topic == commandTopic
}

/*
- topic 规则：
- 数据上报	devices/{id}/telemetry
- 命令下发	devices/{id}/commands
- 在线状态	devices/{id}/state
- 配置同步	devices/{id}/config
- 报警事件	devices/{id}/events
*/
func GetTopic(deviceKey string, topicType string) string {
	switch topicType {
	case "telemetry":
		return "devices/" + deviceKey + "/telemetry"
	case "commands":
		return "devices/" + deviceKey + "/commands"
	case "state":
		return "devices/" + deviceKey + "/state"
	case "config":
		return "devices/" + deviceKey + "/config"
	case "events":
		return "devices/" + deviceKey + "/events"
	default:
		return ""
	}
}
