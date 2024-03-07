package github

import (
	"context"
	"reflect"
	"testing"
)

func TestFindPR(t *testing.T) {
	type args struct {
		ctx        context.Context
		jiraTicket string
	}
	tests := []struct {
		name    string
		args    args
		want    PR
		wantErr bool
	}{
		{
			name: "Debug",
			args: args{
				ctx:        context.Background(),
				jiraTicket: "OM-410",
			},
			want: PR{
				Title:  "OM-410 Add Cache to Groups Service",
				Number: 1507,
				Repo:   "platform",
			},
			wantErr: false,
		},
		{
			name: "Debug",
			args: args{
				ctx:        context.Background(),
				jiraTicket: "OM-337",
			},
			want: PR{
				Title:  "OM-337-allow-omiq-staff-override-in-concat-files",
				Number: 1503,
				Repo:   "platform",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindPR(tt.args.ctx, tt.args.jiraTicket)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindPR() = %v, want %v", got, tt.want)
			}
		})
	}
}
