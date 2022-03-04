package converter

import (
	"container/list"
	"encoding/xml"
	"io"
)

type Element struct {
	//标签名
	Name string
	//该标签的值
	Value string
	//该标签中的属性
	Attrs map[string]string
	//该标签中的子标签
	Children *list.List
	//该标签的父标签
	Parent *Element
	//整个xml文档的根标签
	Root *Element
}

func LoadXML(reader io.Reader) (current *Element) {
	decoder := xml.NewDecoder(reader)
	//标记是否是根节点
	isRoot := true
	for {
		t, er := decoder.Token()
		if er != nil {
			if er == io.EOF {
				break
			} else {
				panic(er)
			}
		}

		switch token := t.(type) {
		case xml.StartElement:
			//一个StartElement即为一个XML标签的开始
			//从中获取去属性，标签名
			element := &Element{
				Name:  token.Name.Local,
				Attrs: make(map[string]string),
			}
			//以属性名为键
			for _, attr := range token.Attr {
				element.Attrs[attr.Name.Local] = attr.Value
			}

			if isRoot {
				isRoot = false
				//根节点的Root为其自身
				element.Root = element
			} else {
				element.Root = current.Root
				element.Parent = current
			}
			//更新current为当前标签
			current = element
		case xml.CharData:
			//CharData即为标签中的值
			if current != nil {
				current.Value = string(token.Copy())
			}
		case xml.EndElement:
			//一个标签的结束标签，即"<demo>value</demo>"中的 "</demo>"
			if current.Parent != nil {
				//当前节点处理完成时，将其加入到其父节点的Children链表中
				childrenList := current.Parent.Children
				if childrenList == nil {
					childrenList = list.New()
				}
				//这里的类型是 Element
				childrenList.PushBack(*current)
				current.Parent.Children = childrenList
				//更新current为当前节点的父节点
				current = current.Parent
			}
		case xml.Comment:
		//xml文档中的注释信息
		case xml.Directive:
		case xml.ProcInst:

		}
	}
	return current.Root
}
