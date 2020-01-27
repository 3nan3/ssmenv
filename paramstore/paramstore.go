package paramstore

import (
	"path/filepath"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/3nan3/ssmenv/env"
)

type ParameterStore struct {
	path string	
	sess *session.Session
	svc  ssmiface.SSMAPI
}

func New(path string) *ParameterStore {
	sess := session.New()
	svc := ssm.New(sess)
	return &ParameterStore{path, sess, svc}
}

func (ps *ParameterStore) GetEnv(envName string) (*env.Env, error) {
	input := &ssm.GetParameterInput {
		Name: aws.String(ps.parameterName(envName)),
		WithDecryption: aws.Bool(true),
	}
	res, err := ps.svc.GetParameter(input)
	if err != nil {
		return nil, err
	}

	envs := env.New()
	putParameters(envs, []*ssm.Parameter{res.Parameter})
	return envs, nil
}

func (ps *ParameterStore) GetEnvs() (*env.Env, error) {
	envs := env.New()

	nextToken := ""
	input := &ssm.GetParametersByPathInput {
		Path: aws.String(ps.path),
		WithDecryption: aws.Bool(true),
	}

	for {
		if nextToken != "" {
			input.SetNextToken(nextToken)
		}
		res, err := ps.svc.GetParametersByPath(input)
		if err != nil {
			return nil, err
		}

		putParameters(envs, res.Parameters)
		if res.NextToken == nil {
			break
		}
		nextToken = *res.NextToken
	}
	return envs, nil
}

func (ps *ParameterStore) PutEnvs(envs *env.Env) (*env.Env, error) {
	oldenvs, err := ps.GetEnvs()
	if err != nil {
		return nil, err
	}
	for name, value := range envs.GetEnvs() {
		input := &ssm.PutParameterInput {
			Name: aws.String(ps.parameterName(name)),
			Overwrite: aws.Bool(true),
			Type: aws.String("SecureString"),
			Value: aws.String(value),
		}
		_, err := ps.svc.PutParameter(input)
		if err != nil {
			return nil, err
		}
	}
	return oldenvs, nil
}

func putParameters(envs *env.Env, params []*ssm.Parameter) {
	for _, param := range params {
		envs.PutEnv(envName(*param.Name), *param.Value)
	}
}

func (ps *ParameterStore) parameterName(envName string) string {
	return filepath.Join(ps.path, envName)
}

func envName(paramName string) string {
	return filepath.Base(paramName)
}