package handlers

import (
	"midjourney/services"
)

func GenerateImage(prompt string) error {
	err := services.GenerateImage(prompt)
	return err
}

func ImageUpscale(index int64, discordMsgId string, msgHash string) error {
	err := services.Upscale(index, discordMsgId, msgHash)
	return err
}

func ImageVariation(index int64, discordMsgId string, msgHash string) error {
	err := services.Variate(index, discordMsgId, msgHash)
	return err
}

func ImageMaxUpscale(discordMsgId string, msgHash string, subType string) error {
	err := services.MaxUpscale(discordMsgId, msgHash, subType)
	return err
}

func ImageZoomOut(discordMsgId string, msgHash string, subType string) error {
	err := services.ZoomOut(discordMsgId, msgHash, subType)
	return err
}

func ImagePan(discordMsgId string, msgHash string, subType string) error {
	err := services.Pan(discordMsgId, msgHash, subType)
	return err
}

func ImageReset(discordMsgId string, msgHash string) error {
	err := services.Reset(discordMsgId, msgHash)
	return err
}

func ImageDescribe(uploadName string) error {
	err := services.Describe(uploadName)
	return err
}
