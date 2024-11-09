package mc_server

import (
	fb_client "Eulogist/core/fb_auth/mv4/client"
	"Eulogist/core/minecraft/netease/protocol/packet"
	raknet_wrapper "Eulogist/core/raknet/wrapper"
	"Eulogist/proxy/persistence_data"
)

// ...
type MinecraftServer struct {
	fbClient              *fb_client.Client
	authResponse          *fb_client.AuthResponse
	getCheckNumEverPassed bool

	PersistenceData *persistence_data.PersistenceData
	Conn            *raknet_wrapper.Raknet[packet.Packet]
}

// ...
type BasicConfig struct {
	ServerCode     string
	ServerPassword string
	Token          string
	AuthServer     string
}