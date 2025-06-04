package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	filepath.Walk("D:\\golang projects\\go_training\\01_basic_training", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		/*这里的path是递归的path并不是一个常量是一个变量随着文件夹层级而改变*/
		fmt.Println(filepath.Join(path, info.Name()))
		return nil
	})
}
