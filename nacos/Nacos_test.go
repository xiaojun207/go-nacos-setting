package nacos

import (
	"fmt"
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
	fmt.Println("---------------------------------------------------------------------------------------")
	for key, value := range properties {
		fmt.Println("onload, key:" + key + ", \tvalue:" + value)
	}
}
