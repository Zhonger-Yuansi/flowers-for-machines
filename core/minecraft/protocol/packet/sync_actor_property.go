package packet

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/nbt"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
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
