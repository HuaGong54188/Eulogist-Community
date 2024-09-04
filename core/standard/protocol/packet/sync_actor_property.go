package packet

import (
	"Eulogist/core/standard/nbt"
	"Eulogist/core/standard/protocol"
)

// SyncActorProperty is an alternative to synced actor data.
type SyncActorProperty struct {
	// PropertyData ...
	PropertyData map[string]any
}

// ID ...
func (*SyncActorProperty) ID() uint32 {
	return IDSyncActorProperty
}

func (pk *SyncActorProperty) Marshal(io protocol.IO) {
	io.NBT(&pk.PropertyData, nbt.NetworkLittleEndian)
}
