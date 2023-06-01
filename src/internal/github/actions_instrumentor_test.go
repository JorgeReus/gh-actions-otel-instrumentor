package github

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestVerifySignature(t *testing.T) {
	testCases := []struct {
		name        string
		raiseErr    bool
		contents    string
		secretToken string
		signature   string
	}{
		{
			name:        "ValidEvent",
			raiseErr:    false,
			contents:    "TEST_STRING",
			secretToken: "TEST_TOKEN",
			signature:   "sha256=325c6b839863bed5f2ab98978715b70b3756c4d41029ded7e96e3d81b20f751c",
		},
		{
			name:        "InvalidEvent",
			raiseErr:    true,
			contents:    "TEST_STRING",
			secretToken: "TEST_TOKEN",
			signature:   "INVALID_SIGNATURE",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := verifySignature(tc.contents, tc.secretToken, tc.signature)
			if tc.raiseErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

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

func TestInstrument(t *testing.T) {
	testCases := []struct {
		name          string
		raiseErr      string
		event         events.LambdaFunctionURLRequest
		webhookSecret string
	}{
		{
			name:          "ValidEvent",
			raiseErr:      "",
			event:         *loadEvent(t, "test_data/valid_event.json"),
			webhookSecret: "Z86H(VY!mV2%npZ{",
		},
		{
			name:          "FailedEvent",
			raiseErr:      "Only completed events are processed",
			event:         *loadEvent(t, "test_data/failed_event.json"),
			webhookSecret: "Z86H(VY!mV2%npZ{",
		},
		{
			name:          "SignaturesDontMatch",
			raiseErr:      "Signatures don't match",
			event:         *loadEvent(t, "test_data/valid_event.json"),
			webhookSecret: "INVALID_WEBHOOK",
		},
		{
			name:          "FailedConclusion",
			raiseErr:      "",
			event:         *loadEvent(t, "test_data/conclusion_failed_event.json"),
			webhookSecret: "TEST_SECRET",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Instrument(context.Background(), tc.event.Body, tc.event.Headers["x-hub-signature-256"], tc.webhookSecret)
			if tc.raiseErr != "" {
				assert.ErrorContains(t, err, tc.raiseErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("Invalid Webhook Body", func(t *testing.T) {
		_, err := Instrument(context.Background(), "INVALID_BODY", "TEST_HEADER", "WEBHOOK_SECRET")
		assert.ErrorContains(t, err, "Error parsing webhook body")
	})
}
