package mc_client

import (
	raknet_connection "Eulogist/core/raknet"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net"

	standardProtocol "Eulogist/core/minecraft/standard/protocol"

	"github.com/sandertv/go-raknet"
)

// CreateListener 在 127.0.0.1:19132 上以 Raknet
// 协议侦听 Minecraft 客户端的连接，
// 这意味着您成功创建了一个 Minecraft 数据包代理服务器
func (m *MinecraftClient) CreateListener() error {
	// 创建一个 Raknet 监听器
	listener, err := raknet.Listen("127.0.0.1:19132")
	if err != nil {
		return fmt.Errorf("创建监听器: %v", err)
	}
	// 获取监听器的地址
	address, ok := listener.Addr().(*net.UDPAddr)
	if !ok {
		return fmt.Errorf("创建监听器: 获取监听器地址失败(请确认你是不是改源码出错了)")
	}
	// 设置 pong data
	listener.PongData([]byte(
		fmt.Sprintf(
			"MCPE;%v;%v;%v;%v;%v;%v;Gophertunnel;%v;%v;%v;%v;",
			"Eulogist", standardProtocol.CurrentProtocol, standardProtocol.CurrentVersion, "0", "1",
			listener.ID(), "Creative", 1, address.Port, address.Port,
		),
	))
	// 初始化变量
	m.listener = listener
	m.connected = make(chan struct{}, 1)
	m.Address = address
	m.Conn = raknet_connection.NewStandardRaknetWrapper()
	// 返回成功
	return nil
}

// WaitConnect 等待 Minecraft 客户端连接到服务器
func (m *MinecraftClient) WaitConnect() error {
	// 接受客户端连接
	conn, err := m.listener.Accept()
	if err != nil {
		return fmt.Errorf("等待连接: %v", err)
	}
	// 丢弃其他连接
	go func() {
		for {
			conn, err := m.listener.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	// 初始化变量
	serverKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	m.Conn.SetConnection(conn, serverKey)
	m.connected <- struct{}{}
	// 返回成功
	return nil
}
