package internal

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.com/projectreferral/auth-api/configs"
	"gitlab.com/projectreferral/auth-api/internal/api/auth"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func SetupEndpoints() {
	_router := mux.NewRouter()

	_router.HandleFunc("/auth", auth.VerifyCredentials).Methods("POST", "OPTIONS")
	_router.HandleFunc("/auth/temp", auth.IssueRegistrationTempToken).Methods("GET")
	//test response that can be used for testing the internal/responses
	_router.HandleFunc("/mock", auth.MockResponse).Methods("GET")

	_router.HandleFunc("/log", displayLog).Methods("GET")

	c := cors.New(cors.Options{
		AllowedMethods: []string{"POST"},
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
		AllowedHeaders: []string{"g-recaptcha-response", "Authorization", "Content-Type","Origin","Accept", "Accept-Encoding", "Accept-Language", "Host", "Connection", "Referer", "Sec-Fetch-Mode", "User-Agent", "Access-Control-Request-Headers", "Access-Control-Request-Method: "},
		OptionsPassthrough: true,
	})

	handler := c.Handler(_router)
	log.Fatal(http.ListenAndServe(configs.PORT,handler))

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