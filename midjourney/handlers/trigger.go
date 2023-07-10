package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp/v3"
)

type RequestTrigger struct {
	Type         string `json:"type"`
	DiscordMsgId string `json:"discordMsgId,omitempty"`
	MsgHash      string `json:"msgHash,omitempty"`
	Prompt       string `json:"prompt,omitempty"`
	Index        int64  `json:"index,omitempty"`
	SubType      string `json:"subType,omitempty"`
}

func MidjourneyBot(c *gin.Context) {
	var body RequestTrigger
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pp.Println(body)

	var err error
	switch body.Type {
	case "generate":
		err = GenerateImage(body.Prompt)
	case "upscale":
		err = ImageUpscale(body.Index, body.DiscordMsgId, body.MsgHash)
	case "variation":
		err = ImageVariation(body.Index, body.DiscordMsgId, body.MsgHash)
	case "maxUpscale":
		err = ImageMaxUpscale(body.DiscordMsgId, body.MsgHash, body.SubType)
	case "zoomout":
		err = ImageZoomOut(body.DiscordMsgId, body.MsgHash, body.SubType)
	case "pan":
		err = ImagePan(body.DiscordMsgId, body.MsgHash, body.SubType)
	case "reset":
		err = ImageReset(body.DiscordMsgId, body.MsgHash)
	case "describe":
		err = ImageDescribe(body.Prompt)
	default:
		err = errors.New("invalid type")
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
