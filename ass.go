package converter

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrEmptyPoll     = errors.New("assProcessor: empty bullet chats")
	DefaultAssConfig = AssConfig{
		FontSize: 32,
		FontName: "黑体",
		//默认白色，30%透明度
		Color:     0x4cffffff,
		RollSpeed: 15,
		FixTime:   5,
		TimeShift: 0,
		Width:     1920,
		Height:    1080,
	}
)

type AssConfig struct {
	//字体大小，默认38
	FontSize int
	//字体名称，默认SimHei(黑体)
	FontName string
	//字体颜色，这里的颜色是一个全局设置，一条具体的 弹幕可以通过ass代码进行修改
	Color ARGB
	//滚动弹幕通过时间，单位秒
	RollSpeed int
	//固定弹幕显示的时间，单位秒
	FixTime int
	//时间偏移，对所有弹幕的时间进行偏移，单位毫秒
	TimeShift int
	//分辨率的宽、高，默认为1920 x 1080
	Width, Height int
}

type AssProcessor struct {
	AssConfig
	//弹幕池
	pool *BulletChatPool
}

func NewAssProcessor(config AssConfig, pool *BulletChatPool) *AssProcessor {
	return &AssProcessor{
		AssConfig: config,
		pool:      pool,
	}
}

func (a *AssProcessor) Write(writer io.Writer) (err error) {
	//单位转换,将秒转换为毫秒
	a.FixTime *= 1000
	a.RollSpeed *= 1000
	//Ass文件格式详见：https://wenku.baidu.com/view/3b93b33ab307e87100f69634.html
	if err = a.writeScriptInfo(writer); err != nil {
		return
	}
	if err = a.writeStyle(writer); err != nil {
		return
	}
	//Events
	_, err = fmt.Fprintln(writer, "\n[Events]\nFormat: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text")

	bulletChats := a.pool.BulletChat
	if bulletChats == nil {
		return ErrEmptyPoll
	}

	//将整个屏幕划分为多个不同的轨道，每条轨道高度为字体的大小，其保存的内容为弹幕出现的时间
	trackLen := a.Height / a.FontSize
	//滚动弹幕的轨道
	rollTrack := make(RollTrack, trackLen)
	//顶部弹幕的轨道
	topTrack := make(FixedTrack, trackLen)
	//底部弹幕的轨道
	bottomTrack := make(FixedTrack, trackLen)

	for node := bulletChats.Front(); node != nil; node = node.Next() {
		bullet := node.Value.(BulletChatNode)
		//时间轴偏移
		bullet.Time += a.TimeShift

		switch bullet.Type {
		case Roll, SupperChat:
			//滚动的起始坐标,x = 屏幕宽度+(字体大小*字数)即屏幕宽度+弹幕的长度
			bulletLen := a.FontSize * bullet.Length
			//对水平方向，相对于文字内容的中心，对垂直方向，相对于文字内容的底部
			sx, sy := a.Width+(bulletLen>>1), a.FontSize
			track, _ := rollTrack.findTrack(bullet.Time, bulletLen, a.RollSpeed, a.Width)
			sy *= track
			//结束坐标
			ex, ey := -bulletLen>>1, sy
			//开始时间和结束时间 0:00:00.00 格式
			st, et := TimeToString(bullet.Time), TimeToString(bullet.Time+a.RollSpeed)
			var text string
			if bullet.Color.RGBEquals(a.Color) {
				text = fmt.Sprintf("{\\move(%d,%d,%d,%d)}%s", sx, sy, ex, ey, bullet.Value)
			} else {
				text = fmt.Sprintf("{\\move(%d,%d,%d,%d)\\c%s}%s", sx, sy, ex, ey, bullet.Color.AssHex(), bullet.Value)
			}
			_, err = fmt.Fprintf(writer, "Dialogue: 0,%s,%s,Roll,,0000,0000,0000,,%s\n", st, et, text)
		case Top:
			//显示的坐标
			//y轴坐标表示文字内容底部距离顶部的距离
			x, y := a.Width/2, a.FontSize
			track, _ := topTrack.findTrack(bullet.Time, a.FixTime)
			y *= track
			st, et := TimeToString(bullet.Time), TimeToString(bullet.Time+a.FixTime)
			var text string
			if bullet.Color.RGBEquals(a.Color) {
				text = fmt.Sprintf("{\\pos(%d,%d)}%s", x, y, bullet.Value)
			} else {
				text = fmt.Sprintf("{\\pos(%d,%d)\\c%s}%s", x, y, bullet.Color.AssHex(), bullet.Value)
			}
			_, err = fmt.Fprintf(writer, "Dialogue: 0,%s,%s,Top,,0000,0000,0000,,%s\n", st, et, text)
		case Bottom:
			//显示的坐标
			//y轴坐标表示文字内容底部距离顶部的距离
			x, y := a.Width/2, a.Height
			track, _ := bottomTrack.findTrack(bullet.Time, a.FixTime)
			y -= a.FontSize * track
			st, et := TimeToString(bullet.Time), TimeToString(bullet.Time+a.FixTime)
			var text string
			if bullet.Color.RGBEquals(a.Color) {
				text = fmt.Sprintf("{\\pos(%d,%d)}%s", x, y, bullet.Value)
			} else {
				text = fmt.Sprintf("{\\pos(%d,%d)\\c%s}%s", x, y, bullet.Color.AssHex(), bullet.Value)
			}
			_, err = fmt.Fprintf(writer, "Dialogue: 0,%s,%s,Bottom,,0000,0000,0000,,%s\n", st, et, text)
		default:
			continue
		}
		if err != nil {
			return
		}
	}
	return
}

func (a *AssProcessor) writeScriptInfo(writer io.Writer) (err error) {
	//Script Info
	_, err = fmt.Fprintf(writer,
		`[Script Info]
; Roomid: %s Name: %s Title: %s
Title: LiveDanmaku ASS file
ScriptType: v4.00+
PlayResX: %d
PlayResY: %d
Collisions: Normal
WrapStyle: 2
Timer: 100.0000
`, a.pool.RoomId, a.pool.Name, a.pool.Title, a.Width, a.Height)
	return
}

func (a *AssProcessor) writeStyle(writer io.Writer) (err error) {
	//V4+ Styles
	_, err = fmt.Fprintln(writer, "\n[V4+ Styles]\nFormat: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding")
	_, err = fmt.Fprintf(writer, "Style: Roll,%s,%d,%s,&H00FFFFFF,&H00000000,&H1E6A5149,-1,0,0,0,100,100,0,0,1,0,1,8,0,0,0,1\n",
		a.FontName, a.FontSize, a.Color.AssHexWithA())
	_, err = fmt.Fprintf(writer, "Style: Top,%s,%d,%s,&H00FFFFFF,&H00000000,&H1E6A5149,-1,0,0,0,100,100,0,0,1,0,1,8,0,0,0,1\n",
		a.FontName, a.FontSize, a.Color.AssHexWithA())
	_, err = fmt.Fprintf(writer, "Style: Bottom,%s,%d,%s,&H00FFFFFF,&H00000000,&H1E6A5149,-1,0,0,0,100,100,0,0,1,0,1,2,0,0,0,1\n",
		a.FontName, a.FontSize, a.Color.AssHexWithA())
	_, err = fmt.Fprintf(writer, "Style: SC,%s,%d,%s,&H00FFFFFF,&H00000000,&H1E6A5149,-1,0,0,0,100,100,0,0,1,0,1,8,0,0,0,1\n",
		a.FontName, a.FontSize, a.Color.AssHexWithA())
	return
}

/*
将整个屏幕横向划分成多个轨道，每条弹幕处于各自的轨道上
这只是一种逻辑上的轨道，用于弹幕排布
*/

// RollTrack 滚动弹幕的轨道
type RollTrack []struct {
	TimePoint int
	BulletLen int
}

// FixedTrack 固定弹幕的轨道
type FixedTrack []int

/*
参数：弹幕时间点，弹幕长度（单位像素），弹幕速度（以时间表示），屏幕宽度
返回值：轨道号，是否是不会发生碰撞的轨道
*/
func (r *RollTrack) findTrack(timePoint, bulletLen, speed, width int) (int, bool) {
	ra := []struct {
		TimePoint int
		BulletLen int
	}(*r)

	//容忍值，当每条轨道都不满足要求时，选择一条可容忍的轨道，这里选择两条弹幕时间差相差最大的
	tolerate := 0
	max := 0
	for i := 0; i < len(ra); i++ {
		d := ra[i]
		//该轨道未被使用
		if d.TimePoint == 0 {
			ra[i].TimePoint = timePoint
			ra[i].BulletLen = bulletLen
			return i, true
		}
		//判断当前轨道上的这条弹幕是否已经完全出现在屏幕上，即弹幕的最右端是否出屏幕
		t := float32(timePoint - d.TimePoint)
		//这里的结果需要是浮点数，因为v表示每毫米移动的多少像素，如果是整数，会丢失精度
		v := float32(d.BulletLen+width) / float32(speed)
		if v*t < float32(d.BulletLen) {
			//该弹幕未完全出现，不使用该轨道
			continue
		}
		//弹幕的移动速度可以表示为（bulletLen+width)/speed
		//所以，如果当前弹幕比前一条弹幕短，则移动速度小于前一条弹幕，也就不可能碰撞
		if bulletLen <= d.BulletLen {
			ra[i].TimePoint = timePoint
			ra[i].BulletLen = bulletLen
			return i, true
		}
		//时间差的最小值，小于这个值时，两条弹幕会碰撞
		//这里可以将两条弹幕的碰撞问题类比成追击问题，即当A到达终点时，B恰好追上A，求B最晚的出发时间
		threshold := speed - width*speed/(width+bulletLen)
		td := timePoint - d.TimePoint
		if td > threshold {
			ra[i].TimePoint = timePoint
			ra[i].BulletLen = bulletLen
			return i, true
		}
		if td > max {
			max = td
			tolerate = i
		}
	}
	return tolerate, false
}

/*
参数：弹幕的时间，固定弹幕的显示时间
返回值：轨道号，是否是不会发生碰撞的轨道
*/
func (f *FixedTrack) findTrack(timePoint, fixedTime int) (int, bool) {
	fa := []int(*f)
	for i := 0; i < len(fa); i++ {
		//两条弹幕的时间差大于显示时间时，说明该弹幕出现时，上一条弹幕已经消失
		if fa[i] == 0 || (timePoint-fa[i]) > fixedTime {
			fa[i] = timePoint
			return i, true
		}
	}
	//默认使用第一条轨道
	return 0, false
}

// TimeToString 将毫秒值时间转为 0:00:00.00 格式
func TimeToString(t int) string {
	h, m, s, ms := 0, 0, 0, 0
	//毫秒值只有两位，所以需要除以10
	ms = (t % 1000) / 10
	t /= 1000
	s = t % 60
	t /= 60
	m = t % 60
	h = t / 60
	return fmt.Sprintf("%d:%02d:%02d.%02d", h, m, s, ms)
}
