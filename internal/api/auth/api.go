package auth

import (
	"encoding/json"
	"fmt"
	"gitlab.com/projectreferral/auth-api/configs"
	"gitlab.com/projectreferral/auth-api/models"
	request "gitlab.com/projectreferral/util/pkg/http_lib"
	"io/ioutil"
	"log"
	"net/http"
)

//Validates the request as human/robot with recaptcha
//Validates the credentials via a request to the Account-API
//Token is issued as a JSON with an expiry time of 2.5days
//This token will allow the user to access the [/GET,/PATCH,/DELETE] endpoints for the Account-API
func VerifyCredentials(w http.ResponseWriter, req *http.Request) {

	if (*req).Method == "OPTIONS" {
		return
	}

	r := &models.ReCaptcha{}

	s := req.Header.Get("g-recaptcha-response")
	RecaptchaVerify(&w, &s, r)

	if !r.Success {
		log.Printf("ReCaptcha verification failed with [%s]\n", r.Error)
		http.Error(w, "Unable to verify recaptcha", http.StatusUnauthorized)
		return
	}

	//empty body
	if req.ContentLength < 1 {
		log.Println("PayLoad empty")
		http.Error(w, "Empty PayLoad", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Println("Error parsing body")
		http.Error(w, "Error parsing body", http.StatusBadRequest)
		return
	}

	//request to account api to verify credentials
	resp, errPost := request.Post(configs.LOGIN_ENDPOINT, body,
		map[string]string{configs.AUTHORIZATION: req.Header.Get(configs.AUTHORIZATION)})

	if errPost != nil {
		log.Println(errPost.Error())
		http.Error(w, errPost.Error(), 400)
		return
	}

	log.Printf("Response to %s returned %d\n", configs.LOGIN_ENDPOINT, resp.StatusCode)
	if resp.StatusCode != 200 {
		errorBody, errParse := ioutil.ReadAll(resp.Body)

		if errParse != nil {
			log.Printf("Error parsing body from [%s]\n", configs.LOGIN_ENDPOINT)
			http.Error(w, "Error parsing body", http.StatusBadRequest)
			return
		}

		http.Error(w, string(errorBody), resp.StatusCode)
		return
	}

	//subject here is the email
	token := IssueToken(configs.EXPIRY, configs.AUTH_AUTHENTICATED, resp.Header.Get(configs.SUBJECT), resp.Body)

	b, err := json.Marshal(token)
	if err != nil {
		log.Println(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

//A temporary token can be requested for registration
//This token will only allow the user to access the /PUT endpoint for the Account-API
func IssueRegistrationTempToken(w http.ResponseWriter, req *http.Request) {
	token := IssueToken(configs.TEMP_EXPIRY, configs.AUTH_REGISTER, "register", nil)

	b, err := json.Marshal(token)

	if err != nil {
		fmt.Sprintf(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

//Response for testing purposes
func MockResponse(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("OK"))
	w.WriteHeader(http.StatusOK)
}
