package deployerai_test

import (
	"context"
	"testing"
	"time"

	"github.com/paul-freeman/deployerai"
)

func TestSelectDeploymentTargetGPT3(t *testing.T) {
	testSelectDeploymentTarget(t, deployerai.ModelGPT3)
}

func TestSelectDeploymentTargetGPT4(t *testing.T) {
	testSelectDeploymentTarget(t, deployerai.ModelGPT4)
}

func testSelectDeploymentTarget(t *testing.T, model deployerai.Model) {
	type args struct {
		ctx               context.Context
		deploymentRequest deployerai.Request
	}
	tests := []struct {
		name    string
		args    args
		want    deployerai.Choice
		wantErr bool
	}{
		{
			name: "Pick a deployment target based on last used time",
			args: args{
				ctx: context.Background(),
				deploymentRequest: deployerai.Request{
					MessageFromUser: "Please deploy pr-1462",
					DeploymentTargets: []deployerai.Target{
						{
							Name:                       "deva",
							CurrentImage:               "pr-1377",
							CurrentImageDeploymentTime: time.Now().Add(-time.Hour * 24 * 3),
							LastRestart:                time.Now().Add(-time.Hour * 24 * 3),
							LastUsed:                   time.Now().Add(-time.Hour * 24 * 2),
						},
						{
							Name:                       "devb",
							CurrentImage:               "pr-1377",
							CurrentImageDeploymentTime: time.Now().Add(-time.Hour * 24 * 3),
							LastRestart:                time.Now().Add(-time.Hour * 24 * 3),
							LastUsed:                   time.Now().Add(-time.Hour * 24 * 3),
						},
					},
					AdditionalNotes: "None",
				},
			},
			want:    deployerai.Choice{DeploymentTargetName: "devb", DeploymentImage: "pr-1462"},
			wantErr: false,
		},
		{
			name: "Listen to additional notes",
			args: args{
				ctx: context.Background(),
				deploymentRequest: deployerai.Request{
					MessageFromUser: "Please deploy pr-1462",
					DeploymentTargets: []deployerai.Target{
						{
							Name:                       "deva",
							CurrentImage:               "pr-1377",
							CurrentImageDeploymentTime: time.Now().Add(-time.Hour * 24 * 3),
							LastRestart:                time.Now().Add(-time.Hour * 24 * 3),
							LastUsed:                   time.Now().Add(-time.Hour * 24 * 2),
						},
						{
							Name:                       "devb",
							CurrentImage:               "pr-1377",
							CurrentImageDeploymentTime: time.Now().Add(-time.Hour * 24 * 3),
							LastRestart:                time.Now().Add(-time.Hour * 24 * 3),
							LastUsed:                   time.Now().Add(-time.Hour * 24 * 3),
						},
					},
					AdditionalNotes: "Avoid putting new builds on devb as it's being used for testing a new feature.",
				},
			},
			want:    deployerai.Choice{DeploymentTargetName: "deva", DeploymentImage: "pr-1462"},
			wantErr: false,
		},
		{
			name: "Pick a deployment target based on one being unused",
			args: args{
				ctx: context.Background(),
				deploymentRequest: deployerai.Request{
					MessageFromUser: "Please deploy pr-1462",
					DeploymentTargets: []deployerai.Target{
						{
							Name:                       "deva",
							CurrentImage:               "pr-1377",
							CurrentImageDeploymentTime: time.Now().Add(-time.Hour * 24 * 3),
							LastRestart:                time.Now().Add(-time.Hour * 24 * 3),
							LastUsed:                   time.Now().Add(-time.Hour * 24 * 2),
						},
						{
							Name:                       "devb",
							CurrentImage:               "pr-1377",
							CurrentImageDeploymentTime: time.Now().Add(-time.Hour * 24 * 3),
							LastRestart:                time.Now().Add(-time.Hour * 24 * 3),
							LastUsed:                   time.Now().Add(-time.Hour * 24 * 3),
						},
						{
							Name:                       "devc",
							CurrentImage:               "pr-1377",
							CurrentImageDeploymentTime: time.Now().Add(-time.Hour * 24 * 3),
							LastRestart:                time.Now().Add(-time.Hour * 24 * 3),
							LastUsed:                   time.Time{},
						},
					},
					AdditionalNotes: "None",
				},
			},
			want:    deployerai.Choice{DeploymentTargetName: "devc", DeploymentImage: "pr-1462"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch got := deployerai.ChooseDeploymentTarget(tt.args.ctx, model, tt.args.deploymentRequest)().(type) {
			case deployerai.Error:
				if !tt.wantErr {
					t.Errorf("SelectDeploymentTarget() error = %v, wantErr %v", got, tt.wantErr)
					return
				}
			case deployerai.Choice:
				if !equalChoices(got, tt.want) {
					t.Errorf("SelectDeploymentTarget() = %v, want %v", got, tt.want)
				} else {
					t.Log(got.Message)
				}
			default:
				t.Errorf("SelectDeploymentTarget() got = %T, want %T", got, deployerai.Choice{})
				return
			}
		})
	}
}

func equalChoices(a, b deployerai.Choice) bool {
	return a.DeploymentTargetName == b.DeploymentTargetName && a.DeploymentImage == b.DeploymentImage
}
