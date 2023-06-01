package aws

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func loadEvent(t *testing.T, path string) *events.LambdaFunctionURLRequest {
	validJson, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err.Error())
	}
	var validEvent events.LambdaFunctionURLRequest
	err = json.Unmarshal(validJson, &validEvent)
	if err != nil {
		t.Fatal(err.Error())
	}
	return &validEvent
}

func TestHandler(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		event    events.LambdaFunctionURLRequest
		secrets  map[string]string
	}{
		{
			name:     "DigestHeaderNotPresent",
			response: "x-hub-signature-256 header not found",
			event:    *loadEvent(t, "../../internal/github/test_data/header_not_present_event.json"),
			secrets: map[string]string{
				"gha-instrumentor/webhook-secret": "TEST_SECRET",
			},
		},
		{
			name:     "HandlerSucess",
			response: "",
			event:    *loadEvent(t, "../../internal/github/test_data/valid_event.json"),
			secrets: map[string]string{
		    "gha-instrumentor/webhook-secret": "Z86H(VY!mV2%npZ{",
			},
		},
		{
			name:     "HandlerSignaturesDontMatch",
			response: "Signatures don't match",
			event:    *loadEvent(t, "../../internal/github/test_data/valid_event.json"),
			secrets: map[string]string{
		    "gha-instrumentor/webhook-secret": "TEST_SECRET",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, _ := Handler(context.Background(), tc.event, tc.secrets)
			if tc.response != "" {
				assert.Equal(t, tc.response, resp.Body)
			}
		})
	}
}
