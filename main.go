package main

import (
	"fmt"

	"github.com/lucas-10101/auth-service/api/conf"
	"github.com/lucas-10101/auth-service/api/logger"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Recover: %v", err)
		}
	}()

	conf.LoadProperties()
	logger.Setup()
}
