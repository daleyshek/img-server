package server

import (
	"log"
	"os"
	"strings"
)

// Config 配置文件
type Config struct {
	Port        string `json:"port"`
	RoutePrefix string `json:"route_prefix"`
	StoragePath string `json:"storage_path"`
	AutoResize  bool   `json:"auto_resize"`
	AuthKey     string `json:"auto_key"`
	HostURL     string `json:"host_url"`
}

// C 配置
var C Config

func init() {
	C = Config{
		Port:        ":7474",
		RoutePrefix: "/ff",
		StoragePath: "./storage/",
		AutoResize:  true,
		AuthKey:     "",
		HostURL:     "",
	}
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
