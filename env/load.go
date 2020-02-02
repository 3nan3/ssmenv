package env

import (
	"os"
	"io"
	"bytes"
	"strings"
	"github.com/joho/godotenv"
)

func (env *Env) LoadDotenv(fileName string) error {
	reader, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer reader.Close()

	return env.load(reader)
}

func (env *Env) LoadEnvVars(text string) error {
	reader := bytes.NewBufferString(text)

	return env.load(reader)
}

func (env *Env) load(io io.Reader) error {
	loaded, err := godotenv.Parse(io)
	if err != nil {
		return err
	}

	replacer := strings.NewReplacer("\n", "\\n")
	for name, value := range loaded {
		rep := replacer.Replace(value)
		env.PutEnv(name, &rep)
	}
	return nil
}
