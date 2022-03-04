package converter

type HotMap struct {
	//每一秒的弹幕数统计,索引为秒数
	Map []int
	//视频的宽度
	Width int
	//视频时长，单位秒
	Duration int
	//考虑到网络延迟，以及发出弹幕所需要的时间，设定一个偏移量
	Shift int
}
