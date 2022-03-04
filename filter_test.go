package converter

import "testing"

func TestEmptyFilter_filter(t *testing.T) {
	type args struct {
		node *BulletChatNode
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Empty BulletChat Value", args{&BulletChatNode{}}, false},
		{"No-Empty BulletChat Value", args{&BulletChatNode{Value: "test"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EmptyFilter{}
			if got := e.filter(tt.args.node); got != tt.want {
				t.Errorf("filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
