package main

import (
	"github.com/xiaojun207/go-nacos-setting/nacos"
	"log"
	"time"
)

func main() {
	nacosSetting := nacos.NacosSetting{
		AppId:         "go-nacos-setting",
		NacosServerIp: "127.0.0.1",
		ClientPort:    8848,
		OnConfigLoad:  OnConfigLoad,
		Username:      "nacos",
		Password:      "nacos",
		DESKey:        "dsfsw422",
	}

	nacosSetting = *nacos.Init(nacosSetting)

	go func() {
		for {
			time.Sleep(time.Second)
			serviceAddress, err := nacosSetting.GetServiceAddress("test-service1")
			if err != nil {
				log.Println("err:", err)
			}
			log.Println("instance:", serviceAddress)

			time.Sleep(time.Second)
			serviceAddress2, err := nacosSetting.GetServiceAddress("test-service2")
			if err != nil {
				log.Println("err2:", err)
			}
			log.Println("instance2:", serviceAddress2)
		}
	}()

	select {}
}

func OnConfigLoad(conf nacos.NacosConfig) {
	log.Println("-----OnConfigLoad----------------------------------------------------------------------------------")
	log.Println(conf.Content)
	log.Println("key", conf.GetValue("key")) // key=DESEncrypt(PsvwSazUUxQ=)
	//log.Println("conf.Content.length:", len(conf.Content))
	//printlnStrMap(conf.YAML)
}
