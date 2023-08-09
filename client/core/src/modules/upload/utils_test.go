package upload

import (
	"testing"
)

func Test_cutStringByWords(t *testing.T) {
	type args struct {
		text       string
		maxSymbols int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Gross",
			args: args{
				text:       "dfg\r\nsdfsd\r\nsdfsd\r\nsdsf",
				maxSymbols: 15,
			},
			want: "dfg\r\nsdfsd\r\nsdfsd",
		},
		{
			name: "Less",
			args: args{
				text:       "dfg\r\nsdfsd\r\nsdfsd\r\nsdsf",
				maxSymbols: 20,
			},
			want: "dfg\r\nsdfsd\r\nsdfsd\r\nsdsf",
		},
		{
			name: "Less with bug",
			args: args{
				text:       "1234567890\r\nsdfsd\r\nsdfsd\r\nsdsf",
				maxSymbols: 5,
			},
			want: "12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cutStringByWords(tt.args.text, tt.args.maxSymbols); got != tt.want {
				t.Errorf("cutStringByWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cutTags(t *testing.T) {
	type args struct {
		tags       string
		maxSymbols int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Gross",
			args: args{
				tags:       "123, 456, 789",
				maxSymbols: 10,
			},
			want: "123, 456",
		},
		{
			name: "Less",
			args: args{
				tags:       "123, 456, 789",
				maxSymbols: 100,
			},
			want: "123, 456, 789",
		},
		{
			name: "Real data ru",
			args: args{
				tags:       "Автомагнитола купить недорого интернет, Автомагнитола купить недорого, Автомагнитола купить, Pioneer mvh x580bt купить, mvh x580bt, Pioneer mvh x580bt, Pioneer mvh, Pioneer, Pioneer mvh x580bt купить Автомагнитола купить недорого интернет, купить автомагнитолу пионер саратов, где купить недорогой автомагнитолу, купить автомагнитолу андроид недорого, где купить магнитолу пионер, автомагнитола pioneer mvh x580bt, автомагнитола пионер все модели и цены, автомагнитолы пионер где купить, магнитола пионер купить",
				maxSymbols: 500,
			},
			want: "Автомагнитола купить недорого интернет, Автомагнитола купить недорого, Автомагнитола купить, Pioneer mvh x580bt купить, mvh x580bt, Pioneer mvh x580bt, Pioneer mvh, Pioneer, Pioneer mvh x580bt купить Автомагнитола купить недорого интернет, купить автомагнитолу пионер саратов, где купить недорогой автомагнитолу, купить автомагнитолу андроид недорого, где купить магнитолу пионер, автомагнитола pioneer mvh x580bt, автомагнитола пионер все модели и цены",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cutTags(tt.args.tags, tt.args.maxSymbols); got != tt.want {
				t.Errorf("cutTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
