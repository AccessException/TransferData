package main

import (
	"os"
	"bufio"
	"strings"
	"log"
	"io"
)

func main(){

	file ,err := os.Open("/Users/apple/go/TransferData/src/entity/user.go")
	if err != nil{
		log.Println("打开文件失败")
	}
	reader := bufio.NewReader(file)
	os.Remove("/Users/apple/go/TransferData/src/entity/model.go")
	f , _:= os.Create("/Users/apple/go/TransferData/src/entity/model.go")
	for{
		// 按行读取文件
		content, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := string(content) // 转换成字符串
		if strings.Contains(line,"User"){
			line = strings.Replace(line,"User","Model",-1)
		}
		f.WriteString(line+"\n")
	}
	f.Close()
	file.Close()

}