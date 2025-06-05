package block_helper

// BaseBlockHelper 描述了各种 Helper 方块的相邻方块
type NearBlockHelper struct {
	Name string
}

func (n NearBlockHelper) BlockName() string {
	return n.Name
}

func (NearBlockHelper) BlockStates() map[string]any {
	return map[string]any{}
}

func (NearBlockHelper) BlockStatesString() string {
	return "[]"
}
