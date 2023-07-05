package handlers

import (
	"context"
	"errors"
	"fmt"
	discord "github.com/bwmarrin/discordgo"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	config "midjourney/initialization"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadToken struct {
	Token    string
	ExpireTs int64
}

var uploadToken *UploadToken = &UploadToken{}

func QiniuUploadImage(attachments []*discord.MessageAttachment) (newAttachments []*discord.MessageAttachment, err error) {

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
		err = doQiniuUploadImage(filePath, key)
		fmt.Println("图片上传结束当前本地时间：", time.Now())
		if err != nil {
			return newAttachments, err
		}

		imagePath := config.GetConfig().QINIU_CDN_HOST + "/" + key

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

func doQiniuUploadImage(localFile string, key string) (err error) {

	defer func() {
		if err != nil {
			fmt.Println(fmt.Printf("doQiniuUploadImage Err: %s", err.Error()))
		}
	}()

	upToken, err := getUploadToken()
	if err != nil {
		return
	}

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := &storage.PutRet{}
	putExtra := storage.PutExtra{}
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		return
	}
	fmt.Println(ret.Key, ret.Hash)
	return
}

func downloadImage(url string, filename string) (filePath string, err error) {

	defer func() {
		if err != nil {
			fmt.Println(fmt.Printf("downloadImage Err: %s", err.Error()))
		}
	}()

	// 发送HTTP GET请求获取图片
	response, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("DownloadFile 发送HTTP请求时发生错误：%s", err.Error())
		return
	}
	defer response.Body.Close()

	dir := config.GetConfig().TMP_DIR

	// 创建保存文件
	filePath = filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("DownloadFile 创建保存文件时发生错误：%s", err.Error())
		return
	}
	defer file.Close()

	// 将HTTP响应的内容写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		err = fmt.Errorf("DownloadFile 写入文件时发生错误：%s", err.Error())
		return
	}

	return
}

func deleteImage(localFile string) {

	// 检查文件是否存在
	_, err := os.Stat(localFile)
	if os.IsNotExist(err) {
		return
	}

	// 删除临时文件
	_ = os.Remove(localFile)
}

func getUploadToken() (upToken string, err error) {

	defer func() {
		if r := recover(); r != nil {
			errStr := fmt.Sprintf("getUploadToken Recover: %#v", r)
			err = errors.New(errStr)
		}
		if err != nil {
			fmt.Println(fmt.Printf("getUploadToken Err: %s", err.Error()))
		}
	}()

	if uploadToken.Token != "" && uploadToken.ExpireTs > (time.Now().Unix()+300) {
		upToken = uploadToken.Token
		return
	}

	putPolicy := storage.PutPolicy{
		Scope:      config.GetConfig().QINIU_BUCKET,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}

	accessKey := config.GetConfig().QINIU_ACCESS_KET
	secretKey := config.GetConfig().QINIU_SECRET_KEY
	mac := qbox.NewMac(accessKey, secretKey)
	upToken = putPolicy.UploadToken(mac)
	if upToken == "" {
		err = errors.New("getUploadToken 获取失败")
		return
	}

	uploadToken = &UploadToken{
		Token:    upToken,
		ExpireTs: time.Now().Unix() + 3300,
	}

	return
}
