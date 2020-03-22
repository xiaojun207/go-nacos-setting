package nacos

import (
	"log"
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

func OnConfigLoad(conf map[string]interface{}) {
	log.Printf("---------------------------------------------------------------------------------------")
	for key, value := range conf {
		log.Println("onload, key:", key, " \tvalue:", value)
	}
}
