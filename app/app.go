package app

import (
	"fmt"

	"yoink/utils/env"
	"yoink/utils/myaws"
	mysqs "yoink/utils/myaws/sqs"
)

func App(){
	err := env.NewEnv(".env.local")
	if err != nil {
		panic(fmt.Errorf("error in parsing env --- %s", err.Error()))
	}
	err = myaws.GetConfig()
	if err != nil {
		panic(fmt.Errorf("error in aws config --- %s", err.Error()))
	}
	mysqs.GetSQSClient()
}