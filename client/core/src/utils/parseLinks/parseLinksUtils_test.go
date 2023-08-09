package parseLinks

import "testing"

func Test_getChannelIDByLink(t *testing.T) {
	type args struct {
		link string
	}
	want := "UCyJrhZm9KXrzRub3-wD2zWg"
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Not clear user link",
			args: args{
				link: "https://www.youtube.com/user/TheBrianMaps/featured",
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "Сlear user link",
			args: args{
				link: "https://www.youtube.com/user/TheBrianMaps",
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "Channel link",
			args: args{
				link: "https://www.youtube.com/channel/UCyJrhZm9KXrzRub3-wD2zWg",
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "Video link",
			args: args{
				link: "https://www.youtube.com/watch?v=Kp1hBOhnRjA",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Invalid link",
			args: args{
				link: "sdfsdfsdf",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetChannelIDByLink(tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("getChannelIDByLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getChannelIDByLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
