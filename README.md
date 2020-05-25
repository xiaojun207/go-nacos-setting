# github.com/xiaojun207/go-nacos-setting
> 本项目，方便大家接入nacos。包括使用，注册中心和配置中心。其中使用了nacos官方sdk。


## 1、使用样例

```
import (
    ...
	"github.com/xiaojun207/go-nacos-setting/nacos"
    ...
)

fun main(){

	nacosSetting := nacos.NacosSetting{
		AppId         :"nacos-demo",
		NacosServerIp :"127.0.0.1",
		ClientPort    :8080,
		ShowLog       :true,
		ConfigType    :"YAML",
		OnConfigLoad  :OnConfigLoad,
	}

	nacos.Init(nacosSetting)

	select {}
}

func OnConfigLoad(conf NacosConfig) {
    log.Printf("---------------------------------------------------------------------------------------")
    version := conf.GetValue("version", "1.0.0")
    log.Println("version:", version)
}

```

## 2、NacosSetting 配置说明
```
NacosSetting{
	AppId           string  // 【必填】，例如：bj-yun-nacos-demo
	NacosServerIp   string  // 【选填】，默认: 127.0.0。1
	NacosServerPort uint64  // 【选填】，默认: 8848
	ClientIp        string  // 【选填】，默认：获取本机IP，可以自己设定
	ClientPort      uint64  // 【选填】，默认：80

	ServiceName     string  // 【选填】，默认：{AppId}
	ClusterName     string  // 【选填】，默认：default

	ConfigDataId    string  // 【选填】，默认：{AppId}
	ConfigGroup     string  // 【选填】，默认：DEFAULT_GROUP
	ConfigType      string  // 【选填】，默认：Properties，支持：JSON、YAML、Properties，所有的配置均以map[string]interface{}回调
	OnConfigLoad 	func(conf map[string]interface{}) // 【选填】，配置更新回调

	ShowLog         bool    // 【选填】，默认：false，因为nacos go sdk设置了log输出到日志文件，不会显示到控制台。当ShowLog=true，日志会显示到控制台
	AESKey  string // 【选填】，默认："", 当不为空时，会检测配置的值，如果是AESEncrypt()，包括起来的，则尝试解密
	DESKey  string // 【选填】，默认："", 当不为空时，会检测配置的值，如果是DESEncrypt()，包括起来的，则尝试解密
}

```
