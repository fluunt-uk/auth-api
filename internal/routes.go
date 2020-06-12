package internal

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/projectreferral/auth-api/configs"
	"gitlab.com/projectreferral/auth-api/internal/api/auth"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func SetupEndpoints() {
	_router := mux.NewRouter()

	_router.HandleFunc("/auth", auth.VerifyCredentials).Methods("GET")
	_router.HandleFunc("/auth/temp", auth.IssueRegistrationTempToken).Methods("GET")
	//test response that can be used for testing the internal/responses
	_router.HandleFunc("/mock", auth.MockResponse).Methods("GET")

	_router.HandleFunc("/log", displayLog).Methods("GET")

	log.Fatal(http.ListenAndServe(configs.PORT, _router))
}

func displayLog(w http.ResponseWriter, r *http.Request){

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	b, _ := ioutil.ReadFile(path + "/logs/authAPI_log.txt")

	w.Write(b)
}