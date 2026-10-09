package mqtt

import (
	"bytes"
	"context"
	"strings"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	app "goiot/internal/app/device"
)

type AuthHook struct {
	mqtt.HookBase
	auth *app.AuthService
}

func NewAuthHook(auth *app.AuthService) *AuthHook {
	return &AuthHook{auth: auth}
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

	ctx, cancel := context.WithTimeout(
		context.Background(),
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
