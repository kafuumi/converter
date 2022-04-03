package converter

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrEmptyPoll = errors.New("assProcessor: empty bullet chats")
)

// AssConfig ASS字幕配置
type AssConfig struct {
	//字体大小，默认38
	Fontsize int
	//字体名称，默认SimHei(黑体)
	FontName string
	//字体颜色，这里的颜色是一个全局设置，一条具体的 弹幕可以通过ass代码进行修改
	Color ARGB
	//描边颜色
	OutlineColor ARGB
	//阴影颜色
	BackColor ARGB
	//滚动弹幕通过时间，单位秒
	RollSpeed int
	//固定弹幕显示的时间，单位秒
	FixTime int
	//时间偏移，对所有弹幕的时间进行偏移，单位秒
	TimeShift int
	//粗体
	IsBold bool
	//描边
	Outline int
	//阴影
	Shadow int
	//分辨率的宽、高，默认为1920 x 1080
	Width, Height int
	//滚动弹幕的显示范围
	RollRange float32
	//固定弹幕的显示范围
	FixedRange float32
	//弹幕的上下间距，受字体影响
	Spacing int
	//同屏弹幕密度,0为无限制
	Density int
	//是否允许弹幕重叠
	Overlay bool
}

type assProcessor struct {
	AssConfig
	//弹幕池
	pool *BulletChatPool
}

func (a *assProcessor) write(writer io.Writer) (err error) {
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

	//滚动弹幕的轨道
	var rTrack track = newRollTrack(&a.AssConfig)
	//顶部弹幕的轨道
	var tTrack track = newFixedTrack(&a.AssConfig)
	//底部弹幕的轨道
	var bTrack track = newFixedTrack(&a.AssConfig)
	//密度限定器
	qualifier := newDensityQualifier(a.Density, a.RollSpeed)

	for node := bulletChats.Front(); node != nil; node = node.Next() {
		bullet := node.Value.(BulletChatNode)
		//时间轴偏移
		bullet.Time += a.TimeShift * 1000
		if bullet.Time < 0 {
			//偏移负数时，置为0
			bullet.Time = 0
		}
		switch bullet.Type {
		case Roll:
			//只限定滚动弹幕密度
			if !qualifier.check(&bullet) {
				continue
			}
			trackId, ok := rTrack.findTrack(&bullet, &a.AssConfig)
			if !ok {
				continue
			}
			//弹幕间距带来的偏移
			offset := trackId * a.Spacing
			//滚动的起始坐标,x = 屏幕宽度+(字体大小*字数)即屏幕宽度+弹幕的长度
			bulletLen := a.Fontsize * bullet.Length
			//对水平方向，相对于文字内容的中心，对垂直方向，相对于文字内容的底部
			sx, sy := a.Width+(bulletLen>>1), a.Fontsize
			sy *= trackId
			sy += offset
			//结束坐标
			ex, ey := -bulletLen>>1, sy
			//开始时间和结束时间 0:00:00.00 格式
			st, et := TimeToString(bullet.Time), TimeToString(bullet.Time+a.RollSpeed)
			_, err = fmt.Fprintf(writer,
				"Dialogue: 0,%s,%s,Roll,,0000,0000,0000,,{\\move(%d,%d,%d,%d)\\c%s}%s\n",
				st, et, sx, sy, ex, ey, bullet.Color.AssHex(), bullet.Value)
		case Top:
			//显示的坐标
			//y轴坐标表示文字内容底部距离顶部的距离
			x, y := a.Width/2, a.Fontsize
			trackId, ok := tTrack.findTrack(&bullet, &a.AssConfig)
			if !ok {
				continue
			}
			//弹幕间距带来的偏移
			offset := trackId * a.Spacing
			y *= trackId
			y += offset
			st, et := TimeToString(bullet.Time), TimeToString(bullet.Time+a.FixTime)
			_, err = fmt.Fprintf(writer,
				"Dialogue: 0,%s,%s,Top,,0000,0000,0000,,{\\pos(%d,%d)\\c%s}%s\n",
				st, et, x, y, bullet.Color.AssHex(), bullet.Value)
		case Bottom:
			//显示的坐标
			//y轴坐标表示文字内容底部距离顶部的距离
			x, y := a.Width/2, a.Height
			trackId, ok := bTrack.findTrack(&bullet, &a.AssConfig)
			if !ok {
				continue
			}
			//弹幕间距带来的偏移
			offset := trackId * a.Spacing
			y -= a.Fontsize * trackId
			//底部弹幕应该向上偏移
			y -= offset
			st, et := TimeToString(bullet.Time), TimeToString(bullet.Time+a.FixTime)
			_, err = fmt.Fprintf(writer,
				"Dialogue: 0,%s,%s,Bottom,,0000,0000,0000,,{\\pos(%d,%d)\\c%s}%s\n",
				st, et, x, y, bullet.Color.AssHex(), bullet.Value)
		case SupperChat:
			//sc
		default:
			continue
		}
		if err != nil {
			return
		}
	}
	return
}

func (a *assProcessor) writeScriptInfo(writer io.Writer) (err error) {
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

func (a *assProcessor) writeStyle(writer io.Writer) (err error) {
	bti := func(b bool) (i string) {
		if b {
			return "-1"
		} else {
			return "0"
		}
	}
	//V4+ Styles
	_, err = fmt.Fprintln(writer, "\n[V4+ Styles]\nFormat: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding")
	_, err = fmt.Fprintf(writer, "Style: Roll,%s,%d,%s,&H00FFFFFF,%s,%s,%s,0,0,0,100,100,0,0,1,%d,%d,8,0,0,0,1\n",
		a.FontName, a.Fontsize, a.Color.AssHexWithA(), a.OutlineColor.AssHexWithA(), a.BackColor.AssHexWithA(), bti(a.IsBold), a.Outline, a.Shadow)
	_, err = fmt.Fprintf(writer, "Style: Top,%s,%d,%s,&H00FFFFFF,%s,%s,%s,0,0,0,100,100,0,0,1,%d,%d,8,0,0,0,1\n",
		a.FontName, a.Fontsize, a.Color.AssHexWithA(), a.OutlineColor.AssHexWithA(), a.BackColor.AssHexWithA(), bti(a.IsBold), a.Outline, a.Shadow)
	_, err = fmt.Fprintf(writer, "Style: Bottom,%s,%d,%s,&H00FFFFFF,%s,%s,%s,0,0,0,100,100,0,0,1,%d,%d,2,0,0,0,1\n",
		a.FontName, a.Fontsize, a.Color.AssHexWithA(), a.OutlineColor.AssHexWithA(), a.BackColor.AssHexWithA(), bti(a.IsBold), a.Outline, a.Shadow)
	_, err = fmt.Fprintf(writer, "Style: SC,%s,%d,%s,&H00FFFFFF,%s,%s,%s,0,0,0,100,100,0,0,1,%d,%d,8,0,0,0,1\n",
		a.FontName, a.Fontsize, a.Color.AssHexWithA(), a.OutlineColor.AssHexWithA(), a.BackColor.AssHexWithA(), bti(a.IsBold), a.Outline, a.Shadow)
	return
}
