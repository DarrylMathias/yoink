package env

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Env struct{
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName string `mapstructure:"DB_NAME"`
	DBSSLRootCert string `mapstructure:"DB_SSL_ROOT_CERT"`

	AwsAccessKeyId string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	CrawlerSqsName string `mapstructure:"CRAWLER_SQS_NAME"`
	IndexerSqsName string `mapstructure:"INDEXER_SQS_NAME"`
	S3BucketName string `mapstructure:"S3_BUCKET_NAME"`

	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
	RedisUsername string `mapstructure:"REDIS_USERNAME"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDatabase string `mapstructure:"REDIS_DATABASE"`
	UpstashRedisRestURL string `mapstructure:"UPSTASH_REDIS_REST_URL"`
	UpstashRedisRestToken string `mapstructure:"UPSTASH_REDIS_REST_TOKEN"`

	ResendAPIKey string `mapstructure:"RESEND_API_KEY"`
}

type Config struct{
	Application string `mapstructure:"APPLICATION"`
	ApplicationPort string `mapstructure:"APPLICATION_PORT"`
	FrontEndURL string `mapstructure:"FRONT_END_URL"`
	Workers string `mapstructure:"WORKERS"`
	NoOfSQSMessages string `mapstructure:"NO_OF_SQS_MESSAGES"`
	PostingThreshold string `mapstructure:"POSTING_THRESHOLD"`
}

var EnvValue *Env
var ConfigValue *Config

func NewEnv() error{
	// env
	viper.SetConfigName(".env")
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

	// config
	viper.SetConfigName(".env.config")
	viper.SetConfigType("env")
    viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil{
		return fmt.Errorf("fatal error config file: %w", err)
	}

	config := new(Config)
	if err := viper.Unmarshal(config); err != nil{
		return fmt.Errorf("fatal error config file: %w", err)
	}
	ConfigValue = config
	
	os.Setenv("AWS_ACCESS_KEY_ID", EnvValue.AwsAccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", EnvValue.AwsSecretAccessKey)
	return nil
}