package initialization

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DISCORD_USER_TOKEN  string
	DISCORD_BOT_TOKEN   string
	DISCORD_SERVER_ID   string
	DISCORD_CHANNEL_ID  string
	CB_URL              string
	MJ_PORT             string
	QINIU_HOST          string
	QINIU_CDN_HOST      string
	QINIU_ACCESS_KET    string
	QINIU_SECRET_KEY    string
	QINIU_REGION        string
	QINIU_BUCKET        string
	ALI_OSS_HOST        string
	ALI_OSS_ACCESS_KET  string
	ALI_OSS_SECRET_KET  string
	ALI_OSS_BUCKET      string
	ALI_OSS_ACCESS_HOST string
	TMP_DIR             string
}

var config *Config

func LoadConfig(cfg string) *Config {
	viper.SetConfigFile(cfg)
	viper.ReadInConfig()
	viper.AutomaticEnv()
	config = &Config{
		DISCORD_USER_TOKEN:  getViperStringValue("DISCORD_USER_TOKEN"),
		DISCORD_BOT_TOKEN:   getViperStringValue("DISCORD_BOT_TOKEN"),
		DISCORD_SERVER_ID:   getViperStringValue("DISCORD_SERVER_ID"),
		DISCORD_CHANNEL_ID:  getViperStringValue("DISCORD_CHANNEL_ID"),
		CB_URL:              getViperStringValue("CB_URL"),
		MJ_PORT:             getDefaultValue("MJ_PORT", "16007"),
		QINIU_HOST:          getViperStringValue("QINIU_HOST"),
		QINIU_CDN_HOST:      getViperStringValue("QINIU_CDN_HOST"),
		QINIU_ACCESS_KET:    getViperStringValue("QINIU_ACCESS_KET"),
		QINIU_SECRET_KEY:    getViperStringValue("QINIU_SECRET_KEY"),
		QINIU_REGION:        getViperStringValue("QINIU_REGION"),
		QINIU_BUCKET:        getViperStringValue("QINIU_BUCKET"),
		ALI_OSS_HOST:        getViperStringValue("ALI_OSS_HOST"),
		ALI_OSS_ACCESS_KET:  getViperStringValue("ALI_OSS_ACCESS_KET"),
		ALI_OSS_SECRET_KET:  getViperStringValue("ALI_OSS_SECRET_KET"),
		ALI_OSS_BUCKET:      getViperStringValue("ALI_OSS_BUCKET"),
		ALI_OSS_ACCESS_HOST: getViperStringValue("ALI_OSS_ACCESS_HOST"),
		TMP_DIR:             getViperStringValue("TMP_DIR"),
	}
	return config
}

func GetConfig() *Config {
	return config
}

func getViperStringValue(key string) string {
	value := viper.GetString(key)
	if value == "" {
		panic(fmt.Errorf("%s MUST be provided in environment or config.yaml file", key))
	}
	return value
}

func getDefaultValue(key string, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}
