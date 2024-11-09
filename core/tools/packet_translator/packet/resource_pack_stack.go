package packet

import (
	neteaseProtocol "Eulogist/core/minecraft/netease/protocol"
	neteasePacket "Eulogist/core/minecraft/netease/protocol/packet"
	packet_translate_struct "Eulogist/core/tools/packet_translator/struct"

	standardProtocol "Eulogist/core/minecraft/standard/protocol"
	standardPacket "Eulogist/core/minecraft/standard/protocol/packet"
)

type ResourcePackStack struct{}

func (pk *ResourcePackStack) ToNetEasePacket(standard standardPacket.Packet) neteasePacket.Packet {
	p := neteasePacket.ResourcePackStack{}
	input := standard.(*standardPacket.ResourcePackStack)

	p.TexturePackRequired = input.TexturePackRequired
	p.BaseGameVersion = input.BaseGameVersion
	p.ExperimentsPreviouslyToggled = input.ExperimentsPreviouslyToggled

	p.BehaviourPacks = packet_translate_struct.ConvertSlice(
		input.BehaviourPacks,
		func(from standardProtocol.StackResourcePack) neteaseProtocol.StackResourcePack {
			return neteaseProtocol.StackResourcePack(from)
		},
	)
	p.TexturePacks = packet_translate_struct.ConvertSlice(
		input.TexturePacks,
		func(from standardProtocol.StackResourcePack) neteaseProtocol.StackResourcePack {
			return neteaseProtocol.StackResourcePack(from)
		},
	)
	p.Experiments = packet_translate_struct.ConvertSlice(
		input.Experiments,
		func(from standardProtocol.ExperimentData) neteaseProtocol.ExperimentData {
			return neteaseProtocol.ExperimentData(from)
		},
	)

	p.Unknown1 = false
	p.Unknown2 = false

	return &p
}

func (pk *ResourcePackStack) ToStandardPacket(netease neteasePacket.Packet) standardPacket.Packet {
	p := standardPacket.ResourcePackStack{}
	input := netease.(*neteasePacket.ResourcePackStack)

	p.TexturePackRequired = input.TexturePackRequired
	p.BaseGameVersion = input.BaseGameVersion
	p.ExperimentsPreviouslyToggled = input.ExperimentsPreviouslyToggled

	p.BehaviourPacks = packet_translate_struct.ConvertSlice(
		input.BehaviourPacks,
		func(from neteaseProtocol.StackResourcePack) standardProtocol.StackResourcePack {
			return standardProtocol.StackResourcePack(from)
		},
	)
	p.TexturePacks = packet_translate_struct.ConvertSlice(
		input.TexturePacks,
		func(from neteaseProtocol.StackResourcePack) standardProtocol.StackResourcePack {
			return standardProtocol.StackResourcePack(from)
		},
	)
	p.Experiments = packet_translate_struct.ConvertSlice(
		input.Experiments,
		func(from neteaseProtocol.ExperimentData) standardProtocol.ExperimentData {
			return standardProtocol.ExperimentData(from)
		},
	)

	return &p
}
