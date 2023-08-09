package serverAPI

import (
	"reflect"
	"testing"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/imroc/req"
	"github.com/spf13/viper"
)

func TestExecuteUserMacroses(t *testing.T) {
	req.SetProxyUrl("http://127.0.0.1:8888")
	type args struct {
		text string
	}

	tests := []struct {
		name       string
		args       args
		wantResult string
		wantErr    bool
		key        string
	}{
		{
			args: args{
				text: "test text: [1-1:test]",
			},
			wantResult: "test text: macrostext",
			wantErr:    false,
			key:        "8a6b25255bcd58970c3c5f59b3e0e641e6c422a2",
		},
		{
			args: args{
				text: "test text: [1-1:test]",
			},
			wantResult: "test text: [1-1:test]",
			wantErr:    true,
			key:        "8a6b25255bcd58970c3c5f59b3e0e641e6c422",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("api_key", tt.key)
			gotResult, err := ExecuteUserMacroses(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteUserMacroses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("ExecuteUserMacroses() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestGetMacroses(t *testing.T) {
	req.SetProxyUrl("http://127.0.0.1:8888")
	tests := []struct {
		name       string
		wantResult []moneytubemodel.Macros
		wantErr    bool
		key        string
	}{
		{
			wantResult: []moneytubemodel.Macros{
				moneytubemodel.Macros{
					Name: "test",
					Data: []string{"macrostext"},
				},
			},
			wantErr: false,
			key:     "8a6b25255bcd58970c3c5f59b3e0e641e6c422a2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("api_key", tt.key)
			gotResult, err := GetMacroses()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMacroses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("GetMacroses() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestSaveMacros(t *testing.T) {
	req.SetProxyUrl("http://127.0.0.1:8888")
	viper.Set("api_key", "8a6b25255bcd58970c3c5f59b3e0e641e6c422a2")
	type args struct {
		macros moneytubemodel.Macros
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				macros: moneytubemodel.Macros{
					Name: "test_for_del",
					Data: []string{"12345"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveMacros(tt.args.macros); (err != nil) != tt.wantErr {
				t.Errorf("SaveMacros() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteMacros(t *testing.T) {
	req.SetProxyUrl("http://127.0.0.1:8888")
	viper.Set("api_key", "8a6b25255bcd58970c3c5f59b3e0e641e6c422a2")
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				name: "test_for_del",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteMacros(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteMacros() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetMacros(t *testing.T) {
	req.SetProxyUrl("http://127.0.0.1:8888")
	viper.Set("api_key", "8a6b25255bcd58970c3c5f59b3e0e641e6c422a2")
	type args struct {
		name string
	}
	tests := []struct {
		name       string
		args       args
		wantResult moneytubemodel.Macros
		wantErr    bool
	}{
		{
			args: args{
				name: "test",
			},
			wantResult: moneytubemodel.Macros{
				Name: "test",
				Data: []string{"macrostext"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := GetMacros(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMacros() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("GetMacros() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
