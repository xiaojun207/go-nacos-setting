package main

import (
	"fmt"
	"go-nacos-setting/nacos"
)

func main() {
	nacosSetting := nacos.NacosSetting{
		AppId:         "nacos-demo",
		NacosServerIp: "127.0.0.1",
		ClientPort:    8080,
	}

	nacos.Init(nacosSetting, OnConfigLoad)

	select {}
}

func OnConfigLoad(properties map[string]string) {
	fmt.Println("---------------------------------------------------------------------------------------")
	for key, value := range properties {
		fmt.Println("onload, key:" + key + ", \tvalue:" + value)
	}
}
