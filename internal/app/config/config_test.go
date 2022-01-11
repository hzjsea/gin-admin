package config

import (
	"encoding/json"
	"testing"
)

func TestMustLoad(t *testing.T) {
	type args struct {
		fpaths []string
	}
	var tests = []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				fpaths: []string{
					"/Users/alpaca/gitProject/goooooooo/gin-admin/configs/config.toml",
					"/Users/alpaca/gitProject/goooooooo/gin-admin/configs/config2.toml",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.args.fpaths)
			MustLoad(tt.args.fpaths...)
			r, _ := json.MarshalIndent(C, "" , "")
			t.Log(string(r))
		})
	}
}
