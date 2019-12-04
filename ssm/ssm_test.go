package ssm

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/stretchr/testify/assert"
)

type mockSsmClient struct {
	ssmiface.SSMAPI
}

func (m *mockSsmClient) GetParametersByPath(in *ssm.GetParametersByPathInput) (out *ssm.GetParametersByPathOutput, err error) {
	params := []*ssm.Parameter{
		{
			Name:  aws.String("/alpha/airflow/airflow_test/secrets/CLIENT_ID"),
			Value: aws.String("test-client-id"),
		},
		{
			Name:  aws.String("/alpha/airflow/airflow_test/secrets/CLIENT_SECRET"),
			Value: aws.String("test-client-secret"),
		},
	}

	return &ssm.GetParametersByPathOutput{
		Parameters: params,
	}, nil
}

func TestGetSsmParams(t *testing.T) {

	const expectedClientIdVal = "test-client-id"
	const expectedClientSecretVal = "test-client-secret"
	ssmSvc := &mockSsmClient{}

	params, err := GetSsmParams(ssmSvc, "/alpha/airflow/airflow_test/secrets/")
	assert.Nil(t, err, err)
	if assert.NotNil(t, params) {
		assert.Equal(t, expectedClientIdVal, params["CLIENT_ID"])
		assert.Equal(t, expectedClientSecretVal, params["CLIENT_SECRET"])
	}
}
