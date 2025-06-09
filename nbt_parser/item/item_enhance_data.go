package nbt_parser_item

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/mitchellh/mapstructure"
)

type SingleItemEnch struct {
	ID    int16 `mapstructure:"id"`
	Level int16 `mapstructure:"lvl"`
}

// Marshal ..
func (s *SingleItemEnch) Marshal(io protocol.IO) {
	io.Int16(&s.ID)
	io.Int16(&s.Level)
}

func parseItemEnchList(enchList []any) (result []SingleItemEnch, err error) {
	for _, value := range enchList {
		var singleItemEnch SingleItemEnch

		val, ok := value.(map[string]any)
		if !ok {
			continue
		}

		err = mapstructure.Decode(&val, &singleItemEnch)
		if err != nil {
			return nil, fmt.Errorf("ParseItemEnchList: %v", err)
		}

		result = append(result, singleItemEnch)
	}
	return
}

func ParseItemEnchList(nbtMap map[string]any) (result []SingleItemEnch, err error) {
	tag, ok := nbtMap["tag"].(map[string]any)
	if !ok {
		return
	}

	ench, ok := tag["ench"].([]any)
	if !ok {
		return
	}

	result, err = parseItemEnchList(ench)
	if err != nil {
		return nil, fmt.Errorf("ParseItemEnchList: %v", err)
	}

	return
}

func ParseItemEnchListNetwork(item protocol.ItemStack) (result []SingleItemEnch, err error) {
	if item.NBTData == nil {
		return
	}

	ench, ok := item.NBTData["ench"].([]any)
	if !ok {
		return
	}

	result, err = parseItemEnchList(ench)
	if err != nil {
		return nil, fmt.Errorf("ParseItemEnchListNetwork: %v", err)
	}

	return
}
