package converter

import (
	"testing"
)

func TestTimeToString(t *testing.T) {
	type args struct {
		t int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"time:100ms to string", args{100}, "0:00:00.10"},
		{"time:10ms to string", args{10}, "0:00:00.01"},
		{"time:1000ms to string", args{1000}, "0:00:01.00"},
		{"time:1100ms to string", args{1100}, "0:00:01.10"},
		{"time:300000ms to string", args{300000}, "0:05:00.00"},
		{"time:3600000ms to string", args{3600000}, "1:00:00.00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeToString(tt.args.t); got != tt.want {
				t.Errorf("TimeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
