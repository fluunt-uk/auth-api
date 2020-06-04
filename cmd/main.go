package main

import (
	"gitlab.com/projectreferral/auth-api/configs"
	"gitlab.com/projectreferral/auth-api/internal"
	"log"
)

func main() {
	log.Println("Running on %s", configs.PORT)
	internal.SetupEndpoints()
}
