package handlers

import (
	"fmt"
	discord "github.com/bwmarrin/discordgo"
    "github.com/aliyun/aliyun-oss-go-sdk/oss"
	config "midjourney/initialization"
	"net/url"
	"strings"
	"time"
)

func AliUploadImage(attachments []*discord.MessageAttachment) (newAttachments []*discord.MessageAttachment, err error) {

	newAttachments = make([]*discord.MessageAttachment, 0, len(attachments))

	for _, attachment := range attachments {
		//下载图片到本地
		fmt.Println("图片下载开始当前本地时间：", time.Now())
		filePath, err := downloadImage(attachment.URL, attachment.Filename)
		if err != nil {
			return newAttachments, err
		}
		fmt.Println("图片下载结束当前本地时间：", time.Now())

		parsedURL, err := url.Parse(attachment.URL)
		key := strings.TrimLeft(parsedURL.Path, "/")

		fmt.Println("图片上传开始当前本地时间：", time.Now())
		err = doAliUploadImage(filePath, key)
		fmt.Println("图片上传结束当前本地时间：", time.Now())
		if err != nil {
			return newAttachments, err
		}

		imagePath := config.GetConfig().ALI_OSS_ACCESS_HOST + "/" + key

		newAttachment := &discord.MessageAttachment{
			ID:          attachment.ID,
			URL:         imagePath,
			ProxyURL:    imagePath,
			Filename:    attachment.Filename,
			ContentType: attachment.ContentType,
			Width:       attachment.Width,
			Height:      attachment.Height,
			Size:        attachment.Size,
			Ephemeral:   attachment.Ephemeral,
		}

		newAttachments = append(newAttachments, newAttachment)

		deleteImage(filePath)
	}

	return
}

func doAliUploadImage(localFile string, key string) (err error) {

	defer func() {
		if err != nil {
			fmt.Println(fmt.Printf("doAliUploadImage Err: %s", err.Error()))
		}
	}()

	host := config.GetConfig().ALI_OSS_HOST
	accessKey := config.GetConfig().ALI_OSS_ACCESS_KET
	secretKey := config.GetConfig().ALI_OSS_SECRET_KET
	bucketKey := config.GetConfig().ALI_OSS_BUCKET

	client, err := oss.New(host, accessKey, secretKey)
	if err != nil {
		return
	}

	// 开启Bucket的传输加速状态。
	// Enabled表示传输加速的开关，取值为true表示开启传输加速，取值为false表示关闭传输加速。
	accConfig := oss.TransferAccConfiguration{}
	accConfig.Enabled = true

	err = client.SetBucketTransferAcc(bucketKey, accConfig)
	if err != nil {
		return
	}

	bucket, err := client.Bucket(bucketKey)
	if err != nil {
		return
	}

	err = bucket.PutObjectFromFile(key, localFile)
	if err != nil {
		return
	}

	return
}

