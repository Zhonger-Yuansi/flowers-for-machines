package packet

import (
	"github.com/OmineDev/flowers-for-machines/core/minecraft/nbt"
	"github.com/OmineDev/flowers-for-machines/core/minecraft/protocol"
)

// EditorNetwork is a packet sent from the server to the client and vise-versa to communicate editor-mode related
// information. It carries a single compound tag containing the relevant information.
type EditorNetwork struct {
	// Payload is a network little endian compound tag holding data relevant to the editor.
	Payload map[string]any
}

// ID ...
func (*EditorNetwork) ID() uint32 {
	return IDEditorNetwork
}

func (pk *EditorNetwork) Marshal(io protocol.IO) {
	io.NBT(&pk.Payload, nbt.NetworkLittleEndian)
}
