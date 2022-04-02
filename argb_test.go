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

func TestParseStringARGB(t *testing.T) {
	type args struct {
		alpha float32
		rgb   string
	}
	tests := []struct {
		name    string
		args    args
		want    ARGB
		wantErr bool
	}{
		{"A:0 RGB:0xf1f210 to ARGB", args{alpha: 0, rgb: "0xf1f210"}, ARGB(0x00f1f210), false},
		{"A:0.3 RGB:0xf1f210 to ARGB", args{alpha: 0.3, rgb: "0xf1f210"}, ARGB(0x4cf1f210), false},
		{"A:1 RGB:#f1f210 to ARGB", args{alpha: 1, rgb: "#f1f210"}, ARGB(0xfff1f210), false},
		{"A:1 RGB:f1f210 to ARGB", args{alpha: 1, rgb: "f1f210"}, ARGB(0xfff1f210), false},
		{"A:1 RGB:xxxx to ARGB", args{alpha: 1, rgb: "xxxx"}, ARGB(0), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStringARGB(tt.args.alpha, tt.args.rgb)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStringARGB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseStringARGB() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hexToI(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    int
		wantErr bool
	}{
		{"Hex:000 to Dec", "000", 0, false},
		{"Hex:fff to Dec", "fff", 0xfff, false},
		{"Hex:0ff to Dec", "0ff", 0xff, false},
		{"Hex:f to Dec", "f", 0xf, false},
		{"Hex:f01f2f to Dec", "f01f2f", 0xf01f2f, false},
		{"Hex:0F1A2B to Dec", "0F1A2B", 0x0f1a2b, false},
		{"Hex:HHH to Dec", "HHH", 0, true},
		{"Hex:ffffffffff to Dec", "ffffffffff", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hexToI(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("hexToI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("hexToI() got = %v, want %v", got, tt.want)
			}
		})
	}
}
