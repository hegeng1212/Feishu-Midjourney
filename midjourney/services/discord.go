package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	config "midjourney/initialization"
	"net/http"
	"path/filepath"

	"github.com/k0kubun/pp/v3"
)

const (
	url             = "https://discord.com/api/v9/interactions"
	uploadUrlFormat = "https://discord.com/api/v9/channels/%s/attachments"
	DiscordId       = "938956540159881230"
	appId           = "936929561302675456"//"1117660105962438706"
	//Version         = "1077969938624553050"
	Version         = "1237876415471554623"//"1166847114203123795"
	SessionId       = "11c73404e48de26b7c76a31349c556c2"//"cb06f61453064c0983f2adae2a88c223"
)

type ImageResponse struct {
	Message string `json:"message"`
}

func GenerateImage(prompt string) error {
	requestBody := ReqTriggerDiscord{
		Type:          2,
		GuildID:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelID:     config.GetConfig().DISCORD_CHANNEL_ID,
		ApplicationId: appId,
		SessionId:     SessionId,
		Data: DSCommand{
			Version: Version,
			Id:      DiscordId,
			Name:    "imagine",
			Type:    1,
			Options: []DSOption{{Type: 3, Name: "prompt", Value: prompt}},
			ApplicationCommand: DSApplicationCommand{
				Id:                       DiscordId,
				ApplicationId:            appId,
				Version:                  Version,
				DefaultPermission:        true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "imagine",
				Description:              "Create images with Midjourney",
				DmPermission:             true,
				Options:                  []DSCommandOption{{Type: 3, Name: "prompt", Description: "The prompt to imagine", Required: true}},
			},
			Attachments: []ReqCommandAttachments{},
		},
	}
	r, err := request(requestBody, url)
	pp.Println(url)
	pp.Println(requestBody)
	pp.Println(r)
	return err
}

func Upscale(index int64, messageId string, messageHash string) error {
	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionId,
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::upsample::%d::%s", index, messageHash),
		},
	}
	_, err := request(requestBody, url)
	return err
}

func MaxUpscale(messageId string, messageHash string, subType string) error {

	customIdTemplate := "MJ::JOB::variation::1::%s::SOLO"
	switch subType {
	case "high":
		customIdTemplate = "MJ::JOB::high_variation::1::%s::SOLO"
	case "low":
		customIdTemplate = "MJ::JOB::low_variation::1::%s::SOLO"
	}

	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionId,
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf(customIdTemplate, messageHash),
		},
	}

	data, _ := json.Marshal(requestBody)

	fmt.Println("max upscale request body: ", string(data))

	_, err := request(requestBody, url)
	return err
}

func ZoomOut(messageId string, messageHash string, subType string) error {
	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionId,
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::Outpaint::%s::1::%s::SOLO", subType, messageHash),
		},
	}

	data, _ := json.Marshal(requestBody)

	fmt.Println("zoom out request body: ", string(data))

	_, err := request(requestBody, url)
	return err
}

func Pan(messageId string, messageHash string, subType string) error {

	customIdTemplate := "MJ::JOB::pan_left::1::%s::SOLO"
	switch subType {
	case "left":
		customIdTemplate = "MJ::JOB::pan_left::1::%s::SOLO"
	case "right":
		customIdTemplate = "MJ::JOB::pan_right::1::%s::SOLO"
	case "up":
		customIdTemplate = "MJ::JOB::pan_up::1::%s::SOLO"
	case "down":
		customIdTemplate = "MJ::JOB::pan_down::1::%s::SOLO"
	}


	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionId,
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf(customIdTemplate, messageHash),
		},
	}

	data, _ := json.Marshal(requestBody)

	fmt.Println("pan request body: ", string(data))

	_, err := request(requestBody, url)
	return err
}

func Variate(index int64, messageId string, messageHash string) error {
	requestBody := ReqVariationDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionId,
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::%d::%s", index, messageHash),
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Reset(messageId string, messageHash string) error {
	requestBody := ReqResetDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionId,
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::reroll::0::%s::SOLO", messageHash),
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Describe(uploadName string) error {
	requestBody := ReqTriggerDiscord{
		Type:          2,
		GuildID:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelID:     config.GetConfig().DISCORD_CHANNEL_ID,
		ApplicationId: appId,
		SessionId:     SessionId,
		Data: DSCommand{
			Version: Version,
			Id:      DiscordId,
			Name:    "describe",
			Type:    1,
			Options: []DSOption{{Type: 11, Name: "image", Value: 0}},
			ApplicationCommand: DSApplicationCommand{
				Id:                       "1092492867185950852",
				ApplicationId:            "936929561302675456",
				Version:                  "1092492867185950853",
				DefaultPermission:        true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "describe",
				Description:              "Writes a prompt based on your image.",
				DmPermission:             true,
				Options:                  []DSCommandOption{{Type: 11, Name: "image", Description: "The image to describe", Required: true}},
			},
			Attachments: []ReqCommandAttachments{{
				Id:             "0",
				Filename:       filepath.Base(uploadName),
				UploadFilename: uploadName,
			}},
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Attachments(name string, size int64) (ResAttachments, error) {
	requestBody := ReqAttachments{
		Files: []ReqFile{{
			Filename: name,
			FileSize: size,
			Id:       "1",
		}},
	}
	uploadUrl := fmt.Sprintf(uploadUrlFormat, config.GetConfig().DISCORD_CHANNEL_ID)
	body, err := request(requestBody, uploadUrl)
	var data ResAttachments
	json.Unmarshal(body, &data)
	return data, err
}

func request(params interface{}, url string) ([]byte, error) {
	requestData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
        pp.Println(string(requestData))
        pp.Println("debug")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", config.GetConfig().DISCORD_USER_TOKEN)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
        //req.Header.Set("Cookie", "__dcfduid=7cb3da50041511eea7fad967352beea2; __sdcfduid=7cb3da51041511eea7fad967352beea26bbdf4f8a8d663056ff30d65e83ec6c0109357486cf86b2a6d7175262ddf7b3c; _ga=GA1.1.208433333.1688970170; OptanonConsent=isIABGlobal=false&datestamp=Wed+Dec+06+2023+10%3A01%3A55+GMT%2B0800+(%E4%B8%AD%E5%9B%BD%E6%A0%87%E5%87%86%E6%97%B6%E9%97%B4)&version=6.33.0&hosts=&landingPath=https%3A%2F%2Fdiscord.com%2F&groups=C0001%3A1%2CC0002%3A1%2CC0003%3A1; _ga_Q149DFWHT7=GS1.1.1701828115.3.0.1701828117.0.0.0; __cfruid=8c99b1775cdeaf56378539ad5b7b77fb65a0df92-1712055442; _cfuvid=1lwU0exUaUGnOFJp5QC4qxLzYmrzvZ5Bmd.Ffwzx7j8-1712055442149-0.0.1.1-604800000; _ga_5CWMJQ1S0X=GS1.1.1712055443.1.0.1712055443.0.0.0; cf_clearance=uEVh0VLr_w_Gn4x87DJBKVpwh0DL08GnO3A0PAmoCUw-1712055444-1.0.1.1-.R7pS_R7610OBp8qp98mpqEJndtfyTxL2CYukzT4GbssalHZ_iMFJbRRpGloxOszYzXEVPUD2c9VzZl..Dym5A")
        //req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
        //req.Header.Set("Accept-Language", "us-EN,us;q=0.9")

        client := &http.Client{}
	pp.Println(req)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bod, respErr := ioutil.ReadAll(response.Body)
	fmt.Println("response:", string(bod), respErr, response.Status, url)

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		resp := &ImageResponse{}
		err = json.Unmarshal(bod, &resp)
		if err != nil {
			return bod, err
		}
		err = fmt.Errorf("statusCode: %d, err: %s", response.StatusCode, resp.Message)

		return bod, err
	}

	return bod, respErr
}
