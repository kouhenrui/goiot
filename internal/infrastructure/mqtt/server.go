package mqtt

import (
	"context"
	"fmt"
	app "goiot/internal/app/device"
	"goiot/internal/config"

	"github.com/gogf/gf/v2/os/glog"
	mqtt "github.com/mochi-mqtt/server/v2"
	mqauth "github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/packets"

	"github.com/mochi-mqtt/server/v2/listeners"
)

type Server struct {
	ctx    context.Context
	broker *mqtt.Server
}

func NewServer(ctx context.Context, auth *app.AuthService, cfg config.MQTTConfig) (*Server, error) {
	s := mqtt.New(
		&mqtt.Options{
			InlineClient: true,
		},
		// &mqtt.Options{
		// InlineClient: true,
		// Capabilities: mqtt.Capabilities{
		// 	Auth: true,
		// },
		// }
	)

	// if err := s.AddHook(NewAuthHook(ctx, auth), nil); err != nil {
	// 	return nil, fmt.Errorf("add mqtt auth hook: %w", err)
	// }
	// 允许所有连接(权限)。
	_ = s.AddHook(new(mqauth.AllowHook), nil)
	_ = s.AddHook(NewMessageHook(ctx), nil)
	tcp := listeners.NewTCP(listeners.Config{
		ID:      "goiot-mqtt",
		Address: cfg.Address,
	})

	if err := s.AddListener(tcp); err != nil {
		return nil, fmt.Errorf("add mqtt listener: %w", err)
	}

	return &Server{ctx: ctx, broker: s}, nil
}

// Start 启动 MQTT 服务器。
func (s *Server) Start() error {
	err := s.broker.Serve()
	if err != nil {
		return fmt.Errorf("start mqtt server: %w", err)
	}
	return s.Publish("devices/dev_001", []byte("你好,我已连接到mqtt,正在发送测试消息"), true, 0)
}

// Close 关闭 MQTT 服务器。
func (s *Server) Close() error {
	return s.broker.Close()
}

// Publish 发布消息到指定主题。
func (s *Server) Publish(
	topic string,
	payload []byte,
	retained bool,
	qos byte,
) error {
	return s.broker.Publish(topic, payload, retained, qos)
}

func (s *Server) OnPublish(
	cl *mqtt.Client,
	pk packets.Packet,
) {
	glog.Info(s.ctx, "mqtt message received",
		"client_id", cl.ID,
		"topic", pk.TopicName)

}
