package deployerai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paul-freeman/deployerai/openai"
)

type Request struct {
	MessageFromUser   string   `json:"message_from_developer"`
	DeploymentTargets []Target `json:"possible_deployment_targets"`
	AdditionalNotes   string   `json:"additional_notes"`
}

type Target struct {
	Name                       string    `json:"deployment_target_name"`
	CurrentImage               string    `json:"currently_deployed_image"`
	CurrentImageDeploymentTime time.Time `json:"current_image_deployment_time"`
	LastRestart                time.Time `json:"last_restart_time"`
	LastUsed                   time.Time `json:"last_used_time"`
}

type Choice struct {
	DeploymentTargetName string `json:"deployment_target_name"`
	DeploymentImage      string `json:"deployment_image"`
	Message              string `json:"message"`
}

type Error struct {
	Message string
}

type Model string

const (
	ModelGPT4 = "gpt-4-0125-preview"
	ModelGPT3 = "gpt-3.5-turbo-0125"
)

func ChooseDeploymentTarget(ctx context.Context, model Model, req Request) tea.Cmd {
	return func() tea.Msg {
		choice, err := chooseDeploymentTarget(ctx, model, req)
		if err != nil {
			return Error{Message: err.Error()}
		}
		return choice
	}
}

func chooseDeploymentTarget(ctx context.Context, model Model, req Request) (Choice, error) {
	const url = "https://api.openai.com/v1/chat/completions"

	userContent, err := makeUserContent(req)
	if err != nil {
		return Choice{}, fmt.Errorf("could not make user content: %v", err)
	}

	// Create the payload.
	payload := openai.Payload{
		Model: string(model),
		Messages: []openai.ReqMessage{
			{
				Role:    "system",
				Content: systemContent,
			},
			{
				Role:    "user",
				Content: userContent,
			},
		},
		ResponseFormat: struct {
			Type string `json:"type"`
		}{
			Type: "json_object",
		},
	}

	// Marshal the payload into JSON.
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return Choice{}, fmt.Errorf("could not marshal OpenAI payload: %v", err)
	}

	// Set up the HTTP request.
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return Choice{}, fmt.Errorf("could not create OpenAI request: %v", err)
	}

	// Add headers.
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	// Create a new HTTP client and send the request.
	client := &http.Client{}
	resp, err := client.Do(httpRequest)
	if err != nil {
		return Choice{}, fmt.Errorf("could not send OpenAI request: %v", err)
	}

	// Read the response body.
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Choice{}, fmt.Errorf("OpenAI request failed: %s", body)
	}

	var response openai.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return Choice{}, fmt.Errorf("could not unmarshal OpenAI response: %v", err)
	}

	if len(response.Choices) == 0 {
		return Choice{}, fmt.Errorf("no choices in OpenAI response")
	}

	var choice Choice
	content := response.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &choice); err != nil {
		return Choice{}, fmt.Errorf("could not unmarshal deployment response: %v", err)
	}

	return choice, nil
}

const systemContent = `
You are a devops engineer at OMIQ. You are responsible for choosing which
development service (or "deployment target") a test build should be deployed
to.

You will be given some JSON data containing a deployment request. The request
will contain a list of deployment targets, each with a name and several time
metrics. There may also be notes from the developer. Your goal is to select
the deployment target that is most appropriate for the given deployment.

If a deployment target has not been used before, or if it has been a long
time since it was last used, it is more likely to be a good choice for the
deployment. If a deployment target has been used recently, it is less likely
to be a good choice for the deployment.

Your response should be given in JSON format with a field named "deployment_target_name"
containing the name of the chosen deployment target and a field named "deployment_image"
containing the name of the image to be deployed. You should also include a field named "message"
containing a message for the developer. It should explain why you chose the deployment target
and any other relevant information.
`

func makeUserContent(req Request) (string, error) {
	reqJson, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("could not marshal deployment request: %v", err)
	}
	return fmt.Sprintf("I have a deployment request for you. Here is the JSON data:\n\n%s", string(reqJson)), nil
}
