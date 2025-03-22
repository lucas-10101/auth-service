package main

import (
	"fmt"

	"github.com/lucas-10101/auth-service/api/conf"
)

func main() {
	defer func() {
		fmt.Println(recover())
	}()

	//api.RunServer()

	conf.LoadProperties()

	fmt.Print(conf.ApplicationProperties)
}
