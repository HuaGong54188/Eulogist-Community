package packet

import (
	neteaseProtocol "Eulogist/core/minecraft/netease/protocol"
	neteasePacket "Eulogist/core/minecraft/netease/protocol/packet"

	standardPacket "Eulogist/core/minecraft/standard/protocol/packet"
)

type CompletedUsingItem struct{}

func (pk *CompletedUsingItem) ToNetEasePacket(standard standardPacket.Packet) neteasePacket.Packet {
	p := neteasePacket.CompletedUsingItem{}
	input := standard.(*standardPacket.CompletedUsingItem)

	p.UsedItemID = input.UsedItemID
	p.UseMethod = input.UseMethod

	p.UnknownItem = neteaseProtocol.ItemInstance{}

	return &p
}

func (pk *CompletedUsingItem) ToStandardPacket(netease neteasePacket.Packet) standardPacket.Packet {
	p := standardPacket.CompletedUsingItem{}
	input := netease.(*neteasePacket.CompletedUsingItem)

	p.UsedItemID = input.UsedItemID
	p.UseMethod = input.UseMethod

	return &p
}
