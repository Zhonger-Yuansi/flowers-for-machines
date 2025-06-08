package nbt_parser_block

import (
	"fmt"
	"strings"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/mapping"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/df-mc/worldupgrader/blockupgrader"
)

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
