package nbt_parser_block

import (
	"fmt"
	"strings"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/mapping"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/df-mc/worldupgrader/blockupgrader"
)

// ParseNBTBlock 从方块实体数据 blockNBT 解析一个方块实体。
// blockName 和 blockStates 分别指示这个方块实体的名称和方块状态
func ParseNBTBlock(blockName string, blockStates map[string]any, blockNBT map[string]any) (block nbt_parser_interface.Block, err error) {
	name := strings.ToLower(blockName)
	if !strings.HasPrefix(name, "minecraft:") {
		name = "minecraft:" + name
	}

	newBlock := blockupgrader.Upgrade(blockupgrader.BlockState{
		Name:       name,
		Properties: blockStates,
	})
	defaultBlock := DefaultBlock{
		Name:   newBlock.Name,
		States: newBlock.Properties,
	}

	blockType, ok := mapping.SupportBlocksPool[newBlock.Name]
	if !ok {
		return &defaultBlock, nil
	}

	switch blockType {
	case mapping.SupportNBTBlockTypeCommandBlock:
		block = &CommandBlock{DefaultBlock: defaultBlock}
	case mapping.SupportNBTBlockTypeContainer:
		block = &Container{DefaultBlock: defaultBlock}
	case mapping.SupportNBTBlockTypeSign:
		block = &Sign{DefaultBlock: defaultBlock}
	case mapping.SupportNBTBlockTypeFrame:
		block = &Frame{DefaultBlock: defaultBlock}
	case mapping.SupportNBTBlockTypeStructureBlock:
		block = &StructureBlock{DefaultBlock: defaultBlock}
	case mapping.SupportNBTBlockTypeBanner:
		block = &Banner{DefaultBlock: defaultBlock}
	case mapping.SupportNBTBlockTypeLectern:
		block = &Lectern{DefaultBlock: defaultBlock}
	case mapping.SupportNBTBlockTypeJukeBox:
		block = &JukeBox{DefaultBlock: defaultBlock}
	case mapping.SupportNBTBlockTypeBrewingStand:
		block = &BrewingStand{DefaultBlock: defaultBlock}
	default:
		panic("ParseNBTBlock: Should nerver happened")
	}

	err = block.Parse(blockNBT)
	if err != nil {
		return nil, fmt.Errorf("ParseNBTBlock: %v", err)
	}
	return block, nil
}

func init() {
	nbt_parser_interface.ParseBlock = ParseNBTBlock
}
