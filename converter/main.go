package main

import (
	"fmt"
	"os"

	"github.com/Hami-Lemon/converter"
)

func main() {
	src, err := os.Open("D:\\ProgrameStudy\\converter\\converter\\test\\test.xml")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("文件不存在：%v\n", "")
			return
		} else {
			panic(err)
		}
	}
	chain := converter.NewFilterChain()
	kf := &converter.KeyWordFilter{Keyword: []string{"?", "？"}}
	chain.AddFilter(converter.NewTypeConverter("stb -> rrr")).AddFilter(kf)
	pool := converter.LoadPool(src, chain)
	_ = src.Close()
	if pool == nil {
		fmt.Println("弹幕为空")
		return
	}

	dst, err := os.Create("./test/test.ass")
	if err != nil {
		panic(err)
	}
	assConfig := converter.DefaultAssConfig
	err = pool.Convert(dst, assConfig)
	if err != nil {
		panic(err)
	}
	_ = dst.Close()
}
