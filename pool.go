package converter

import (
	"bufio"
	"container/list"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

type BulletChatType int8

const (
	//Roll 滚动弹幕
	Roll BulletChatType = iota + 1
	//Top 顶部弹幕
	Top
	//Bottom 底部弹幕
	Bottom
	// SupperChat SC弹幕
	SupperChat
)

//BulletChatNode 代表一条弹幕
type BulletChatNode struct {
	//该弹幕出现的时间，单位毫秒
	Time int
	//弹幕内容
	Value string
	//弹幕内容的字数
	Length int
	//弹幕类型
	Type BulletChatType
	//弹幕颜色,RGB颜色值的十进制形式，例：#FFFFFF 对应为 16777215
	Color ARGB
	//价格,只有弹幕为sc时有效
	Price int
	//显示时间，只有弹幕为sc时有效，单位秒
	ShowTime int
}

type BulletChatPool struct {
	//房间号
	RoomId string
	//主播id
	Name string
	//直播间的标题
	Title string
	//所有弹幕
	BulletChat *list.List
	//录播姬版本号
	Ver string
}

func LoadPool(src io.Reader, chain *FilterChain) *BulletChatPool {
	dom := loadXML(src)
	return parseBulletChat(dom, chain)
}

func parseBulletChat(root *element, chain *FilterChain) *BulletChatPool {
	if root == nil {
		return nil
	}
	pool := new(BulletChatPool)
	children := root.children
	if children == nil {
		return nil
	}
	pool.BulletChat = list.New()
	for node := children.Front(); node != nil; node = node.Next() {
		e := node.Value.(element)
		var bulletChat *BulletChatNode
		switch e.name {
		//录播姬信息
		case "BililiveRecorder":
			pool.Ver = e.attrs["version"]
			//房间信息
		case "BililiveRecorderRecordInfo":
			pool.RoomId = e.attrs["roomid"]
			pool.Name = e.attrs["name"]
			pool.Title = e.attrs["title"]
			//普通弹幕
		case "d":
			bulletChat = plainChat(e)
			//sc
		case "sc":
			bulletChat = superChat(e)
		case "gift":
			//礼物信息
		}

		if bulletChat != nil {
			//当filter为空或过滤器返回true，会直接加入该弹幕，
			if chain == nil || chain.filter(bulletChat) {
				//这里保存的类型时 BulletChatNode 而不是指针
				pool.BulletChat.PushBack(*bulletChat)
			}
		}
	}
	return pool
}

func plainChat(src element) *BulletChatNode {
	//弹幕内容为空，不添加
	if src.value == "" {
		return nil
	}
	//Time初始化为-1，用于测试时判断是否解析成功
	node := &BulletChatNode{
		Time: -1,
	}
	//弹幕内容
	node.Value = src.value
	//弹幕字数
	node.Length = utf8.RuneCountInString(node.Value)
	//从属性p上解析 开始时间点，类型，颜色
	p := strings.Split(src.attrs["p"], ",")
	if len(p) < 8 {
		panic("弹幕格式错误")
	}
	//弹幕类型
	t, _ := strconv.Atoi(p[1])
	switch t {
	case 1:
		//从右往左
		node.Type = Roll
	case 6:
		//从左往右
		node.Type = Roll
	case 5:
		//顶部弹幕
		node.Type = Top
	case 4:
		//底部弹幕
		node.Type = Bottom
	case 7:
	//高级弹幕
	default:
		//其它未知的类型视为滚动弹幕
		node.Type = Roll
	}
	//弹幕颜色
	color, _ := strconv.Atoi(p[3])
	node.Color = ARGB(color)
	//弹幕出现的时间点，10.111 形式，单位秒，这里将其转换为单位为毫米的整数
	timePoint, _ := strconv.ParseFloat(p[0], 32)
	node.Time = int(timePoint * 1000)
	return node
}

func superChat(src element) *BulletChatNode {
	//弹幕内容为空，不添加
	if src.value == "" {
		return nil
	}
	//Time初始化为-1，用于测试时判断是否解析成功
	node := &BulletChatNode{
		Time: -1,
	}
	timePoint, _ := strconv.ParseFloat(src.attrs["ts"], 32)
	price, _ := strconv.Atoi(src.attrs["price"])

	node.Value = src.value
	//弹幕字数
	node.Length = utf8.RuneCountInString(node.Value)
	node.Price = price
	node.Time = int(timePoint * 1000)
	node.Type = SupperChat
	//不同价格的sc有不同的颜色和显示时间
	if price < 50 {
		//显示一分钟
		node.ShowTime = 60
		node.Color = 0x2a60b2
	} else if price < 100 {
		//显示两分钟
		node.ShowTime = 120
		node.Color = 0x427d9e
	} else if price < 500 {
		//显示五分钟
		node.ShowTime = 300
		node.Color = 0xe2b52b
	} else if price < 1000 {
		//显示30分钟
		node.ShowTime = 1800
		node.Color = 0xe09443
	} else if price < 2000 {
		//显示一小时
		node.ShowTime = 3600
		node.Color = 0xe54d4d
	} else {
		//显示两小时
		node.ShowTime = 7200
		node.Color = 0xab1a32
	}
	return node
}

// Convert 转换为ass字幕文件
func (b *BulletChatPool) Convert(dst io.Writer, config AssConfig) error {
	writer, ok := dst.(*bufio.Writer)
	if !ok {
		//512KB的缓冲区
		bufSize := 1024 * 512
		writer = bufio.NewWriterSize(dst, bufSize)
	}
	ap := &assProcessor{AssConfig: config, pool: b}
	err := ap.write(writer)
	if err != nil {
		return err
	}
	return writer.Flush()
}
