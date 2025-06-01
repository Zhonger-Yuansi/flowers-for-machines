package game_interface

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/TriM-Organization/bedrock-world-operator/block"
	"github.com/go-gl/mathgl/mgl32"
)

// UseItemOnBlocks 是机器人在使
// 用手持物品对方块进行操作时的通用结构体
type UseItemOnBlocks struct {
	HotbarSlotID resources_control.SlotID // 指代机器人当前已选择的快捷栏编号
	BlockPos     protocol.BlockPos        // 指代被操作方块的位置
	BlockName    string                   // 指代被操作方块的名称
	BlockStates  map[string]any           // 指代被操作方块的方块状态
}

// BotClick 是基于 ResourcesWrapper
// 和 Commands 实现的已简化的点击实现。
//
// 由于点击操作与机器人手持物品强相关，
// 本处也集成了切换手持物品的实现
type BotClick struct {
	r *ResourcesWrapper
	c *Commands
}

// NewBotClick 基于 wrapper 和 commands 创建并返回一个新的 BotClick
func NewBotClick(wrapper *ResourcesWrapper, commands *Commands) *BotClick {
	return &BotClick{r: wrapper, c: commands}
}

// 切换客户端的手持物品栏为 hotBarSlotID 。
// 若提供的 hotBarSlotID 大于 8 ，则会重定向为 0
func (b *BotClick) ChangeSelectedHotbarSlot(hotbarSlotID uint8) error {
	if hotbarSlotID > 8 {
		hotbarSlotID = 0
	}

	err := b.r.WritePacket(&packet.PlayerHotBar{
		SelectedHotBarSlot: uint32(hotbarSlotID),
		WindowID:           0,
		SelectHotBarSlot:   true,
	})
	if err != nil {
		return fmt.Errorf("ChangeSelectedHotbarSlot: %v", err)
	}

	return nil
}

// clickBlock ..
func (b *BotClick) clickBlock(
	request UseItemOnBlocks,
	blockFace int32,
	position mgl32.Vec3,
) error {
	// Step 1: 取得被点击方块的方块运行时 ID
	blockRuntimeID, found := block.StateToRuntimeID(request.BlockName, request.BlockStates)
	if !found {
		return fmt.Errorf(
			"clickBlock: Can't found the block runtime ID of block %#v (block states = %#v)",
			request.BlockName, request.BlockStates,
		)
	}

	// Step 2: 取得当前手持物品的信息
	item, inventoryExisted := b.r.Inventories().GetItemStack(0, request.HotbarSlotID)
	if !inventoryExisted {
		return fmt.Errorf("clickBlock: Should never happened")
	}

	// Step 3: 发送点击操作
	err := b.r.WritePacket(&packet.InventoryTransaction{
		LegacyRequestID:    0,
		LegacySetItemSlots: []protocol.LegacySetItemSlot(nil),
		Actions:            []protocol.InventoryAction{},
		TransactionData: &protocol.UseItemTransactionData{
			LegacyRequestID:    0,
			LegacySetItemSlots: nil,
			Actions:            nil,
			ActionType:         protocol.UseItemActionClickBlock,
			BlockPosition:      request.BlockPos,
			BlockFace:          blockFace,
			HotBarSlot:         int32(request.HotbarSlotID),
			HeldItem:           *item,
			Position:           position,
			BlockRuntimeID:     blockRuntimeID,
		},
	})
	if err != nil {
		return fmt.Errorf("clickBlock: %v", err)
	}

	// Step 4: 额外操作 (自 v1.20.50 以外的必须更改)
	{
		// !!! NOTE - MUST SEND AUTH INPUT TWICE !!!
		// await changes and send auth
		// input to submit changes
		err = b.c.AwaitChangesGeneral()
		if err != nil {
			return fmt.Errorf("clickBlock: %v", err)
		}
		err = b.r.WritePacket(&packet.PlayerAuthInput{InputData: packet.InputFlagStartFlying})
		if err != nil {
			return fmt.Errorf("clickBlock: %v", err)
		}
		err = b.c.AwaitChangesGeneral()
		if err != nil {
			return fmt.Errorf("clickBlock: %v", err)
		}
		err = b.r.WritePacket(&packet.PlayerAuthInput{InputData: packet.InputFlagStartFlying})
		if err != nil {
			return fmt.Errorf("clickBlock: %v", err)
		}
	}

	return nil
}

/*
让客户端点击 request 所指代的方块，
并且指定当次交互时玩家的位置为 position 。

position 不一定需要是真实的，
客户端可以上传欺骗性的数据，
服务器不会对它们进行验证。

该函数在通常情况下被用于十分精细的操作，
例如为告示牌的特定面附加发光效果。

此函数不会自动切换物品栏，也不会等待租赁服响应更改
*/
func (b *BotClick) ClickBlockWitchPosition(
	request UseItemOnBlocks,
	position mgl32.Vec3,
) error {
	err := b.clickBlock(request, 0, position)
	if err != nil {
		return fmt.Errorf("ClickBlockWitchPosition: %v", err)
	}
	return nil
}

/*
让客户端点击 request 所指代的方块。

你可以对容器使用这样的操作，这会使得容器被打开。

你亦可以对物品展示框使用这样的操作，
这会使得物品被放入或令展示框内的物品旋转。

此函数不会自动切换物品栏，也不会等待租赁服响应更改
*/
func (b *BotClick) ClickBlock(request UseItemOnBlocks) error {
	err := b.clickBlock(request, 0, mgl32.Vec3{})
	if err != nil {
		return fmt.Errorf("ClickBlock: %v", err)
	}
	return nil
}

// 使用快捷栏 hotbarSlotID 进行一次空点击操作。
// 此函数不会自动切换物品栏，也不会等待租赁服响应更改
func (b *BotClick) ClickAir(hotbarSlotID resources_control.SlotID) error {
	// Step 1: 获取手持物品栏物品数据信息
	item, inventoryExisted := b.r.Inventories().GetItemStack(0, hotbarSlotID)
	if !inventoryExisted {
		return fmt.Errorf("ClickAir: Should never happened")
	}

	// Step 2: 发送点击数据包
	err := b.r.WritePacket(
		&packet.InventoryTransaction{
			TransactionData: &protocol.UseItemTransactionData{
				ActionType: protocol.UseItemActionClickAir,
				HotBarSlot: int32(hotbarSlotID),
				HeldItem:   *item,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("ClickAir: %v", err)
	}

	// Step 3: 额外操作 (自 v1.20.50 以外的必须更改)
	err = b.c.AwaitChangesGeneral()
	if err != nil {
		return fmt.Errorf("ClickAir: %v", err)
	}
	err = b.r.WritePacket(&packet.PlayerAuthInput{InputData: packet.InputFlagStartFlying})
	if err != nil {
		return fmt.Errorf("ClickAir: %v", err)
	}

	return nil
}

/*
PlaceBlock 使客户端创建一个新方块。

request 指代实际被点击的方块，但这并不代表新方块被创建的位置。
我们通过点击 request 处的方块，并指定点击的面为 blockFace ，
然后租赁服根据这些信息，在另外相应的位置创建这些新的方块。

此函数不会自动切换物品栏，也不会等待租赁服响应更改
*/
func (b *BotClick) PlaceBlock(
	request UseItemOnBlocks,
	blockFace int32,
) error {
	err := b.clickBlock(request, blockFace, mgl32.Vec3{})
	if err != nil {
		return fmt.Errorf("PlaceBlock: %v", err)
	}
	return nil
}
