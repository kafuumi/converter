package converter

/*
将整个屏幕横向划分成多个轨道，每条弹幕处于各自的轨道上
这只是一种逻辑上的轨道，用于弹幕排布
*/

type track interface {
	findTrack(node *BulletChatNode, config *AssConfig) (int, bool)
}

type rollTrackNode struct {
	timePoint int
	bulletLen int
}

// rollTrack 滚动弹幕的轨道
type rollTrack []rollTrackNode

// newRollTrack 根据弹幕设置计算出轨道数
func newRollTrack(config *AssConfig) rollTrack {
	//增加弹幕间距可以等价为增加文字高度
	fh := float32(config.Fontsize + config.Spacing)
	h := float32(config.Height) * config.RollRange
	//最后一行弹幕不受间距影响，所以这里将屏幕高度放大
	h += float32(config.Spacing)
	trackNum := int(h / fh)
	t := make(rollTrack, trackNum)
	return t
}

/*
参数：弹幕时间点，弹幕长度（单位像素），弹幕速度（以时间表示），屏幕宽度
返回值：轨道号，是否是不会发生碰撞的轨道
*/
func (r rollTrack) findTrack(node *BulletChatNode, config *AssConfig) (int, bool) {
	ra := []rollTrackNode(r)
	//滚动的起始坐标,x = 屏幕宽度+(字体大小*字数)即屏幕宽度+弹幕的长度
	bulletLen := config.Fontsize * node.Length
	timePoint, speed, width := node.Time, config.RollSpeed, config.Width
	//容忍值，当每条轨道都不满足要求时，选择一条可容忍的轨道，这里选择两条弹幕时间差相差最大的
	tolerate := 0
	max := 0
	for i := 0; i < len(ra); i++ {
		d := ra[i]
		//该轨道未被使用
		if d.timePoint == 0 {
			ra[i].timePoint = timePoint
			ra[i].bulletLen = bulletLen
			return i, true
		}
		//判断当前轨道上的这条弹幕是否已经完全出现在屏幕上，即弹幕的最右端是否出屏幕
		t := float32(timePoint - d.timePoint)
		//这里的结果需要是浮点数，因为v表示每毫米移动的多少像素，如果是整数，会丢失精度
		v := float32(d.bulletLen+width) / float32(speed)
		if v*t < float32(d.bulletLen) {
			//该弹幕未完全出现，不使用该轨道
			continue
		}
		//弹幕的移动速度可以表示为（bulletLen+width)/speed
		//所以，如果当前弹幕比前一条弹幕短，则移动速度小于前一条弹幕，也就不可能碰撞
		if bulletLen <= d.bulletLen {
			ra[i].timePoint = timePoint
			ra[i].bulletLen = bulletLen
			return i, true
		}
		//时间差的最小值，小于这个值时，两条弹幕会碰撞
		//这里可以将两条弹幕的碰撞问题类比成追击问题，即当A到达终点时，B恰好追上A，求B最晚的出发时间
		threshold := speed - width*speed/(width+bulletLen)
		td := timePoint - d.timePoint
		if td > threshold {
			ra[i].timePoint = timePoint
			ra[i].bulletLen = bulletLen
			return i, true
		}
		//允许弹幕重叠时才执行
		if config.Overlay {
			if td > max {
				max = td
				tolerate = i
			}
		}
	}

	return tolerate, config.Overlay
}

// fixedTrack 固定弹幕的轨道
type fixedTrack []int

// newFixedTrack 根据弹幕设置计算出轨道数
func newFixedTrack(config *AssConfig) fixedTrack {
	//增加弹幕间距可以等价为增加文字高度
	fh := float32(config.Fontsize + config.Spacing)
	h := float32(config.Height) * config.FixedRange
	//最后一行弹幕不受间距影响，所以这里将屏幕高度放大
	h += float32(config.Spacing)
	trackNum := int(h / fh)
	t := make(fixedTrack, trackNum)
	return t
}

/*
参数：弹幕的时间，固定弹幕的显示时间
返回值：轨道号，是否是不会发生碰撞的轨道
*/
func (f fixedTrack) findTrack(node *BulletChatNode, config *AssConfig) (int, bool) {
	timePoint, fixedTime := node.Time, config.FixTime
	fa := []int(f)
	for i := 0; i < len(fa); i++ {
		//两条弹幕的时间差大于显示时间时，说明该弹幕出现时，上一条弹幕已经消失
		if fa[i] == 0 || (timePoint-fa[i]) > fixedTime {
			fa[i] = timePoint
			return i, true
		}
	}
	//默认使用第一条轨道
	return 0, config.Overlay
}
