package mqtt

import (
	"fmt"
	app "goiot/internal/app/device"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
)

type Server struct {
	broker *mqtt.Server
}

func NewServer(auth *app.AuthService, address string) (*Server, error) {
	s := mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	if err := s.AddHook(NewAuthHook(auth), nil); err != nil {
		return nil, fmt.Errorf("add mqtt auth hook: %w", err)
	}

	tcp := listeners.NewTCP(listeners.Config{
		ID:      "mqtt-tcp",
		Address: address,
	})

	if err := s.AddListener(tcp); err != nil {
		return nil, fmt.Errorf("add mqtt listener: %w", err)
	}

	return &Server{broker: s}, nil
}

func (s *Server) Start() error {
	return s.broker.Serve()
}

func (s *Server) Close() error {
	return s.broker.Close()
}

func (s *Server) Publish(
	topic string,
	payload []byte,
	retained bool,
	qos byte,
) error {
	return s.broker.Publish(topic, payload, retained, qos)
}
