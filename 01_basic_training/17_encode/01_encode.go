package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// Config /*解析文件*/
type Config struct {
	Database string
	Url      string
	Port     int
	Username string
	Password string
}

func main() {
	/*序列化*/
	config := &Config{
		Database: "oracle",
		Url:      "localhost",
		Port:     3326,
		Username: "root",
		Password: "123456",
	}

	out, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%++v", string(out))
	/*反序列化*/
	openFile, err := os.ReadFile("D:\\golang projects\\go_training\\01_basic_training\\17_encode\\config.yaml")
	if err != nil {
		panic(err)
	}
	var config1 Config
	err = yaml.Unmarshal(openFile, &config1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", config1)
}
