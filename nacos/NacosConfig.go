package nacos

import (
	"github.com/xiaojun207/go-base-utils/utils"
	"strings"
)

type NacosConfig struct {
	DataId     string `param:"dataId"`
	Group      string `param:"group"`
	Content    string `param:"content"`
	Namespace  string `param:"namespace"`
	ConfigType string `param:"configType"`
	Properties map[string]string
	JSON       map[string]interface{}
	YAML       map[string]interface{}
	OnChange   func(namespace, group, dataId, data string)
	AESKey     string
	DESKey     string
}

func (e *NacosConfig) LoadData(content string) {
	e.Content = content
	// JSON、YAML、Properties
	if e.ConfigType == "Properties" {
		e.Properties = Properties(e.Content)
	} else if e.ConfigType == "YAML" {
		e.YAML = Yaml(e.Content)
	} else if e.ConfigType == "JSON" {
		e.JSON = make(map[string]interface{})
		utils.JsonToMap(e.Content, e.JSON)
	} else {
		e.Properties = Properties(e.Content)
	}
}

func (e *NacosConfig) GetValue(key string) string {
	value := e.Properties[key]
	value = e.DESDecrypt(value)
	value = e.AESDecrypt(value)
	return value
}

func (e *NacosConfig) GetString(key, defalueValue string) string {
	value := e.GetValue(key)
	if value == "" {
		value = defalueValue
	}
	return value
}

func (e *NacosConfig) GetFloat64(key string, defalueValue float64) float64 {
	str := e.GetValue(key)
	return utils.StrToFloat64Def(str, defalueValue)
}

func (e *NacosConfig) GetBool(key string, defalueValue bool) bool {
	str := e.GetValue(key)
	return utils.StrToBoolDef(str, defalueValue)
}

func (e *NacosConfig) GetInt(key string, defalueValue int) int {
	str := e.GetValue(key)
	return utils.StrToIntDef(str, defalueValue)
}

func (e *NacosConfig) GetInt64(key string, defalueValue int64) int64 {
	str := e.GetValue(key)
	return utils.StrToInt64Def(str, defalueValue)
}

func (e *NacosConfig) DESDecrypt(value string) string {
	if e.DESKey != "" && strings.HasPrefix(value, "DESEncrypt(") && strings.HasSuffix(value, ")") {
		tmp := strings.TrimPrefix(value, "DESEncrypt(")
		tmp = strings.TrimSuffix(tmp, ")")
		value = utils.DESDecrypt(tmp, e.DESKey)
	}
	return value
}

func (e *NacosConfig) AESDecrypt(value string) string {
	if e.AESKey != "" && strings.HasPrefix(value, "AESEncrypt(") && strings.HasSuffix(value, ")") {
		tmp := strings.TrimPrefix(value, "DESEncrypt(")
		tmp = strings.TrimSuffix(tmp, ")")
		value = utils.DESDecrypt(tmp, e.AESKey)
	}
	return value
}
