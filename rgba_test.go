package converter

import (
	"reflect"
	"testing"
)

func TestParseRGB(t *testing.T) {
	type args struct {
		rgb uint32
	}
	tests := []struct {
		name string
		args args
		want ARGB
	}{
		{"RGB:0xf1f210 to ARGB", args{0xf1f210}, 0xf1f210},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseRGB(tt.args.rgb); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRGB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestARGB_AssHexWithA(t *testing.T) {
	tests := []struct {
		name string
		a    ARGB
		want string
	}{
		{"ARGB:0xf1f2f3f4 to AssHexWithA", 0xf1f2f3f4, "&HF1F4F3F2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.AssHexWithA(); got != tt.want {
				t.Errorf("AssHexWithA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestARGB_AssHex(t *testing.T) {
	tests := []struct {
		name string
		a    ARGB
		want string
	}{
		{"ARGB:0xf2f3f4 to AssHexWithA", 0xf2f3f4, "&HF4F3F2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.AssHex(); got != tt.want {
				t.Errorf("AssHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestARGB_RGBEquals(t *testing.T) {
	type args struct {
		t ARGB
	}
	tests := []struct {
		name string
		a    ARGB
		args args
		want bool
	}{
		{"#00ff00aa == #00ff00aa", 0x00ff00aa, args{0x00ff00aa}, true},
		{"#bbff00aa == #00ff00aa", 0xbbff00aa, args{0x00ff00aa}, true},
		{"#ffff00aa == #ffff00bb", 0xffff00aa, args{0xffff00bb}, false},
		{"#ffff00aa == #00ff00bb", 0xffff00aa, args{0x00ff00bb}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.RGBEquals(tt.args.t); got != tt.want {
				t.Errorf("RGBEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}
