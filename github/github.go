package github

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	ErrNoPR        = fmt.Errorf("no PR found for ticket")
	ErrMultiplePRs = fmt.Errorf("found multiple PRs for ticket")
)

type PR struct {
	Title  string
	Number int
	Repo   string
}

func FindPR(ctx context.Context, jiraTicket string) (PR, error) {
	jiraTicket = strings.ToUpper(jiraTicket)

	apiPRs, err := getOpenPRTitles(ctx, "platform")
	if err != nil {
		return PR{}, fmt.Errorf("could not get open api PRs: %w", err)
	}
	webappPRs, err := getOpenPRTitles(ctx, "webapp")
	if err != nil {
		return PR{}, fmt.Errorf("could not get open webapp PRs: %w", err)
	}

	var foundPR *PR

	for i, pr := range apiPRs {
		ticket := strings.ToUpper(jiraRegex.FindString(pr.Title))
		if ticket != jiraTicket {
			continue
		}
		if foundPR != nil {
			fmt.Println(pr)
			fmt.Println(foundPR)
			return PR{}, ErrMultiplePRs
		}
		foundPR = &apiPRs[i]
	}
	for i, pr := range webappPRs {
		ticket := strings.ToUpper(jiraRegex.FindString(pr.Title))
		if ticket != jiraTicket {
			continue
		}
		if foundPR != nil {
			fmt.Println(pr)
			fmt.Println(foundPR)
			return PR{}, ErrMultiplePRs
		}
		foundPR = &webappPRs[i]
	}

	if foundPR == nil {
		return PR{}, ErrNoPR
	}

	return *foundPR, nil
}

func getOpenPRTitles(ctx context.Context, repo string) ([]PR, error) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}

	tokenService := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tokenClient := oauth2.NewClient(ctx, tokenService)

	client := github.NewClient(tokenClient)

	size := 1000
	prs := make([]PR, size)
	pulls, _, err := client.PullRequests.List(ctx, "omiq-ai", repo, &github.PullRequestListOptions{State: "open"})
	if err != nil {
		return nil, err
	}

	for _, pr := range pulls {
		prs = append(prs, PR{Title: *pr.Title, Number: *pr.Number, Repo: repo})
	}

	return slices.Clip(prs), nil
}

var jiraRegex *regexp.Regexp

func init() {
	jiraRegex = regexp.MustCompile(`(?i)^OM[ -]\d+`)
}
