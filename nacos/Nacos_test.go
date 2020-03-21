package nacos

import (
	"log"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	nacosSetting := NacosSetting{
		AppId:         "nacos-demo",
		NacosServerIp: "127.0.0.1",
		ClientPort:    8080,
	}

	Init(nacosSetting, OnConfigLoad)

	select {}
}

func OnConfigLoad(properties map[string]string) {
	log.Printf("---------------------------------------------------------------------------------------")
	for key, value := range properties {
		log.Println("onload, key:" + key + ", \tvalue:" + value)
	}
}
