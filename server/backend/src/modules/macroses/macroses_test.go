package macroses

import (
	"testing"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/server/backend/src/model"

	"github.com/spf13/viper"
)

func TestExecuteUserMacroses(t *testing.T) {
	logger.Init("logs", "", "")
	viper.AddConfigPath("../../")
	viper.ReadInConfig()
	model.Init()

	type args struct {
		key  string
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "error",
			args: args{
				key:  "909115bd1c38454a6ad49f2f077a158add9052fd",
				text: "[-1-0:tags]",
			},
		},
		{
			name: "less",
			args: args{
				key:  "909115bd1c38454a6ad49f2f077a158add9052fd",
				text: "[less-100:tags]",
			},
		},
		{
			name: "tless",
			args: args{
				key:  "909115bd1c38454a6ad49f2f077a158add9052fd",
				text: "[tless-100:tags]",
			},
		},
		{
			name: "count",
			args: args{
				key:  "909115bd1c38454a6ad49f2f077a158add9052fd",
				text: "[3-3:tags]",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExecuteUserMacroses(tt.args.key, tt.args.text); got != tt.want {
				t.Errorf("ExecuteUserMacroses() = %v, want %v", got, tt.want)
			}
		})
	}
}
