package auth

import (
	"encoding/json"
	"gitlab.com/projectreferral/auth-api/configs"
	"gitlab.com/projectreferral/auth-api/models"
	request "gitlab.com/projectreferral/util/pkg/http_lib"
	"gitlab.com/projectreferral/util/pkg/security"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func IssueToken(expiry time.Duration, audience string, subject string, body io.ReadCloser) security.TokenResponse {
	t := time.Now()
	e := t.Add(expiry * time.Minute)
	var u models.UserResponse

	//assign the claims to our customer model
	token := &security.TokenClaims{
		Issuer:		configs.SERVICE_ID,
		Subject:	subject,
		//treat audience as scope(permissions the token has access to)
		Audience:   audience,
		IssuedAt:   t.Unix(),
		Expiration: e.Unix(),
		NotBefore:  t.Unix(),
		Id:         "NOT_SET",
	}

	if subject != "register" {
		errJson := json.NewDecoder(body).Decode(&u)

		if errJson != nil {
			log.Println("Error parsing data to UserResponse object")
		}
	}

	tr := security.TokenResponse{
		//GenerateToken is our security library
		AccessToken:	security.GenerateToken(token),
		TokenType:		configs.BEARER,
		ExpiresIn:		configs.EXPIRY,
		//No support for refresh tokens as of yet
		RefreshToken: 	"N/A",
		UserData:		u,
	}

	return tr
}

func RecaptchaVerify(w *http.ResponseWriter, token *string, r *models.ReCaptcha){

	form := url.Values{}
	form.Add("response", *token)
	form.Add("secret", configs.RECAPTCHA_SECRET)

	reqVer, errReq := request.Post(configs.RECAPTCHA_VERIFY, []byte(form.Encode()), map[string]string{
		"Content-Type": "application/x-www-form-urlencoded"})

	if errReq != nil {
		log.Printf("Request to [%s] failed\n", configs.RECAPTCHA_VERIFY)
		http.Error(*w, errReq.Error(), 400)
		return
	}

	json.NewDecoder(reqVer.Body).Decode(&r)
}
