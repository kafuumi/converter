package converter

// densityQualifier 弹幕密度限定器
type densityQualifier struct {
	//最大弹幕密度
	limit int
	//时间区间，单位毫秒
	timeSlice int
	//当前时间区间的开始时间点
	point int
	//当前时间区间的弹幕密度
	density int
}

func newDensityQualifier(limit, timeSlice int) *densityQualifier {
	//划分成更细的时间区间，减少两个时间区间之间弹幕的间隔距离
	split := 5
	l, t := limit/split, timeSlice/split
	return &densityQualifier{limit: l, timeSlice: t}
}

//根据密度限定，判断是否显示该弹幕
func (d *densityQualifier) check(node *BulletChatNode) bool {
	//0为无限密度，即不对密度做限制
	if d.limit == 0 {
		return true
	}
	//进入到下一个时间区间
	if d.point+d.timeSlice <= node.Time {
		d.point = node.Time
		d.density = 1
		return true
	}
	//达到最大密度
	if d.density >= d.limit {
		return false
	}
	d.density++
	return true
}
