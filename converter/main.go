package main

import (
	"bufio"
	"fmt"
	"github.com/Hami-Lemon/converter"
	"os"
)

func main() {
	src, err := os.Open("./test/test.xml")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("文件不存在：%v\n", "")
			return
		} else {
			panic(err)
		}
	}

	domRoot := converter.LoadXML(src)
	_ = src.Close()
	if domRoot == nil {
		fmt.Println("解析XML文件失败")
		return
	}
	pool := converter.ParseBulletChat(domRoot, &converter.EmptyFilter{})
	if pool == nil {
		fmt.Println("弹幕为空")
		return
	}

	dst, err := os.Create("./test/test.ass")
	if err != nil {
		panic(err)
	}
	//512KB的缓冲区
	bufSize := 1024 * 512
	writer := bufio.NewWriterSize(dst, bufSize)
	assConfig := converter.DefaultAssConfig
	assProcessor := converter.NewAssProcessor(assConfig, pool)
	err = assProcessor.Write(writer)
	if err != nil {
		panic(err)
	}
	_ = writer.Flush()
	_ = dst.Close()
}
