package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

/*
	使用注意事项
	确保结构体字段首字母大写（公开可见）
	使用 yaml:"tag_name" 标签指定映射关系
	支持嵌套结构和数组类型
	使用 yaml.UnmarshalStrict 可严格校验未知字段
	常见问题排查：
	检查 YAML 缩进是否正确（必须使用空格）
	确认结构体标签与 YAML key 完全匹配
	复杂类型需要使用 yaml.v3 的 inline 标签实现嵌套展开
	使用 omitempty 标签处理可选字段

*/

// ServerConfig /*yaml映射实体*/
type ServerConfig struct {
	Port    int    `yaml:"port"`
	Env     string `yaml:"env"`
	Timeout string `yaml:"timeout,omitempty"`
}

// Config 嵌套结构体示例
type Config struct {
	Server   ServerConfig `yaml:"server"`
	Features []string     `yaml:"features"`
}

func main() {
	config, _ := LoadConfig("D:\\golang projects\\go_training\\01_basic_training\\10_struct\\app.yml")
	fmt.Printf("%+v", config)
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
