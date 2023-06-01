package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/stretchr/testify/assert"
  "golang.org/x/exp/maps"
)

func enterSecretsManagerTest() (*testtools.AwsmStubber, *AWSSecretsProvider) {
	stubber := testtools.NewStubber()
	provider := NewAWSSecretsProvider(*stubber.SdkConfig)
	return stubber, provider
}

func StubGetSecret(input *secretsmanager.GetSecretValueInput, ouput *secretsmanager.GetSecretValueOutput, raiseErr *testtools.StubError) testtools.Stub {
	return testtools.Stub{
		OperationName: "GetSecretValue",
		Input:         input,
		Output:        ouput,
		Error:         raiseErr,
	}
}

func TestUnitGetSecrets(t *testing.T) {
	testCases := []struct {
		name     string
		raiseErr *testtools.StubError
		secretPairs    map[string]string
	}{
		{
			name:     "NoErrors",
			raiseErr: nil,
			secretPairs:    map[string]string{
        "TEST_PATH": "TEST_SECRET",
      },
		},
		{
			name:     "RaiseError",
			raiseErr: &testtools.StubError{Err: errors.New("TestError")},
			secretPairs:    map[string]string{
        "TEST_PATH": "",
      },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stubber, secretsProvider := enterSecretsManagerTest()

			for path, value := range tc.secretPairs {
				stubber.Add(StubGetSecret(&secretsmanager.GetSecretValueInput{
					SecretId: aws.String(path),
				},
					&secretsmanager.GetSecretValueOutput{
						SecretString: aws.String(value),
					}, tc.raiseErr))
			}

      secrets, err := secretsProvider.GetSecrets(context.Background(), maps.Keys(tc.secretPairs))
			if tc.raiseErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			for k, v := range tc.secretPairs {
        assert.Equal(t, secrets[k], v)
      }

			testtools.ExitTest(stubber, t)
		})
	}
}
