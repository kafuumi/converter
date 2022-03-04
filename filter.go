package converter

//BulletChatFilter
/*弹幕过滤器，返回为True时，才会该弹幕加入到弹幕池中，
参数为指针类型，所以也可以在过滤器中，修改弹幕内容*/
type BulletChatFilter interface {
	filter(node *BulletChatNode) bool
}

// EmptyFilter 空弹幕过滤器，过滤内容为空的弹幕
type EmptyFilter struct{}

func (e *EmptyFilter) filter(node *BulletChatNode) bool {
	return node.Value != ""
}
