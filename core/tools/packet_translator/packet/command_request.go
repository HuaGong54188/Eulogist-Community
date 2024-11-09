package packet

import (
	neteaseProtocol "Eulogist/core/minecraft/netease/protocol"
	neteasePacket "Eulogist/core/minecraft/netease/protocol/packet"

	standardProtocol "Eulogist/core/minecraft/standard/protocol"
	standardPacket "Eulogist/core/minecraft/standard/protocol/packet"
)

type CommandRequest struct{}

func (pk *CommandRequest) ToNetEasePacket(standard standardPacket.Packet) neteasePacket.Packet {
	p := neteasePacket.CommandRequest{}
	input := standard.(*standardPacket.CommandRequest)

	p.CommandLine = input.CommandLine
	p.CommandOrigin = neteaseProtocol.CommandOrigin(input.CommandOrigin)
	p.Internal = input.Internal
	p.Version = input.Version

	p.UnLimited = false

	return &p
}

func (pk *CommandRequest) ToStandardPacket(netease neteasePacket.Packet) standardPacket.Packet {
	p := standardPacket.CommandRequest{}
	input := netease.(*neteasePacket.CommandRequest)

	p.CommandLine = input.CommandLine
	p.CommandOrigin = standardProtocol.CommandOrigin(input.CommandOrigin)
	p.Internal = input.Internal
	p.Version = input.Version

	return &p
}
