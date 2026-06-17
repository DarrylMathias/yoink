package env

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Env struct{
	DBHost string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName string `mapstructure:"DB_NAME"`
	AwsAccessKeyId string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
}

var EnvValue *Env

func NewEnv(file string) error{
	viper.SetConfigName(file)
	viper.SetConfigType("env")
    viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil{
		return fmt.Errorf("fatal error config file: %w", err)
	}

	env := new(Env)
	if err := viper.Unmarshal(env); err != nil{
		return fmt.Errorf("fatal error config file: %w", err)
	}
	EnvValue = env
	
	os.Setenv("AWS_ACCESS_KEY_ID", EnvValue.AwsAccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", EnvValue.AwsSecretAccessKey)
	return nil
}