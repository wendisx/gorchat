package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/wendisx/gorchat/internal/constant"
)

type Env map[string]string

func NewEnv(file string) Env {
	env, err := godotenv.Read(file)
	if err != nil {
		log.Fatalf("[init] -- (config/dotenv) %s init failed.", file)
	} else {
		mode := env[constant.SERVER_MODE]
		log.Printf("[init] -- (config/dotenv) mode: %s", mode)
	}
	return env
}
