package serial

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/goburrow/modbus"
)

type Driver struct {
	mu      sync.Mutex
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

func (d *Driver) Connect(
	ctx context.Context,
	cfg Config,
) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.handler != nil {
		return fmt.Errorf("driver already connected")
	}

	h := modbus.NewRTUClientHandler(cfg.Port)
	h.BaudRate = cfg.BaudRate
	h.DataBits = cfg.DataBits
	h.StopBits = cfg.StopBits
	h.SlaveId = cfg.SlaveID
	h.Timeout = 3 * time.Second

	switch cfg.Parity {
	case "N":
		h.Parity = "N"
	case "E":
		h.Parity = "E"
	case "O":
		h.Parity = "O"
	default:
		return fmt.Errorf("unsupported parity: %s", cfg.Parity)
	}

	if err := h.Connect(); err != nil {
		return fmt.Errorf("open modbus serial port: %w", err)
	}

	d.handler = h
	d.client = modbus.NewClient(h)
	return nil
}

func (d *Driver) Read(
	ctx context.Context,
	point string,
) (any, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.client == nil {
		return nil, fmt.Errorf("device is not connected")
	}

	// 示例：从保持寄存器地址 0 读取一个寄存器。
	// 实际项目中应通过 point 查找点位配置。
	if point != "temperature" {
		return nil, fmt.Errorf("unknown point: %s", point)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := d.client.ReadHoldingRegisters(0, 1)
	if err != nil {
		return nil, fmt.Errorf("read holding registers: %w", err)
	}

	if len(data) < 2 {
		return nil, fmt.Errorf("invalid register response")
	}

	raw := uint16(data[0])<<8 | uint16(data[1])

	// 示例设备约定：温度以 0.1 摄氏度为单位。
	// 必须根据真实设备手册确定比例和数据类型。
	return float64(int16(raw)) / 10, nil
}

func (d *Driver) Write(
	ctx context.Context,
	point string,
	value any,
) error {
	return fmt.Errorf("write operation not implemented")
}

func (d *Driver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.handler == nil {
		return nil
	}

	err := d.handler.Close()
	d.handler = nil
	d.client = nil
	return err
}
