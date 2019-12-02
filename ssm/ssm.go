package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"path/filepath"
)

func GetSsmParams(ssmClient ssmiface.SSMAPI, parameterPath string) (params map[string]string, err error) {
	ssmInput := &ssm.GetParametersByPathInput{
		Path:           aws.String(parameterPath),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}

	paramResp, err := ssmClient.GetParametersByPath(ssmInput)
	if err != nil {
		return nil, err
	}

	params = make(map[string]string)
	for _, param := range paramResp.Parameters {
		paramName := filepath.Base(*param.Name)
		params[paramName] = *param.Value
	}

	return params, nil
}
