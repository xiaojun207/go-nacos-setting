# go-nacos-setting
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
	}
	nacos.Init(nacosSetting, OnConfigLoad)

	select {}
}

func OnConfigLoad(conf map[string]interface{}) {
    log.Printf("---------------------------------------------------------------------------------------")
    for key, value := range conf {
        log.Println("onload, key:" , key , " \tvalue:" , value)
    }
}

```

## 2、NacosSetting 配置说明
```
    NacosSetting{
        AppId           string  // 必填，例如：bj-yun-nacos-demo
        NacosServerIp   string  // 默认: 127.0.0。1
        NacosServerPort uint64  // 默认: 8848
        ClientIp        string  // 默认：获取本机IP，可以自己设定
        ClientPort      uint64  // 默认：80
    
        ServiceName     string  // 默认：{AppId}
        ClusterName     string  // 默认：default
    
        ConfigDataId    string  // 默认：{AppId}
        ConfigGroup     string  // 默认：DEFAULT_GROUP
        ConfigType      string  // 默认：Properties，支持：JSON、YAML、Properties，所有的配置均以map[string]interface{}回调
        ShowLog         bool    // 默认：false，因为nacos go sdk设置了log输出到日志文件，不会显示到控制台。当ShowLog=true，日志会显示到控制台
    }

```