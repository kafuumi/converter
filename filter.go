package converter

import "strings"

//BulletChatFilter
/*弹幕过滤器，返回为True时，才会该弹幕加入到弹幕池中，
参数为指针类型，所以也可以在过滤器中，修改弹幕内容*/
type BulletChatFilter interface {
	filter(node *BulletChatNode) bool
}

// KeyWordFilter 关键字过滤,包含指定关键字的弹幕将被过滤掉
type KeyWordFilter struct {
	Keyword []string
}

func (k *KeyWordFilter) AppendKeyword(keyword string) {
	if k.Keyword == nil {
		k.Keyword = make([]string, 0)
	}
	k.Keyword = append(k.Keyword, keyword)
}

func (k *KeyWordFilter) filter(node *BulletChatNode) bool {
	for _, key := range k.Keyword {
		if strings.Contains(node.Value, key) {
			return false
		}
	}
	return true
}

// TypeConverter 类型过滤
// 符号定义：s：sc, r: 滚动弹幕，t：顶部弹幕，b：底部弹幕，_：过滤掉该类型（只能出现在目标类型）
// 转换语句：原类型 -> 目标类型，一次可以指定多个类型，原类型和目标类型应该相对 “->" 处于对称位置。
// 例如：
// s -> r ：将sc转换成滚动类型，
// stb -> ___ ：过滤掉sc,顶部弹幕，底部弹幕（这里目标类型不能省略为一个 ”_“）
// stb -> tb_：过滤掉sc，顶部弹幕转底部弹幕，底部弹幕转顶部弹幕
type TypeConverter struct {
	//转换表 s r t b _
	table map[BulletChatType]BulletChatType
}

func NewTypeConverter(table string) *TypeConverter {
	mapper := func(c byte) BulletChatType {
		switch c {
		case 's':
			return SupperChat
		case 'r':
			return Roll
		case 't':
			return Top
		case 'b':
			return Bottom
		case '_':
			return 0
		}
		//无效
		return -1
	}
	t := []byte(table)
	l := len(t)
	converter := &TypeConverter{table: make(map[BulletChatType]BulletChatType)}
	for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
		ms, mt := mapper(t[i]), mapper(t[j])
		if mt == -1 || ms == -1 {
			break
		}
		converter.table[ms] = mt
	}
	return converter
}

func (t *TypeConverter) filter(node *BulletChatNode) bool {
	bts := node.Type
	if btt, ok := t.table[bts]; ok {
		//过滤掉该类型
		if btt == 0 {
			return false
		}
		node.Type = btt
	}
	return true
}

// FilterChain 过滤器链，读取弹幕时，通过过滤器链中的过滤器执行过滤操作
type FilterChain struct {
	filters []BulletChatFilter
}

func NewFilterChain() *FilterChain {
	return &FilterChain{filters: make([]BulletChatFilter, 0)}
}

// AddFilter 添加一个过滤器
func (f *FilterChain) AddFilter(filter BulletChatFilter) *FilterChain {
	f.filters = append(f.filters, filter)
	return f
}

func (f *FilterChain) filter(node *BulletChatNode) bool {
	for _, filter := range f.filters {
		ok := filter.filter(node)
		if !ok {
			return false
		}
	}
	return true
}
