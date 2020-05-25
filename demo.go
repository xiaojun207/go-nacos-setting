package main

import (
	"github.com/xiaojun207/go-nacos-setting/nacos"
	"log"
	"time"
)

func main() {
	nacosSetting := nacos.NacosSetting{
		AppId:         "nacos-demo",
		NacosServerIp: "127.0.0.1",
		ClientPort:    8080,
		ShowLog:       false,
		OnConfigLoad:  OnConfigLoad,
		Username:      "nacos",
		Password:      "nacos",
	}

	nacos.Init(nacosSetting)

	go func() {
		for {
			time.Sleep(time.Second)
			instance, err := nacos.GetInstance("nacos-demo", "default")
			if err != nil {
				log.Println("err:", err)
			}
			log.Println("instance:", instance)
		}
	}()

	select {}
}

func OnConfigLoad(conf nacos.NacosConfig) {
	log.Println("-----OnConfigLoad----------------------------------------------------------------------------------")
	log.Println("conf.Content.length:", len(conf.Content))
	//printlnStrMap(conf.YAML)
}
