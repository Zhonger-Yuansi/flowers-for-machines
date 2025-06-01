package game_interface

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
	"github.com/google/uuid"
)

// StructureBackup 是基于 Commands 包装的结构备份与恢复相关的实现
type StructureBackup struct {
	api *Commands
}

// NewStructureBackup 根据 api 返回并创建一个新的 StructureBackup
func NewStructureBackup(api *Commands) *StructureBackup {
	return &StructureBackup{api: api}
}

// BackupStructure 通过使用 structure 命令保存 pos 处的方块。
// 返回的 uuid 是标识该结构的唯一标识符
func (s *StructureBackup) BackupStructure(pos protocol.BlockPos) (result uuid.UUID, err error) {
	api := s.api

	uniqueId := uuid.New()
	request := fmt.Sprintf(
		`structure save "%s" %d %d %d %d %d %d`,
		utils.MakeUUIDSafeString(uniqueId),
		pos[0], pos[1], pos[2],
		pos[0], pos[1], pos[2],
	)
	resp, isTimeout, err := api.SendWSCommandWithTimeout(request, DefaultTimeoutCommandRequest)

	if isTimeout {
		err = api.SendSettingsCommand(request, true)
		if err != nil {
			return result, fmt.Errorf("BackupStructure: %v", err)
		}
		err = api.AwaitChangesGeneral()
		if err != nil {
			return result, fmt.Errorf("BackupStructure: %v", err)
		}
		return uniqueId, nil
	}

	if err != nil {
		return result, fmt.Errorf("BackupStructure: %v", err)
	}

	if resp.SuccessCount == 0 {
		return result, fmt.Errorf(
			"BackupStructure: Backup (%d,%d,%d) failed because the success count of the command %#v is 0",
			pos[0], pos[1], pos[2], request,
		)
	}

	return uniqueId, nil
}

// DeleteStructure 删除标识符为 uniqueID 的结构。
// 即便目标结构不存在，此函数在通常情况下也仍然会返回空错误
func (s *StructureBackup) DeleteStructure(uniqueID uuid.UUID) error {
	err := s.api.SendSettingsCommand(
		fmt.Sprintf(
			`structure delete "%v"`,
			utils.MakeUUIDSafeString(uniqueID),
		),
		false,
	)
	if err != nil {
		return fmt.Errorf("DeleteStructure: %v", err)
	}
	return nil
}

// RevertAndDeleteStructure 在 pos 处恢复先前备份的结构，
// 其中，uniqueID 是该结构的唯一标识符
func (s *StructureBackup) RevertStructure(uniqueID uuid.UUID, pos protocol.BlockPos) error {
	api := s.api
	request := fmt.Sprintf(
		`structure load "%v" %d %d %d`,
		utils.MakeUUIDSafeString(uniqueID),
		pos[0],
		pos[1],
		pos[2],
	)
	resp, isTimeOut, err := api.SendWSCommandWithTimeout(request, DefaultTimeoutCommandRequest)

	if isTimeOut {
		err = api.SendSettingsCommand(request, true)
		if err != nil {
			return fmt.Errorf("RevertStructure: %v", err)
		}
		err = api.AwaitChangesGeneral()
		if err != nil {
			return fmt.Errorf("RevertStructure: %v", err)
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("RevertStructure: %v", err)
	}

	if resp.SuccessCount == 0 {
		return fmt.Errorf(
			"RevertStructure: Revert structure %#v on (%d,%d,%d) failed because the success count of the command %#v is 0",
			utils.MakeUUIDSafeString(uniqueID), pos[0], pos[1], pos[2], request,
		)
	}

	return nil
}
