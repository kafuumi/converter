package converter

import (
	"errors"
	"strings"
)

const (
	hexTable        = "0123456789ABCDEF"
	maskA    uint32 = 0xff000000
	maskR    uint32 = 0x00ff0000
	maskG    uint32 = 0x0000ff00
	maskB    uint32 = 0x000000ff
	maskRGB         = ^maskA
)

var (
	ErrInvalidHex  = errors.New("invalid number")
	ErrHexOverflow = errors.New("hex number can not gather than 0xffffff")
)

// ARGB 32位无符号数，包含RGB和透明度的颜色值，按大端序排列，顺序为A,R,G,B
//透明度ff表示完全透明，00表示完全不透明，例如：完全不透明的红色表示为0x00ff0000
type ARGB uint32

// ParseRGB 通过不带透明度的rgb值解析
func ParseRGB(rgb uint32) ARGB {
	//透明度位 置0
	rgb &= maskRGB
	return ARGB(rgb)
}

// ParseStringARGB 从字符串中解析
func ParseStringARGB(alpha float32, rgb string) (ARGB, error) {
	rgbStr := strings.ToLower(rgb)
	var rgbValue uint
	var err error
	rgbStr = strings.TrimPrefix(rgbStr, "0x")
	rgbStr = strings.TrimPrefix(rgbStr, "#")
	value, err := hexToI(rgbStr)
	rgbValue = uint(value)

	if err != nil {
		return 0, err
	}
	a := uint(alpha * 255.0)
	a <<= 24
	rgbValue |= a
	return ARGB(rgbValue), nil
}

func hexToI(hex string) (int, error) {
	cToI := func(c byte) int {
		v := -1
		switch {
		case c >= '0' && c <= '9':
			//数字
			v = int(c - '0')
		case c >= 'A' && c <= 'F':
			//大写字母
			v = int(c-'A') + 10
		case c >= 'a' && c <= 'f':
			//小写字母
			v = int(c-'a') + 10
		}
		return v
	}
	hexLen := len(hex)
	if hexLen == 0 {
		return 0, nil
	} else if hexLen > 6 {
		return 0, ErrHexOverflow
	}

	dec := 0
	for i := 0; i < hexLen; i++ {
		if num := cToI(hex[i]); num != -1 {
			num <<= (hexLen - 1 - i) * 4
			dec |= num
		} else {
			return 0, ErrInvalidHex
		}
	}
	return dec, nil
}

// AssHexWithA 返回ass文件中的颜色格式，包含透明度信息
func (a ARGB) AssHexWithA() string {
	dst := make([]byte, 10)
	//前缀
	dst[0] = '&'
	dst[1] = 'H'
	i := uint32(a)
	A := byte((i & maskA) >> 24)
	R := byte((i & maskR) >> 16)
	G := byte((i & maskG) >> 8)
	B := byte(i & maskB)

	//A
	dst[2] = hexTable[A>>4]
	dst[3] = hexTable[A&0x0f]
	//B
	dst[4] = hexTable[B>>4]
	dst[5] = hexTable[B&0x0f]
	//G
	dst[6] = hexTable[G>>4]
	dst[7] = hexTable[G&0x0f]
	//R
	dst[8] = hexTable[R>>4]
	dst[9] = hexTable[R&0x0f]
	return string(dst)
}

// AssHex 返回ass文件中颜色格式，不包含透明度信息
func (a ARGB) AssHex() string {
	dst := make([]byte, 8)
	//前缀
	dst[0] = '&'
	dst[1] = 'H'
	i := uint32(a)
	R := byte((i & maskR) >> 16)
	G := byte((i & maskG) >> 8)
	B := byte(i & maskB)

	//B
	dst[2] = hexTable[B>>4]
	dst[3] = hexTable[B&0x0f]
	//G
	dst[4] = hexTable[G>>4]
	dst[5] = hexTable[G&0x0f]
	//R
	dst[6] = hexTable[R>>4]
	dst[7] = hexTable[R&0x0f]
	return string(dst)
}

// RGBEquals 比较颜色是否相等，不比较透明度
func (a ARGB) RGBEquals(t ARGB) bool {
	ai := uint32(a) & maskRGB
	ti := uint32(t) & maskRGB
	return ai == ti
}
