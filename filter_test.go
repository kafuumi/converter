package converter

import "testing"

func TestTypeConverter_filter(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		args     *BulletChatNode
		want     bool
		nodeType BulletChatType
	}{
		//命名：转换表达式:输入格式===>期待输出格式
		{"s->_:S==>_", "s -> _", &BulletChatNode{Type: SupperChat},
			false, 0},
		{"s->r:s==>r", "s -> r", &BulletChatNode{Type: SupperChat},
			true, Roll},
		{"srtb->srtb:s==>b", "srtb -> srtb", &BulletChatNode{Type: SupperChat},
			true, Bottom},
		{"srtb->srtb:r==>t", "srtb -> srtb", &BulletChatNode{Type: Roll},
			true, Top},
		{"srtb->srtb:t==>r", "srtb -> srtb", &BulletChatNode{Type: Top},
			true, Roll},
		{"srtb->srtb:b==>s", "srtb -> srtb", &BulletChatNode{Type: Bottom},
			true, SupperChat},
		{"srtb->____:b==>_", "srtb -> ____", &BulletChatNode{Type: Bottom},
			false, 0},
		{"s==>r:s==>r", "s ==> r", &BulletChatNode{Type: SupperChat},
			true, Roll},
		{"s==>r:r==>r", "s ==> r", &BulletChatNode{Type: Roll},
			true, Roll},
		//错误格式
		{"_->b:r==>r", "_ -> b", &BulletChatNode{Type: Roll},
			true, Roll},
		{"srt->b:s==>b", "srt -> b", &BulletChatNode{Type: SupperChat},
			true, Bottom},
		{"srt->b:r==>r", "srt -> b", &BulletChatNode{Type: Roll},
			true, Roll},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewTypeConverter(tt.filter)
			if got := f.filter(tt.args); got != tt.want {
				t.Errorf("filter()=%v, want=%v, %v ==> %v",
					got, tt.want, tt.args.Type, tt.nodeType)
			} else if got && tt.args.Type != tt.nodeType {
				t.Errorf("converter to %v, want=%v", tt.args.Type, tt.nodeType)
			}
		})
	}
}
