package converter

const (
	HexTable        = "0123456789ABCDEF"
	MaskA    uint32 = 0xff000000
	MaskR    uint32 = 0x00ff0000
	MaskG    uint32 = 0x0000ff00
	MaskB    uint32 = 0x000000ff
	MaskRGB         = ^MaskA
)

// ARGB 包含RGB和透明度的颜色值
type ARGB uint32

func ParseRGB(rgb uint32) ARGB {
	//透明度位 置0
	rgb &= MaskRGB
	return ARGB(rgb)
}

func (a ARGB) AssHexWithA() string {
	dst := make([]byte, 10)
	//前缀
	dst[0] = '&'
	dst[1] = 'H'
	i := uint32(a)
	A := byte((i & MaskA) >> 24)
	R := byte((i & MaskR) >> 16)
	G := byte((i & MaskG) >> 8)
	B := byte(i & MaskB)

	//A
	dst[2] = HexTable[A>>4]
	dst[3] = HexTable[A&0x0f]
	//B
	dst[4] = HexTable[B>>4]
	dst[5] = HexTable[B&0x0f]
	//G
	dst[6] = HexTable[G>>4]
	dst[7] = HexTable[G&0x0f]
	//R
	dst[8] = HexTable[R>>4]
	dst[9] = HexTable[R&0x0f]
	return string(dst)
}

func (a ARGB) AssHex() string {
	dst := make([]byte, 8)
	//前缀
	dst[0] = '&'
	dst[1] = 'H'
	i := uint32(a)
	R := byte((i & MaskR) >> 16)
	G := byte((i & MaskG) >> 8)
	B := byte(i & MaskB)

	//B
	dst[2] = HexTable[B>>4]
	dst[3] = HexTable[B&0x0f]
	//G
	dst[4] = HexTable[G>>4]
	dst[5] = HexTable[G&0x0f]
	//R
	dst[6] = HexTable[R>>4]
	dst[7] = HexTable[R&0x0f]
	return string(dst)
}

// RGBEquals 比较颜色是否相等，不比较透明度
func (a ARGB) RGBEquals(t ARGB) bool {
	ai := uint32(a) & MaskRGB
	ti := uint32(t) & MaskRGB
	return ai == ti
}
