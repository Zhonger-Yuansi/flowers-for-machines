package client

import (
	"context"
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/bunker/auth"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/py_rpc"
	cts "github.com/Happy2018new/the-last-problem-of-the-humankind/core/py_rpc/mod_event/client_to_server"
	cts_mc "github.com/Happy2018new/the-last-problem-of-the-humankind/core/py_rpc/mod_event/client_to_server/minecraft"
	cts_mc_p "github.com/Happy2018new/the-last-problem-of-the-humankind/core/py_rpc/mod_event/client_to_server/minecraft/preset"
	cts_mc_v "github.com/Happy2018new/the-last-problem-of-the-humankind/core/py_rpc/mod_event/client_to_server/minecraft/vip_event_system"
	mei "github.com/Happy2018new/the-last-problem-of-the-humankind/core/py_rpc/mod_event/interface"
)

// createConnection 通过 authenticator 连接到租赁服，
// 并初始化 Minecraft 连接
func createConnection(
	ctx context.Context,
	authenticator minecraft.Authenticator,
) (conn *minecraft.Conn, err error) {
	// prepare
	var dialer minecraft.Dialer
	var authResponse auth.AuthResponse

	// create connection
	dialer = minecraft.Dialer{
		Authenticator: authenticator,
	}
	conn, authResponse, err = dialer.DialContext(ctx, "raknet")
	if err != nil {
		return nil, err
	}

	// get constant and send pre-login packet
	runtimeid := fmt.Sprintf("%d", conn.GameData().EntityUniqueID)
	conn.WritePacket(&packet.ClientCacheStatus{
		Enabled: false,
	})

	// send netease related packet
	conn.WritePacket(&packet.NeteaseJson{
		Data: []byte(
			fmt.Sprintf(
				`{"eventName":"LOGIN_UID","resid":"","uid":"%d"}`,
				conn.IdentityData().Uid,
			),
		),
	})

	{
		modUUIDs := make([]any, 0)
		botComponent := make(map[string]int64, 0)
		for modUUID, outfitType := range authResponse.BotComponent {
			modUUIDs = append(modUUIDs, modUUID)
			if outfitType != nil {
				botComponent[modUUID] = int64(*outfitType)
			}
		}
		conn.WritePacket(&packet.PyRpc{
			Value: py_rpc.Marshal(&py_rpc.SyncUsingMod{
				modUUIDs,
				conn.ClientData().SkinID,
				"",
				true,
				botComponent,
			}),
			OperationType: packet.PyRpcOperationTypeSend,
		})
	}

	conn.WritePacket(&packet.PyRpc{
		Value:         py_rpc.Marshal(&py_rpc.SyncVipSkinUUID{nil}),
		OperationType: packet.PyRpcOperationTypeSend,
	})
	conn.WritePacket(&packet.PyRpc{
		Value:         py_rpc.Marshal(&py_rpc.ClientLoadAddonsFinishedFromGac{}),
		OperationType: packet.PyRpcOperationTypeSend,
	})

	{
		event := cts_mc_p.GetLoadedInstances{PlayerRuntimeID: runtimeid}
		module := cts_mc.Preset{Module: &mei.DefaultModule{Event: &event}}
		park := cts.Minecraft{Default: mei.Default{Module: &module}}
		conn.WritePacket(&packet.PyRpc{
			Value: py_rpc.Marshal(&py_rpc.ModEvent{
				Package: &park,
				Type:    py_rpc.ModEventClientToServer,
			}),
			OperationType: packet.PyRpcOperationTypeSend,
		})
	}

	conn.WritePacket(&packet.PyRpc{
		Value:         py_rpc.Marshal(&py_rpc.ArenaGamePlayerFinishLoad{}),
		OperationType: packet.PyRpcOperationTypeSend,
	})

	{
		event := cts_mc_v.PlayerUiInit{RuntimeID: runtimeid}
		module := cts_mc.VIPEventSystem{Module: &mei.DefaultModule{Event: &event}}
		park := cts.Minecraft{Default: mei.Default{Module: &module}}
		conn.WritePacket(&packet.PyRpc{
			Value: py_rpc.Marshal(&py_rpc.ModEvent{
				Package: &park,
				Type:    py_rpc.ModEventClientToServer,
			}),
			OperationType: packet.PyRpcOperationTypeSend,
		})
	}

	// return
	return
}
