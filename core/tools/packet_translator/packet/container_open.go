package packet

import (
	neteaseProtocol "Eulogist/core/minecraft/netease/protocol"
	neteasePacket "Eulogist/core/minecraft/netease/protocol/packet"

	standardProtocol "Eulogist/core/minecraft/standard/protocol"
	standardPacket "Eulogist/core/minecraft/standard/protocol/packet"
)

type ContainerOpen struct{}

func (pk *ContainerOpen) ToNetEasePacket(standard standardPacket.Packet) neteasePacket.Packet {
	p := neteasePacket.ContainerOpen{}
	input := standard.(*standardPacket.ContainerOpen)

	p.WindowID = input.WindowID
	p.ContainerType = input.ContainerType
	p.ContainerPosition = neteaseProtocol.BlockPos(input.ContainerPosition)
	p.ContainerEntityUniqueID = input.ContainerEntityUniqueID

	p.Unknown1 = false

	return &p
}

func (pk *ContainerOpen) ToStandardPacket(netease neteasePacket.Packet) standardPacket.Packet {
	p := standardPacket.ContainerOpen{}
	input := netease.(*neteasePacket.ContainerOpen)

	p.WindowID = input.WindowID
	p.ContainerType = input.ContainerType
	p.ContainerPosition = standardProtocol.BlockPos(input.ContainerPosition)
	p.ContainerEntityUniqueID = input.ContainerEntityUniqueID

	return &p
}
