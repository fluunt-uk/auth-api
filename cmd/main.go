package main

import (
	"fmt"
	"gitlab.com/projectreferral/auth-api/configs"
	"gitlab.com/projectreferral/auth-api/internal"
	"log"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	f, err := os.OpenFile(path + "/logs/authAPI_log.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Printf("Running on %s\n", configs.PORT)
	internal.SetupEndpoints()
}
