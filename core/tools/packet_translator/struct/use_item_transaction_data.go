package packet_translate_struct

import (
	neteaseProtocol "Eulogist/core/minecraft/netease/protocol"
	standardProtocol "Eulogist/core/minecraft/standard/protocol"
	packet_translate_pool "Eulogist/core/tools/packet_translator/pool"
)

// 将 from 转换为 neteaseProtocol.UseItemTransactionData
func ConvertToNetEaseUseItemTransactionData(from *standardProtocol.UseItemTransactionData) *neteaseProtocol.UseItemTransactionData {
	return &neteaseProtocol.UseItemTransactionData{
		LegacyRequestID: from.LegacyRequestID,
		LegacySetItemSlots: ConvertSlice(
			from.LegacySetItemSlots,
			func(from standardProtocol.LegacySetItemSlot) neteaseProtocol.LegacySetItemSlot {
				return neteaseProtocol.LegacySetItemSlot{
					ContainerID: packet_translate_pool.StandardContainerIDToNetEaseContainerID[from.ContainerID],
					Slots:       from.Slots,
				}
			},
		),
		Actions: ConvertSlice(
			from.Actions,
			ConvertToNetEaseInventoryAction,
		),
		ActionType:      from.ActionType,
		BlockPosition:   neteaseProtocol.BlockPos(from.BlockPosition),
		BlockFace:       from.BlockFace,
		HotBarSlot:      from.HotBarSlot,
		HeldItem:        ConvertToNetEaseItemInstance(from.HeldItem),
		Position:        from.Position,
		ClickedPosition: from.ClickedPosition,
		BlockRuntimeID:  from.BlockRuntimeID,
	}
}

// 将 from 转换为 standardProtocol.UseItemTransactionData
func ConvertToStandardUseItemTransactionData(from *neteaseProtocol.UseItemTransactionData) *standardProtocol.UseItemTransactionData {
	return &standardProtocol.UseItemTransactionData{
		LegacyRequestID: from.LegacyRequestID,
		LegacySetItemSlots: ConvertSlice(
			from.LegacySetItemSlots,
			func(from neteaseProtocol.LegacySetItemSlot) standardProtocol.LegacySetItemSlot {
				return standardProtocol.LegacySetItemSlot{
					ContainerID: packet_translate_pool.NetEaseContainerIDStandardContainerID[from.ContainerID],
					Slots:       from.Slots,
				}
			},
		),
		Actions: ConvertSlice(
			from.Actions,
			ConvertToStandardInventoryAction,
		),
		ActionType:      from.ActionType,
		BlockPosition:   standardProtocol.BlockPos(from.BlockPosition),
		BlockFace:       from.BlockFace,
		HotBarSlot:      from.HotBarSlot,
		HeldItem:        ConvertToStandardItemInstance(from.HeldItem),
		Position:        from.Position,
		ClickedPosition: from.ClickedPosition,
		BlockRuntimeID:  from.BlockRuntimeID,
	}
}
