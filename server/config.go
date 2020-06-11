package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Config 配置文件
type Config struct {
	Port        string `json:"port"`
	StoragePath string `json:"storage_path"`
	AutoResize  bool   `json:"auto_resize"`
	AuthKey     string `json:"auto_key"`
	HostURL     string `json:"host_url"`
}

// C 配置
var C Config

const configFileName = "img-server.json"

func init() {
	f, err := os.Open(configFileName)
	if err != nil {
		fmt.Println("未找到配置文件，已自动生成")
		f, err = os.OpenFile(configFileName, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("创建配置文件失败", err)
		}
		C.StoragePath = "storage/"
		C.AutoResize = true
		C.Port = "8080"
		js, _ := json.MarshalIndent(C, "", "  ")
		f.Write(js)
		f.Close()
		os.Exit(0)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("配置文件无效", err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &C)
	if err != nil {
		log.Println("配置文件不是有效的json格式", err)
		os.Exit(1)
	}
	C.StoragePath = getPath(C.StoragePath)
	initStoragePath()
}

func initStoragePath() {
	_, err := os.OpenFile(C.StoragePath, os.O_RDONLY, os.ModeDir)
	if err != nil {
		err = os.Mkdir(C.StoragePath, 0755)
		if err != nil {
			log.Println("无法创建保存目录", err)
			os.Exit(1)
		}
	}
}

func getPath(p string) string {
	if strings.HasSuffix(p, "/") {
		return p
	}
	return p + "/"
}
