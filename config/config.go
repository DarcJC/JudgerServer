package config

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var fileConfig map[string]interface{}

func init() {
    // 读取配置文件
    b, err := ioutil.ReadFile("config.yml")
    if err != nil {
        panic("Error occur when try to read config.yml.")
    }
    yaml.Unmarshal(b, &fileConfig)
}

// GetConfig 获取配置项
func GetConfig(key string) (string, error) {
    value := os.Getenv(key)
    if value != "" {
        return value, nil
    }

    // 不存在该环境变量则从配置文件中读取
    v := fileConfig[key]
    if v != nil {
        if v2, ok := v.(string); ok {
            return v2, nil
        }
    }

    return "", errors.New("Could not found the setting object")
}
