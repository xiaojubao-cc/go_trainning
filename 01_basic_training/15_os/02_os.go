package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	sourceFile, _ := os.OpenFile("src.txt", os.O_RDWR|os.O_CREATE, 0666)
	defer sourceFile.Close()
	dstFile, _ := os.OpenFile("dst.txt", os.O_WRONLY|os.O_CREATE, 0666)
	defer dstFile.Close()
	/*指定缓冲区大小进行复制*/
	io.CopyBuffer(dstFile, sourceFile, make([]byte, 1024))
	/*重命名文件或者文件夹*/
	err := os.Rename("D:\\golang projects\\go_training\\dst.txt", "D:\\golang projects\\go_training\\target.txt")
	if err != nil {
		fmt.Errorf("exception:%s", err)
	}
	/*删除单个文件os.remove(),删除当前文件夹下所有的目录和子目录os.removeAll()*/
	/*刷盘os.sync()*/
}
