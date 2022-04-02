package converter

import "fmt"

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
