package main

import (
	"go-nacos-setting/nacos"
	"log"
)

func main() {
	nacosSetting := nacos.NacosSetting{
		AppId:         "nacos-demo",
		NacosServerIp: "127.0.0.1",
		ClientPort:    8080,
		ShowLog:       true,
		OnConfigLoad:  OnConfigLoad,
		//Username: "nacos",
		//Password: "nacos",
	}

	nacos.Init(nacosSetting)
	instance, err := nacos.GetInstance("nacos-demo", "default")
	log.Println("instance:", instance)
	log.Println("err:", err)

	select {}
}

func OnConfigLoad(conf nacos.NacosConfig) {
	log.Println("-----OnConfigLoad----------------------------------------------------------------------------------")
	log.Println("conf.Content.length:", len(conf.Content))
	//printlnStrMap(conf.YAML)
}
