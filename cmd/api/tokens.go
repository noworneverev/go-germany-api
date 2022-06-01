package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	user, err := app.models.DB.GetUser(creds.Username)
	if err != nil {
		app.errorJSON(w, errors.New("cannot find user"))
		return
	}

	// hasedPassword := validUser.Password
	hasedPassword := user.Password

	err = bcrypt.CompareHashAndPassword([]byte(hasedPassword), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	var claims jwt.Claims
	// claims.Subject = fmt.Sprint(validUser.ID)
	claims.Subject = fmt.Sprint(user.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "go-germany.com"
	claims.Audiences = []string{"go-germany.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error signing"))
		return
	}

	app.writeJSON(w, http.StatusOK, string(jwtBytes), "response")

}
