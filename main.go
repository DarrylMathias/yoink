package main

import (
	"fmt"
	"yoink/utils/env"
	"yoink/utils/myaws"
)

func main() {
	err := env.NewEnv(".env.local")
	if err != nil {
		panic(fmt.Errorf("error in parsing env --- %s", err.Error()))
	}
	err = myaws.GetConfig()
	if err != nil {
		panic(fmt.Errorf("error in aws config --- %s", err.Error()))
	}
	myaws.GetSQSClient()
}
