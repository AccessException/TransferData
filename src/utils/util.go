package utils

import (
	"path"
	"runtime"
	"strings"
	//"github.com/joho/godotenv"
	//"os"
	"log"
	"github.com/spf13/viper"
	"github.com/fsnotify/fsnotify"
	"fmt"
)

var Config ConfigJson

func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")
	currentPath := GetCurrentDirectory()
	currentPath = path.Join(currentPath, "../")
	log.Println("当前路径：",currentPath)
	viper.AddConfigPath(currentPath)   // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	Reload()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		Reload()
	})

}

// 获取当前路径
func GetCurrentDirectory() string {
	_, filename, _, _ := runtime.Caller(1)
	dir := path.Dir(filename)
	return strings.Replace(dir, "\\", "/", -1)
}

func Reload(){
	Config.MongoAddress = viper.GetString("config.mongoAddress")
	Config.MongoDataBaseName = viper.GetString("config.mongoDataBaseName")
	Config.MongoCollectionName = viper.GetString("config.mongoCollectionName")
	Config.TiDBDataBaseName = viper.GetString("config.tiDBDataBaseName")
	Config.TiDBTableName = viper.GetString("config.tiDBTableName")
}


type ConfigJson struct {
	MongoAddress string `"json:mongoAddress"`
	MongoDataBaseName string `"json:mongoDataBaseName"`
	MongoCollectionName string `"json:mongoCollectionName"`
	TiDBDataBaseName string `"json:tiDBDataBaseName"`
	TiDBTableName string `"json:tiDBTableName"`
}

// 返回model
func Creater(model string) interface{}{
	
	return nil
}

