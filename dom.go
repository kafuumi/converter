package converter

import (
	"container/list"
	"encoding/xml"
	"io"
)

type element struct {
	//标签名
	name string
	//该标签的值
	value string
	//该标签中的属性
	attrs map[string]string
	//该标签中的子标签
	children *list.List
	//该标签的父标签
	parent *element
	//整个xml文档的根标签
	root *element
}

//解析xml文件
func loadXML(reader io.Reader) (current *element) {
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
			e := &element{
				name:  token.Name.Local,
				attrs: make(map[string]string),
			}
			//以属性名为键
			for _, attr := range token.Attr {
				e.attrs[attr.Name.Local] = attr.Value
			}

			if isRoot {
				isRoot = false
				//根节点的Root为其自身
				e.root = e
			} else {
				e.root = current.root
				e.parent = current
			}
			//更新current为当前标签
			current = e
		case xml.CharData:
			//CharData即为标签中的值
			if current != nil {
				current.value = string(token.Copy())
			}
		case xml.EndElement:
			//一个标签的结束标签，即"<demo>value</demo>"中的 "</demo>"
			if current.parent != nil {
				//当前节点处理完成时，将其加入到其父节点的Children链表中
				childrenList := current.parent.children
				if childrenList == nil {
					childrenList = list.New()
				}
				//这里的类型是 element
				childrenList.PushBack(*current)
				current.parent.children = childrenList
				//更新current为当前节点的父节点
				current = current.parent
			}
		case xml.Comment:
		//xml文档中的注释信息
		case xml.Directive:
		case xml.ProcInst:

		}
	}
	return current.root
}
