package block_helper

// AnvilBlockHelper 描述了一个铁砧
type AnvilBlockHelper struct{}

func (AnvilBlockHelper) BlockName() string {
	return "minecraft:anvil"
}

func (AnvilBlockHelper) BlockStates() map[string]any {
	return map[string]any{
		"damage":                       "undamaged",
		"minecraft:cardinal_direction": "east",
	}
}

func (AnvilBlockHelper) BlockStatesString() string {
	return `["damage"="undamaged","minecraft:cardinal_direction"="east"]`
}
