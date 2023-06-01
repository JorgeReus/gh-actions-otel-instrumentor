package aws

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func NewAWSSecretsProvider(sdkConfig aws.Config) *AWSSecretsProvider {
	return &AWSSecretsProvider{
		client: secretsmanager.NewFromConfig(sdkConfig),
	}
}

type AWSSecretsProvider struct {
	client *secretsmanager.Client
}

func (p *AWSSecretsProvider) GetSecrets(ctx context.Context, secretPaths []string) (map[string]string, error) {
	errBufferedCh := make(chan error, len(secretPaths))
	var wg sync.WaitGroup
	secrets := map[string]string{}
	wg.Add(len(secretPaths))
	for _, k := range secretPaths {
		go func(k string) {
			defer wg.Done()
			secretTokenInput := secretsmanager.GetSecretValueInput{
				SecretId: aws.String(k),
			}
			secretToken, err := p.client.GetSecretValue(ctx, &secretTokenInput)
			if err != nil {
				fmt.Println(err)
				errBufferedCh <- err
				return
			}
			secrets[k] = *secretToken.SecretString
		}(k)
	}
	wg.Wait()
	select {
	case err := <-errBufferedCh:
		return nil, err
	default:
		// Do nothing without blocking
	}
	return secrets, nil
}
