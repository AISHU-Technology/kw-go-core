package filter

import (
	"fmt"
	"testing"
)

func TestSqlFilter(t *testing.T) {
	type args struct {
		matchStr string
		exactly  bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test sql 0",
			args: args{
				matchStr: `%'`,
			},
			want: `\%\'`,
		},
		{
			name: "test sql 1",
			args: args{
				matchStr: `update or 1`,
			},
			want: `\update \or 1`,
		},
		{
			name: "test sql 1",
			args: args{
				matchStr: `\test\ttt\\\update`,
			},
			want: `\\\\test\\\\ttt\\\\\\\\\\\\\update`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SqlFilter(tt.args.matchStr, tt.args.exactly)
			if got != tt.want {
				fmt.Printf("SafeString() = %v, want %v", got, tt.want)
			} else {
				fmt.Printf("==success== %v, want %v", got, tt.want)
			}
		})
	}
}
