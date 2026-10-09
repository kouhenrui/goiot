package device

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"goiot/internal/domain/device"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *device.DeviceRepository
}

func NewAuthService(repo *device.DeviceRepository) *AuthService {
	return &AuthService{repo: repo}
}

type RegisterMQTTDeviceResult struct {
	DeviceKey string `json:"device_key"`
	Secret    string `json:"secret"`
}

// randomToken 生成密码学安全的随机字符串。
func randomToken(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// RegisterMQTTDevice 注册 MQTT 设备。
// Secret 只在本次响应中返回，数据库只保存 bcrypt 哈希。
func (s *AuthService) RegisterMQTTDevice(
	ctx context.Context,
	productID uint,
	name string,
) (*RegisterMQTTDeviceResult, error) {
	keySuffix, err := randomToken(6)
	if err != nil {
		return nil, fmt.Errorf("generate device key: %w", err)
	}

	secret, err := randomToken(32)
	if err != nil {
		return nil, fmt.Errorf("generate device secret: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(secret),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("hash device secret: %w", err)
	}

	d := &device.Device{
		ProductID:      productID,
		Name:           name,
		DeviceKey:      "dev" + keySuffix,
		ConnectionType: device.ConnectionMQTT,
		SecretHash:     string(hash),
		Status:         device.StatusInactive,
	}

	if err := s.repo.Create(ctx, d); err != nil {
		return nil, fmt.Errorf("create device: %w", err)
	}

	return &RegisterMQTTDeviceResult{
		DeviceKey: d.DeviceKey,
		Secret:    secret,
	}, nil
}

// AuthenticateMQTT 验证 MQTT 客户端身份。
func (s *AuthService) AuthenticateMQTT(
	ctx context.Context,
	deviceKey string,
	clientID string,
	secret string,
) bool {
	if deviceKey == "" || secret == "" || clientID != deviceKey {
		return false
	}

	d, err := s.repo.FindByKey(ctx, deviceKey)
	if err != nil {
		return false
	}

	if d.ConnectionType != device.ConnectionMQTT {
		return false
	}

	if d.Status == device.StatusDisabled || d.SecretHash == "" {
		return false
	}

	return bcrypt.CompareHashAndPassword(
		[]byte(d.SecretHash),
		[]byte(secret),
	) == nil
}
