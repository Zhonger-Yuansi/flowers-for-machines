package block_helper

// AirBlock 描述了一个空气方块
type AirBlock struct{}

func (AirBlock) BlockName() string {
	return "minecraft:air"
}

func (AirBlock) BlockStates() map[string]any {
	return map[string]any{}
}

func (AirBlock) BlockStatesString() string {
	return `[]`
}
