package nbt_block

import (
	"fmt"
	"strings"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
	nbt_assigner_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_item"
	nbt_parser_block "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/block"
	nbt_hash "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/hash"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	nbt_parser_item "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/item"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
)

type Container struct {
	NBTBlockBase
	data nbt_parser_block.Container
}

func (Container) Offset() protocol.BlockPos {
	return protocol.BlockPos{0, 0, 0}
}

// spawnContainer 在操作台中心生成空的 container。
// 确保这个容器的自定义物品名称也被考虑在内
func (c *Container) spawnContainer(container nbt_parser_block.Container) error {
	// 准备
	api := c.console.API()
	useCommandToPlaceBlock := true
	successFunc := func() {
		c.console.UseHelperBlock(nbt_console.RequesterUser, nbt_console.ConsoleIndexCenterBlock, block_helper.ContainerBlockHelper{
			OpenInfo: block_helper.ContainerBlockOpenInfo{
				Name:                  container.BlockName(),
				States:                container.BlockStates(),
				ConsiderOpenDirection: container.ConsiderOpenDirection(),
				ShulkerFacing:         container.NBT.ShulkerFacing,
			},
			IsEmpty: len(container.NBT.Items) == 0,
		})
	}

	// 尝试基容器缓存
	hit, err := c.cache.BaseContainerCache().LoadCache(container.BlockName(), container.BlockStates(), container.CustomName)
	if err != nil {
		return fmt.Errorf("spawnContainer: %v", err)
	}
	if hit {
		return nil
	}

	// 先将目标位置替换为空气
	err = api.SetBlock().SetBlock(c.console.Center(), "minecraft:air", "[]")
	if err != nil {
		return fmt.Errorf("spawnContainer: %v", err)
	}
	c.console.UseHelperBlock(nbt_console.RequesterUser, nbt_console.ConsoleIndexCenterBlock, block_helper.Air{})

	// 检查是否需要复杂的工序来放置容器
	if len(c.data.CustomName) > 0 {
		useCommandToPlaceBlock = false
	}
	if strings.Contains(c.data.BlockName(), "shulker") {
		if c.data.NBT.ShulkerFacing != 1 {
			useCommandToPlaceBlock = false
		}
	}

	// 如果需要复杂的工序
	if !useCommandToPlaceBlock {
		// 先把目标物品获取到物品栏
		err := api.Replaceitem().ReplaceitemInInventory(
			"@s",
			game_interface.ReplacePathHotbarOnly,
			game_interface.ReplaceitemInfo{
				Name:     c.data.BlockName(),
				Count:    1,
				MetaData: 0,
				Slot:     c.console.HotbarSlotID(),
			},
			"",
			true,
		)
		if err != nil {
			return fmt.Errorf("spawnContainer: %v", err)
		}
		c.console.UseInventorySlot(nbt_console.RequesterUser, c.console.HotbarSlotID(), true)

		// 这个容器具有自定义的物品名称，需要进一步特殊处理
		if len(c.data.CustomName) > 0 {
			index, err := c.console.FindOrGenerateNewAnvil()
			if err != nil {
				return fmt.Errorf("spawnContainer: %v", err)
			}

			success, err := c.console.OpenContainerByIndex(index)
			if err != nil {
				return fmt.Errorf("spawnContainer: %v", err)
			}
			if !success {
				return fmt.Errorf("spawnContainer: Failed to open the anvil to rename container")
			}

			success, _, _, err = api.ItemStackOperation().OpenTransaction().
				RenameInventoryItem(c.console.HotbarSlotID(), c.data.CustomName).
				Commit()
			if err != nil {
				_ = api.ContainerOpenAndClose().CloseContainer()
				return fmt.Errorf("spawnContainer: %v", err)
			}
			if !success {
				_ = api.ContainerOpenAndClose().CloseContainer()
				return fmt.Errorf("spawnContainer: The server rejected the container rename operation")
			}

			err = api.ContainerOpenAndClose().CloseContainer()
			if err != nil {
				return fmt.Errorf("spawnContainer: %v", err)
			}
		}

		// 放置目标容器
		_, offsetPos, err := api.BotClick().PlaceBlockHighLevel(c.console.Center(), c.console.HotbarSlotID(), c.data.NBT.ShulkerFacing)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		successFunc()
		c.console.UpdatePosition(c.console.Center())
		*c.console.NearBlockByIndex(nbt_console.ConsoleIndexCenterBlock, offsetPos) = block_helper.NearBlock{
			Name: game_interface.BasePlaceBlock,
		}

		// 将该容器保存到基容器缓存命中系统
		err = c.cache.BaseContainerCache().StoreCache(container.CustomName)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}

		// 返回值
		return nil
	}

	// 目标容器可以直接通过简单的 setblock 放置
	err = api.SetBlock().SetBlock(c.console.Center(), container.BlockName(), container.BlockStatesString())
	if err != nil {
		return fmt.Errorf("spawnContainer: %v", err)
	}
	successFunc()

	return nil
}

// itemTransition 将已置于操作台中心的 srcContainer 转移为 dstContainer。
// 应当保证 dstContainer 中的物品可以只通过移动 srcContainer 中的物品得到
func (c *Container) itemTransition(
	srcContainer nbt_parser_block.Container,
	dstContainer nbt_parser_block.Container,
) error {
	api := c.console.API()

	// 清空背包
	_, err := api.Commands().SendWSCommandWithResp("clear")
	if err != nil {
		return fmt.Errorf("itemTransition: %v", err)
	}
	c.console.CleanInventory()

	// 打开加载好的容器
	success, err := c.console.OpenContainerByIndex(nbt_console.ConsoleIndexCenterBlock)
	if err != nil {
		return fmt.Errorf("itemTransition: %v", err)
	}
	if !success {
		return fmt.Errorf("itemTransition: Failed to open the container for %#v (stage 1)", srcContainer)
	}

	// 将该容器内的物品移动到背包
	transaction := api.ItemStackOperation().OpenTransaction()
	for index, item := range srcContainer.NBT.Items {
		_ = transaction.MoveToInventory(
			resources_control.SlotID(item.Slot),
			resources_control.SlotID(index),
			item.Item.ItemCount(),
		)
	}

	// 提交更改
	success, _, _, err = transaction.Commit()
	if err != nil {
		_ = api.ContainerOpenAndClose().CloseContainer()
		return fmt.Errorf("itemTransition: %v", err)
	}
	if !success {
		_ = api.ContainerOpenAndClose().CloseContainer()
		return fmt.Errorf("itemTransition: The server rejected the stack request action")
	}

	// 关闭容器
	err = api.ContainerOpenAndClose().CloseContainer()
	if err != nil {
		return fmt.Errorf("itemTransition: %v", err)
	}

	// 生成新容器
	err = c.spawnContainer(dstContainer)
	if err != nil {
		return fmt.Errorf("itemTransition: %v", err)
	}

	// 打开新容器
	success, err = c.console.OpenContainerByIndex(nbt_console.ConsoleIndexCenterBlock)
	if err != nil {
		return fmt.Errorf("itemTransition: %v", err)
	}
	if !success {
		return fmt.Errorf("itemTransition: Failed to open the container for %#v (stage 2)", srcContainer)
	}
	defer api.ContainerOpenAndClose().CloseContainer()

	// 准备
	itemTypeIndex := game_interface.ItemType(0)
	itemTypeMapping := make(map[uint64]game_interface.ItemType)
	src := make([]game_interface.ItemInfoWithSlot, 0)
	dst := make([]game_interface.ItemInfoWithSlot, 0)

	// 处理源
	for _, item := range srcContainer.NBT.Items {
		hashNumber := nbt_hash.NBTItemTypeHash(item.Item)

		if _, ok := itemTypeMapping[hashNumber]; !ok {
			itemTypeMapping[hashNumber] = itemTypeIndex
			itemTypeIndex++
		}

		src = append(src, game_interface.ItemInfoWithSlot{
			Slot: resources_control.SlotID(item.Slot),
			ItemInfo: game_interface.ItemInfo{
				Count:    item.Item.ItemCount(),
				ItemType: itemTypeMapping[hashNumber],
			},
		})
	}

	// 处理目的地
	for _, item := range dstContainer.NBT.Items {
		hashNumber := nbt_hash.NBTItemTypeHash(item.Item)

		if _, ok := itemTypeMapping[hashNumber]; !ok {
			itemTypeMapping[hashNumber] = itemTypeIndex
			itemTypeIndex++
		}

		dst = append(dst, game_interface.ItemInfoWithSlot{
			Slot: resources_control.SlotID(item.Slot),
			ItemInfo: game_interface.ItemInfo{
				Count:    item.Item.ItemCount(),
				ItemType: itemTypeMapping[hashNumber],
			},
		})
	}

	// 进行物品状态转移
	success, err = api.ItemTransition().TransitionToContainer(src, dst)
	if err != nil {
		return fmt.Errorf("itemTransition: %v", err)
	}
	if !success {
		return fmt.Errorf("itemTransition: Failed to do transition")
	}

	return nil
}

func (c *Container) Make() error {
	api := c.console.API()

	// Step 1: 检查该容器是否命中集合校验和
	{
		// 尝试从底层缓存命中系统加载
		structure, hit, isSetHashHit, err := c.cache.NBTBlockCache().LoadCache(
			nbt_hash.CompletelyHashNumber{
				HashNumber:    nbt_hash.NBTBlockHash(&c.data),
				SetHashNumber: nbt_hash.ContainerSetHash(&c.data),
			},
		)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		if hit {
			panic("Make: Should nerver happened")
		}

		// 如果我们命中了集合哈希校验和
		if isSetHashHit {
			container, ok := structure.Block.(*nbt_parser_block.Container)
			if !ok {
				panic("Make: Should nerver happened")
			}

			err = c.itemTransition(*container, c.data)
			if err != nil {
				return fmt.Errorf("Make: %v", err)
			}

			return nil
		}
	}

	// Step 2: 构造物品树 (仅限复杂物品或需要处理的子方块)
	itemTypeIndex := game_interface.ItemType(0)
	itemTypes := make(map[uint64]game_interface.ItemType)
	itemGroups := make(map[uint64][]nbt_parser_block.ItemWithSlot)
	for _, item := range c.data.NBT.Items {
		if !item.Item.IsComplex() {
			continue
		}
		hashNumber := nbt_hash.NBTItemNBTHash(item.Item)
		itemGroups[hashNumber] = append(itemGroups[hashNumber], item)
		if _, ok := itemTypes[hashNumber]; !ok {
			itemTypes[hashNumber] = itemTypeIndex
			itemTypeIndex++
		}
	}

	// Step 3.1: 找出部分命中和没有命中的需要处理的子方块 (找出集合)
	allSubBlocks := make([]int, 0)
	allSubBlocksSet := make(map[uint64]bool)
	subBlockPartHit := make([]int, 0)
	subBlockNotHit := make([]int, 0)
	for index, item := range c.data.NBT.Items {
		underlying := item.Item.UnderlyingItem().(*nbt_parser_item.DefaultItem)
		if underlying.Block.SubBlock == nil || !underlying.Block.SubBlock.NeedSpecialHandle() {
			continue
		}

		hashNumber := nbt_hash.NBTItemNBTHash(item.Item)
		if _, ok := allSubBlocksSet[hashNumber]; ok {
			continue
		}
		allSubBlocksSet[hashNumber] = true

		_, hit, partHit := c.cache.NBTBlockCache().CheckCache(nbt_hash.CompletelyHashNumber{
			HashNumber:    nbt_hash.NBTBlockHash(underlying.Block.SubBlock),
			SetHashNumber: nbt_hash.ContainerSetHash(underlying.Block.SubBlock),
		})

		if !hit && partHit {
			subBlockPartHit = append(subBlockPartHit, index)
		}
		if !hit && !partHit {
			subBlockNotHit = append(subBlockNotHit, index)
		}
		allSubBlocks = append(allSubBlocks, index)
	}

	// Step 3.2: 处理部分命中的子方块 (容器)
	for _, index := range subBlockPartHit {
		item := c.data.NBT.Items[index]
		underlying := item.Item.UnderlyingItem().(*nbt_parser_item.DefaultItem)

		wantContainer, ok := underlying.Block.SubBlock.(*nbt_parser_block.Container)
		if !ok {
			panic("Make: Should nerver happened")
		}

		structure, _, partHit, err := c.cache.NBTBlockCache().LoadCache(nbt_hash.CompletelyHashNumber{
			HashNumber:    nbt_hash.NBTBlockHash(wantContainer),
			SetHashNumber: nbt_hash.ContainerSetHash(wantContainer),
		})
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		if !partHit {
			panic("Make: Should nerver happened")
		}

		container, ok := structure.Block.(*nbt_parser_block.Container)
		if !ok {
			panic("Make: Should nerver happened")
		}

		err = c.itemTransition(*container, *wantContainer)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}

		err = c.cache.NBTBlockCache().StoreCache(wantContainer, c.console.Center())
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
	}

	// Step 3.3: 处理没有命中的子方块
	for _, index := range subBlockNotHit {
		item := c.data.NBT.Items[index]
		underlying := item.Item.UnderlyingItem().(*nbt_parser_item.DefaultItem)

		wantContainer, ok := underlying.Block.SubBlock.(*nbt_parser_block.Container)
		if !ok {
			panic("Make: Should nerver happened")
		}

		_, _, _, err := nbt_assigner_interface.PlaceNBTBlock(c.console, c.cache, wantContainer)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
	}

	// Step 4: 生成当前容器
	err := c.spawnContainer(c.data)
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}

	// Step 5: 将子方块放入容器
	if len(allSubBlocks) > 0 {
		// 清空物品栏
		_, err := api.Commands().SendWSCommandWithResp("clear")
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}

		// 占用所有物品栏，
		// 因为我们无法确保数据匹配
		for index := range 36 {
			c.console.UseInventorySlot(nbt_console.RequesterUser, resources_control.SlotID(index), true)
		}

		// 通过 Pick block 得到所有的子方块
		for _, index := range allSubBlocks {
			underlying := c.data.NBT.Items[index].Item.UnderlyingItem()
			subBlock := underlying.(*nbt_parser_item.DefaultItem).Block.SubBlock

			structure, hit, partHit := c.cache.NBTBlockCache().CheckCache(nbt_hash.CompletelyHashNumber{
				HashNumber:    nbt_hash.NBTBlockHash(subBlock),
				SetHashNumber: nbt_hash.ContainerSetHash(subBlock),
			})
			if !hit || partHit {
				panic("Make: Should nerver happened")
			}

			index, _, block := c.console.FindSpaceToPlaceNewContainer(false, true)
			if block == nil {
				index = nbt_console.ConsoleIndexFirstHelperBlock
			}

			err = api.StructureBackup().RevertStructure(
				structure.UniqueID,
				c.console.BlockPosByIndex(index),
			)
			if err != nil {
				return fmt.Errorf("Make: %v", err)
			}
			c.console.UseHelperBlock(nbt_console.RequesterUser, index, block_helper.ComplexBlock{
				Name:   structure.Block.BlockName(),
				States: structure.Block.BlockStates(),
			})

			err = c.console.CanReachOrMove(c.console.BlockPosByIndex(index))
			if err != nil {
				return fmt.Errorf("Make: %v", err)
			}

			success, currentSlot, err := api.BotClick().PickBlock(c.console.BlockPosByIndex(index), true)
			if err != nil || !success {
				_ = api.BotClick().ChangeSelectedHotbarSlot(nbt_console.DefaultHotbarSlot)
				c.console.UpdateHotbarSlotID(nbt_console.DefaultHotbarSlot)
			}
			if err != nil {
				return fmt.Errorf("Make: %v", err)
			}
			if !success {
				return fmt.Errorf("Make: Failed to get sub block %#v by pick block", subBlock)
			}
			c.console.UpdateHotbarSlotID(currentSlot)
		}

		// c.console.API().Commands().AwaitChangesGeneral() // idk ?

		// 现在所有子方块都被 Pick Block 到背包了
		allItemStack, inventoryExisted := api.Resources().Inventories().GetAllItemStack(0)
		if !inventoryExisted {
			panic("Make: Should nerver happened")
		}
		if len(allItemStack) != len(allSubBlocks) {
			panic("Make: Should nerver happened")
		}

		// 打开操作台中心处容器
		success, err := c.console.OpenContainerByIndex(nbt_console.ConsoleIndexCenterBlock)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		if !success {
			return fmt.Errorf("Make: Failed to open the container %#v when move sub block in it", c.data)
		}

		// 将背包中的每个子方块移动到对应的父节点处
		transaction := api.ItemStackOperation().OpenTransaction()
		for srcSlot, value := range allItemStack {
			newItem, err := nbt_parser_item.ParseItemNetwork(
				value.Stack,
				api.Resources().ConstantPacket().ItemNameByNetworkID(value.Stack.NetworkID),
			)
			if err != nil {
				_ = api.ContainerOpenAndClose().CloseContainer()
				return fmt.Errorf("Make: %v", err)
			}

			hashNumber := nbt_hash.NBTItemNBTHash(newItem)
			dstSlot := itemGroups[hashNumber][0].Slot

			_ = transaction.MoveToContainer(srcSlot, resources_control.SlotID(dstSlot), 1)
		}

		// 提交更改
		success, _, _, err = transaction.Commit()
		if err != nil {
			_ = api.ContainerOpenAndClose().CloseContainer()
			return fmt.Errorf("Make: %v", err)
		}
		if !success {
			_ = api.ContainerOpenAndClose().CloseContainer()
			return fmt.Errorf("Make: The server rejected the stack request action when move sub block in it")
		}

		// 关闭容器
		err = api.ContainerOpenAndClose().CloseContainer()
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
	}

	// Step 6.1: 计算出哪些物品是需要制作的非子方块复杂物品
	complexItemExcludeSubBlock := make([]nbt_parser_interface.Item, 0)
	for _, value := range itemGroups {
		if _, ok := value[0].Item.(*nbt_parser_item.DefaultItem); ok {
			continue
		}
		complexItemExcludeSubBlock = append(complexItemExcludeSubBlock, value[0].Item)
	}

	// Step 6.2: 制作非子方块的复杂物品
	variousItems := nbt_assigner_interface.MakeNBTItemMethod(c.console, c.cache, complexItemExcludeSubBlock...)
	for _, item := range variousItems {
		for {
			resultSlot, err := item.Make()
			if err != nil {
				return fmt.Errorf("Make: %v", err)
			}
			if len(resultSlot) == 0 {
				break
			}

			success, err := c.console.OpenContainerByIndex(nbt_console.ConsoleIndexCenterBlock)
			if err != nil {
				return fmt.Errorf("Make: %v", err)
			}
			if !success {
				return fmt.Errorf("Make: Failed to open the container %#v when make complex item", c.data)
			}

			transaction := api.ItemStackOperation().OpenTransaction()
			for hashNumber, slotID := range resultSlot {
				dstSlot := resources_control.SlotID(itemGroups[hashNumber][0].Slot)
				_ = transaction.MoveToContainer(slotID, dstSlot, 1)
			}

			success, _, _, err = transaction.Commit()
			if err != nil {
				_ = api.ContainerOpenAndClose().CloseContainer()
				return fmt.Errorf("Make: %v", err)
			}
			if !success {
				_ = api.ContainerOpenAndClose().CloseContainer()
				return fmt.Errorf("Make: The server rejected item stack request action when make complex item")
			}
			for _, slotID := range resultSlot {
				c.console.UseInventorySlot(nbt_console.RequesterUser, slotID, false)
			}

			err = api.ContainerOpenAndClose().CloseContainer()
			if err != nil {
				return fmt.Errorf("Make: %v", err)
			}
		}
	}

	// Step 7.1: 检测是否需要物品分裂
	needItemCopy := false
	for _, value := range itemGroups {
		if len(value) > 0 {
			needItemCopy = true
			break
		}
		if value[0].Item.ItemCount() > 1 {
			needItemCopy = true
			break
		}
	}

	// Step 7.2: 物品分裂 (复杂物品复制)
	if needItemCopy {
		// 清理背包
		_, err = api.Commands().SendWSCommandWithResp("clear")
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		c.console.CleanInventory()

		// 打开容器
		success, err := c.console.OpenContainerByIndex(nbt_console.ConsoleIndexCenterBlock)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		if !success {
			return fmt.Errorf("Make: Failed to open the container %#v when do item copy", c.data)
		}

		// 将容器中现存的所有物品拿回
		transaction := api.ItemStackOperation().OpenTransaction()
		for _, value := range itemGroups {
			_ = transaction.MoveToInventory(
				resources_control.SlotID(value[0].Slot),
				resources_control.SlotID(value[0].Slot),
				1,
			)
		}

		// 提交更改
		success, _, _, err = transaction.Commit()
		if err != nil {
			_ = api.ContainerOpenAndClose().CloseContainer()
			return fmt.Errorf("Make: %v", err)
		}
		if !success {
			_ = api.ContainerOpenAndClose().CloseContainer()
			return fmt.Errorf("Make: Failed to move item from container %#v when do item copy", c.data)
		}

		// 关闭容器
		err = api.ContainerOpenAndClose().CloseContainer()
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}

		// 构造基物品
		baseItems := make([]game_interface.ItemInfoWithSlot, 0)
		for _, value := range itemGroups {
			baseItems = append(baseItems, game_interface.ItemInfoWithSlot{
				Slot: resources_control.SlotID(value[0].Slot),
				ItemInfo: game_interface.ItemInfo{
					Count:    1,
					ItemType: itemTypes[nbt_hash.NBTItemNBTHash(value[0].Item)],
				},
			})
		}

		// 构造蓝图
		targetItems := make([]*game_interface.ItemInfo, 27)
		for _, value := range itemGroups {
			for _, val := range value {
				targetItems[val.Slot] = &game_interface.ItemInfo{
					Count:    val.Item.ItemCount(),
					ItemType: itemTypes[nbt_hash.NBTItemNBTHash(val.Item)],
				}
			}
		}

		// 物品分裂
		err = api.ItemCopy().CopyItem(
			game_interface.UseItemOnBlocks{
				HotbarSlotID: c.console.HotbarSlotID(),
				BlockPos:     c.console.Center(),
				BlockName:    c.data.BlockName(),
				BlockStates:  c.data.BlockStates(),
			},
			baseItems, targetItems,
		)
		if err != nil {
			api.Commands().SendWSCommandWithResp("clear")
			c.console.CleanInventory()
			return fmt.Errorf("Make: %v", err)
		}

		// 清理背包
		_, err = api.Commands().SendWSCommandWithResp("clear")
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		c.console.CleanInventory()
	}

	// Step 8.1: 填充剩余物品
	for _, item := range c.data.NBT.Items {
		if item.Item.IsComplex() {
			continue
		}
		underlying := item.Item.UnderlyingItem().(*nbt_parser_item.DefaultItem)

		err = api.Replaceitem().ReplaceitemInContainerAsync(
			c.console.Center(),
			game_interface.ReplaceitemInfo{
				Name:     item.Item.ItemName(),
				Count:    item.Item.ItemCount(),
				MetaData: item.Item.ItemMetadata(),
				Slot:     resources_control.SlotID(item.Slot),
			},
			utils.MarshalItemComponent(underlying.Enhance.ItemComponent),
		)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
	}

	// Step 8.2: 等待更改
	err = api.Commands().AwaitChangesGeneral()
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}

	// Step 9.1: 找出所有需要修改物品名称或需要附魔的物品
	enchOrRenameList := make([]int, 0)
	for index, value := range c.data.NBT.Items {
		if value.Item.NeedEnchOrRename() {
			enchOrRenameList = append(enchOrRenameList, index)
		}
	}

	// Step 9.2: 将需要修改物品名称或需要附魔的物品移动到背包
	if len(enchOrRenameList) > 0 {
		success, err := c.console.OpenContainerByIndex(nbt_console.ConsoleIndexCenterBlock)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		if !success {
			return fmt.Errorf("Make: Failed to open the container %#v when do ench or rename operation", c.data)
		}

		transaction := api.ItemStackOperation().OpenTransaction()
		for _, index := range enchOrRenameList {
			item := c.data.NBT.Items[index]
			_ = transaction.MoveToInventory(
				resources_control.SlotID(item.Slot),
				resources_control.SlotID(item.Slot+9),
				item.Item.ItemCount(),
			)
		}

		success, _, _, err = transaction.Commit()
		if err != nil {
			_ = api.ContainerOpenAndClose().CloseContainer()
			return fmt.Errorf("Make: %v", err)
		}
		if !success {
			_ = api.ContainerOpenAndClose().CloseContainer()
			return fmt.Errorf("Make: The server rejected the stack request action when do ench or rename operation")
		}

		err = api.ContainerOpenAndClose().CloseContainer()
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
	}

	// Step 9.3: 物品附魔或重命名操作
	if len(enchOrRenameList) > 0 {
		multipleItems := [27]*nbt_parser_interface.Item{}
		for _, index := range enchOrRenameList {
			item := c.data.NBT.Items[index]
			multipleItems[item.Slot] = &item.Item
		}
		err = nbt_item.EnchAndRenameMultiple(c.console, multipleItems)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
	}

	// Step 9.4: 将物品移动回容器
	if len(enchOrRenameList) > 0 {
		success, err := c.console.OpenContainerByIndex(nbt_console.ConsoleIndexCenterBlock)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		if !success {
			return fmt.Errorf("Make: Failed to open the container %#v when finish ench or rename operation", c.data)
		}

		transaction := api.ItemStackOperation().OpenTransaction()
		for _, index := range enchOrRenameList {
			item := c.data.NBT.Items[index]
			_ = transaction.MoveToContainer(
				resources_control.SlotID(item.Slot+9),
				resources_control.SlotID(item.Slot),
				item.Item.ItemCount(),
			)
		}

		success, _, _, err = transaction.Commit()
		if err != nil {
			_ = api.ContainerOpenAndClose().CloseContainer()
			return fmt.Errorf("Make: %v", err)
		}
		if !success {
			_ = api.ContainerOpenAndClose().CloseContainer()
			return fmt.Errorf("Make: The server rejected the stack request action when finish ench or rename operation")
		}

		err = api.ContainerOpenAndClose().CloseContainer()
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
	}

	// Step 10: 返回值
	return nil
}
