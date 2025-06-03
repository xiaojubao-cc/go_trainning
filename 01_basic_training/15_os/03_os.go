package main

import (
	"fmt"
	"os"
	"path/filepath"
)

/*文件夹操作*/
func main() {
	/*创建文件夹*/
	//os.Mkdir("newFolder", os.ModePerm)
	/*创建指定父路径的文件夹*/
	//os.MkdirAll("D:\\golang projects\\go_training\\01_basic_training\\16_os", os.ModePerm)
	/*读取文件夹路径*/
	filepath.Walk("D:\\golang projects\\go_training\\", func(path string, info os.FileInfo, err error) error {
		fmt.Printf("current path: %s\n", path)
		fmt.Printf("current file info isDir:%v\n", info.IsDir())
		fmt.Printf("current file info name:%v\n", info.Name())
		fmt.Printf("current file info size:%v\n", info.Size())
		fmt.Printf("current file info mode:%v\n", info.Mode())
		fmt.Printf("current file info modTime:%v\n", info.ModTime())
		fmt.Printf("current file info sys:%v\n", info.Sys())
		return nil
	})
}
