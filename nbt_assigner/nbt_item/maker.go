package nbt_item

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	nbt_assigner_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	nbt_parser_item "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/item"
)

func init() {
	nbt_assigner_interface.MakeNBTItemMethod = MakeNBTItemMethod
	nbt_assigner_interface.EnchMultiple = EnchMultiple
	nbt_assigner_interface.RenameMultiple = RenameMultiple
	nbt_assigner_interface.EnchAndRenameMultiple = EnchAndRenameMultiple
}

func MakeNBTItemMethod(
	console *nbt_console.Console,
	multipleItems ...nbt_parser_interface.Item,
) (result nbt_assigner_interface.Item, supported bool) {
	if len(multipleItems) == 0 {
		return nil, false
	}

	switch multipleItems[0].(type) {
	case *nbt_parser_item.Book:
		result = &Book{api: console}
	case *nbt_parser_item.Banner:
		result = &Banner{
			api:             console,
			maxSlotCanUse:   BannerMaxSlotCanUse,
			maxBannerToMake: BannerMaxBannerToMake,
		}
	case *nbt_parser_item.Shield:
		result = &Shield{api: console}
	default:
		return nil, false
	}

	result.Append(multipleItems...)
	return result, true
}

func EnchMultiple(
	console *nbt_console.Console,
	multipleItems [27]*nbt_parser_interface.Item,
) error {
	api := console.API()

	enchItems := make([]resources_control.SlotID, 0)
	enchItemsCount := make(map[resources_control.SlotID]uint8)

	for index, value := range multipleItems {
		if value == nil {
			continue
		}

		slotID := resources_control.SlotID(index + 9)
		defaultItem := (*value).UnderlyingItem().(*nbt_parser_item.DefaultItem)

		if len(defaultItem.Enhance.EnchList) > 0 {
			enchItems = append(enchItems, slotID)
			enchItemsCount[slotID] = defaultItem.ItemCount()
		}
	}

	if len(enchItems) > 0 {
		success, err := api.ContainerOpenAndClose().OpenInventory()
		if err != nil {
			return fmt.Errorf("EnchMultiple: %v", err)
		}
		if !success {
			return fmt.Errorf("EnchMultiple: Failed to open the inventory")
		}
		defer api.ContainerOpenAndClose().CloseContainer()
	}

	for {
		if len(enchItems) == 0 {
			break
		}

		currentRound := enchItems[0:min(len(enchItems), 9)]
		transaction := api.ItemStackOperation().OpenTransaction()

		for dstSlotID, srcSlotID := range currentRound {
			_ = transaction.MoveBetweenInventory(
				srcSlotID,
				resources_control.SlotID(dstSlotID),
				enchItemsCount[srcSlotID],
			)
		}

		success, _, _, err := transaction.Commit()
		if err != nil {
			return fmt.Errorf("EnchMultiple: %v", err)
		}
		if !success {
			return fmt.Errorf("EnchMultiple: The server rejected the item stack operation (Ench stage 1)")
		}

		for index, originSlotID := range currentRound {
			item := multipleItems[originSlotID]
			defaultItem := (*item).UnderlyingItem().(*nbt_parser_item.DefaultItem)

			currentSlotID := resources_control.SlotID(index)
			if console.HotbarSlotID() != currentSlotID {
				err = api.BotClick().ChangeSelectedHotbarSlot(currentSlotID)
				if err != nil {
					return fmt.Errorf("EnchMultiple: %v", err)
				}
				console.UpdateHotbarSlotID(currentSlotID)
			}

			for _, ench := range defaultItem.Enhance.EnchList {
				err = api.Commands().SendSettingsCommand(fmt.Sprintf("enchant @s %d %d", ench.ID, ench.Level), true)
				if err != nil {
					return fmt.Errorf("EnchMultiple: %v", err)
				}
			}
		}

		err = api.Commands().AwaitChangesGeneral()
		if err != nil {
			return fmt.Errorf("EnchMultiple: %v", err)
		}

		for currentSlotID, originSlotID := range currentRound {
			_ = transaction.MoveBetweenInventory(
				resources_control.SlotID(currentSlotID),
				originSlotID,
				enchItemsCount[originSlotID],
			)
		}

		success, _, _, err = transaction.Commit()
		if err != nil {
			return fmt.Errorf("EnchMultiple: %v", err)
		}
		if !success {
			return fmt.Errorf("EnchMultiple: The server rejected the item stack operation (Ench stage 2)")
		}

		enchItems = enchItems[len(currentRound):]
	}

	return nil
}

func RenameMultiple(
	console *nbt_console.Console,
	multipleItems [27]*nbt_parser_interface.Item,
) error {
	api := console.API()

	renameItems := make([]resources_control.SlotID, 0)
	renameItemsNewName := make([]string, 0)

	for index, value := range multipleItems {
		if value == nil {
			continue
		}

		slotID := resources_control.SlotID(index + 9)
		defaultItem := (*value).UnderlyingItem().(*nbt_parser_item.DefaultItem)
		displayName := defaultItem.Enhance.DisplayName

		if len(displayName) > 0 {
			renameItems = append(renameItems, slotID)
			renameItemsNewName = append(renameItemsNewName, displayName)
		}
	}

	if len(renameItems) == 0 {
		return nil
	}

	index, err := console.FindOrGenerateNewAnvil()
	if err != nil {
		return fmt.Errorf("RenameMultiple: %v", err)
	}

	success, err := console.OpenContainerByIndex(index)
	if err != nil {
		return fmt.Errorf("RenameMultiple: %v", err)
	}
	if !success {
		return fmt.Errorf("RenameMultiple: Failed to open the anvil")
	}
	defer api.ContainerOpenAndClose().CloseContainer()

	transaction := api.ItemStackOperation().OpenTransaction()
	for index, slotID := range renameItems {
		_ = transaction.RenameInventoryItem(
			slotID,
			renameItemsNewName[index],
		)
	}

	success, _, _, err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("RenameMultiple: %v", err)
	}
	if !success {
		return fmt.Errorf("RenameMultiple: The server rejected the renaming operation")
	}

	return nil
}

func EnchAndRenameMultiple(
	console *nbt_console.Console,
	multipleItems [27]*nbt_parser_interface.Item,
) error {
	err := EnchMultiple(console, multipleItems)
	if err != nil {
		return fmt.Errorf("EnchAndRenameMultiple: %v", err)
	}
	err = RenameMultiple(console, multipleItems)
	if err != nil {
		return fmt.Errorf("EnchAndRenameMultiple: %v", err)
	}
	return nil
}
