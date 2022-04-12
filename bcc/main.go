package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Hami-Lemon/converter"
)

const Version = "Bullet Chat Converter (bcc) version: 0.2.1"

var (
	xml      string //待转换的xml文件
	location string //该程序所在的目录，从该目录下读取配置文件
	logger   = log.New(os.Stdout, "[log] ", log.Ldate|log.Lshortfile)
)

func flags() {
	flag.StringVar(&xml, "x", "", "待转换的xml文件，如果该路径是一个目录，则处理目录下的所有xml文件，默认为当前目录")
	flag.Usage = func() {
		fmt.Println(Version)
		fmt.Printf("将B站录播姬录制的XML文件转换成ass文件。\n\n")
		fmt.Println("用法：bcc -x xml")
		flag.PrintDefaults()
	}
	flag.Parse()
	if xml == "" {
		var err error
		xml, err = os.Getwd()
		if err != nil {
			logger.Fatalln(err)
		}
	}
}

func main() {
	flags()
	xmlState, err := os.Stat(xml)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Fatalf("文件：%s不存在\n", xml)
		} else {
			logger.Fatalln(err)
		}
	}

	xmls := make([]string, 0)
	if xmlState.IsDir() {
		if xml[len(xml)-1] != os.PathSeparator {
			xml += string(os.PathSeparator)
		}
		if entries, err := os.ReadDir(xml); err != nil {
			logger.Fatalln(err)
		} else {
			for _, entry := range entries {
				if !entry.IsDir() {
					name := entry.Name()
					if strings.HasSuffix(name, ".xml") {
						xmls = append(xmls, xml+name)
					}
				}
			}
		}
	} else {
		if strings.HasSuffix(xml, ".xml") {
			xmls = append(xmls, xml)
		} else {
			logger.Fatalln("不支持的文件格式。")
		}
	}

	location = filepath.Dir(os.Args[0]) + string(os.PathSeparator) + "setting.json"
	f, err := os.Open(location)
	var setting Setting
	if err != nil {
		if os.IsNotExist(err) {
			setting = DefaultSetting
		} else {
			logger.Fatalln(err)
		}
	} else {
		setting = ReadSetting(f)
	}
	assConfig := setting.GetAssConfig()
	chain := converter.NewFilterChain()
	keywordFilter, typeFilter := setting.GetFilter()
	chain.AddFilter(keywordFilter).AddFilter(typeFilter)
	waitGroup := sync.WaitGroup{}
	var success int32 = 0
	var failed int32 = 0
	for _, file := range xmls {
		waitGroup.Add(1)
		go func(xml string) {
			src, _ := os.Open(xml)
			if src == nil {
				atomic.AddInt32(&failed, 1)
				return
			}
			defer func() {
				_ = src.Close()
				waitGroup.Done()
			}()
			pool := converter.LoadPool(src, chain)
			dotIndex := strings.LastIndex(xml, ".")
			if dotIndex == -1 {
				dotIndex = len(xml)
			}
			dstFile := xml[:dotIndex] + ".ass"
			dst, err := os.Create(dstFile)
			if err != nil {
				atomic.AddInt32(&failed, 1)
				logger.Println(err)
				return
			}
			if err := pool.Convert(dst, assConfig); err == nil {
				fmt.Printf("[ok] %s ==> %s\n", xml, dstFile)
				atomic.AddInt32(&success, 1)
			} else {
				atomic.AddInt32(&failed, 1)
				fmt.Printf("[failed] %s\n", xml)
			}
		}(file)
	}
	waitGroup.Wait()
	fmt.Printf("xml文件总数：%d, 转换成功数：%d 转换失败数：%d\n", len(xmls), success, failed)
}
