package main

import "testing"

func Test_validateReleaseCharts(t *testing.T) {
	type args struct {
		apps map[string]*release
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Valid stable Helm chart",
			args: args{apps: map[string]*release{"jenkins": {
				Name:    "jenkins",
				Chart:   "stable/jenkins",
				Version: "0.16.1",
			}}},
			want: true,
		}, {
			name: "Valid local Helm chart",
			args: args{apps: map[string]*release{"example_chart": {
				Name:    "example_chart",
				Chart:   "./test_files/example_chart",
				Version: "0.1.0",
			}}},
			want: true,
		}, {
			name:    "Invalid stable Helm chart",
			args:    args{apps: map[string]*release{"jenkins": {Name: "jenkins", Chart: "stable/invalid123"}}},
			wantErr: true,
		}, {
			name:    "Invalid local Helm chart",
			args:    args{apps: map[string]*release{"test": {Name: "jenkins", Chart: "./test_files"}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateReleaseCharts(tt.args.apps)
			if got != tt.want {
				t.Errorf("validateReleaseCharts() got = %v, want %v", got, tt.want)
			}
			if (err == "") == tt.wantErr {
				t.Errorf("validateReleaseCharts() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
