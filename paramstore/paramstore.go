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

	return parameterToEnv([]*ssm.Parameter{res.Parameter}), nil
}

func (ps *ParameterStore) GetEnvs() (*env.Env, error) {
	input := &ssm.GetParametersByPathInput {
		Path: aws.String(ps.path),
		WithDecryption: aws.Bool(true),
	}
	res, err := ps.svc.GetParametersByPath(input)
	if err != nil {
		return nil, err
	}

	return parameterToEnv(res.Parameters), nil
}

func parameterToEnv(params []*ssm.Parameter) *env.Env {
	envs := env.New()
	for _, param := range params {
		envs.PutEnv(envName(*param.Name), *param.Value)
	}
	return envs
}

func (ps *ParameterStore) parameterName(envName string) string {
	return filepath.Join(ps.path, envName)
}

func envName(paramName string) string {
	return filepath.Base(paramName)
}