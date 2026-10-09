package serial

import "context"

// Driver 串口驱动接口
type DriverI interface {
	Connect(ctx context.Context, cfg Config) error
	Read(ctx context.Context, point string) (any, error)
	Write(ctx context.Context, point string, value any) error
	Close() error
}

// Config 串口配置
type Config struct {
	Port     string `json:"port"`            // 串口名称，例如 COM1、/dev/ttyUSB0
	BaudRate int    `json:"baud_rate"`       // 波特率，例如 9600、115200
	DataBits int    `json:"data_bits"`       // 数据位，例如 7、8
	StopBits int    `json:"stop_bits"`       // 停止位，例如 1、2
	Parity   string `json:"parity"`          // 校验位，例如 N、E、O
	SlaveID  byte   `json:"slave_id"`        // 从站ID，例如 1、2、3
	Timeout  int    `json:"timeout_seconds"` // 超时时间，单位秒
}

type enumPort string

const (
	PortCOM1  enumPort = "COM1"
	PortCOM2  enumPort = "COM2"
	PortCOM3  enumPort = "COM3"
	PortCOM4  enumPort = "COM4"
	PortCOM5  enumPort = "COM5"
	PortCOM6  enumPort = "COM6"
	PortCOM7  enumPort = "COM7"
	PortCOM8  enumPort = "COM8"
	PortCOM9  enumPort = "COM9"
	PortCOM10 enumPort = "COM10"
)

// enumRate 串口波特率枚举
type enumRate int

const (
	Rate1200   enumRate = 1200
	Rate2400   enumRate = 2400
	Rate4800   enumRate = 4800
	Rate9600   enumRate = 9600
	Rate19200  enumRate = 19200
	Rate38400  enumRate = 38400
	Rate57600  enumRate = 57600
	Rate115200 enumRate = 115200
)

// enumDataBits 串口数据位枚举
type enumDataBits int

const (
	DataBits5 enumDataBits = 5
	DataBits6 enumDataBits = 6
	DataBits7 enumDataBits = 7
	DataBits8 enumDataBits = 8
)

// enumStopBits 串口停止位枚举
type enumStopBits int

const (
	StopBits1 enumStopBits = 1
	StopBits2 enumStopBits = 2
)

// enumParity 串口校验位枚举
type enumParity string

const (
	ParityNone  enumParity = "N"
	ParityEven  enumParity = "E"
	ParityOdd   enumParity = "O"
	ParityMark  enumParity = "M"
	ParitySpace enumParity = "S"
)
