package internal

import (
	"gitlab.com/projectreferral/auth-api/configs"
	"gitlab.com/projectreferral/auth-api/internal/api/auth"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func SetupEndpoints() {
	_router := mux.NewRouter()

	_router.HandleFunc("/auth", auth.VerifyCredentials).Methods("GET")
	_router.HandleFunc("/auth/temp", auth.IssueRegistrationTempToken).Methods("GET")
	//test response that can be used for testing the internal/responses
	_router.HandleFunc("/mock", auth.MockResponse).Methods("GET")

	log.Fatal(http.ListenAndServe(configs.PORT, _router))
}
