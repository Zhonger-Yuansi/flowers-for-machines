package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/nbt"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
	"github.com/gin-gonic/gin"
)

func CheckAlive(c *gin.Context) {
	c.Writer.WriteString("Still Alive")
}

func ProcessExist(c *gin.Context) {
	_ = mcClient.Conn().Close()
	go func() {
		time.Sleep(time.Second)
		os.Exit(0)
	}()
}

func PlaceNBTBlock(c *gin.Context) {
	var request PlaceNBTBlockRequest
	var blockNBT map[string]any

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusOK, PlaceNBTBlockResponse{
			Success:   false,
			ErrorType: ResponseErrorTypeParseError,
			ErrorInfo: fmt.Sprintf("Failed to parse request; err = %v", err),
		})
	}

	blockNBTBytes, err := base64.StdEncoding.DecodeString(request.BlockNBTBase64String)
	if err != nil {
		c.JSON(http.StatusOK, PlaceNBTBlockResponse{
			Success:   false,
			ErrorType: ResponseErrorTypeParseError,
			ErrorInfo: fmt.Sprintf("Failed to parse block NBT base64 string; err = %v", err),
		})
	}
	err = nbt.UnmarshalEncoding(blockNBTBytes, &blockNBT, nbt.LittleEndian)
	if err != nil {
		c.JSON(http.StatusOK, PlaceNBTBlockResponse{
			Success:   false,
			ErrorType: ResponseErrorTypeParseError,
			ErrorInfo: fmt.Sprintf("Block NBT bytes is broken; err = %v", err),
		})
	}

	canFast, uniqueID, offset, err := wrapper.PlaceNBTBlock(
		request.BlockName,
		utils.ParseBlockStatesString(request.BlockStatesString),
		blockNBT,
	)
	if err != nil {
		c.JSON(http.StatusOK, PlaceNBTBlockResponse{
			Success:   false,
			ErrorType: ResponseErrorTypeRuntimeError,
			ErrorInfo: fmt.Sprintf("Runtime error: Failed to place NBT block; err = %v", err),
		})
	}

	c.JSON(http.StatusOK, PlaceNBTBlockResponse{
		Success:           true,
		CanFast:           canFast,
		StructureUniqueID: uniqueID.String(),
		StructureName:     utils.MakeUUIDSafeString(uniqueID),
		OffsetX:           offset.X(),
		OffsetY:           offset.Y(),
		OffsetZ:           offset.Z(),
	})
}
