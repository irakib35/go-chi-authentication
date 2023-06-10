package controllers

import (
	"encoding/base64"
	"gochitest/database"
	"gochitest/models"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/unrolled/render"
)

var tokenAuth *jwtauth.JWTAuth
var jrender *render.Render

const Secret = "jwt-secret>"

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(Secret), nil)
	jrender = render.New()
}

func Users(w http.ResponseWriter, r *http.Request) {
	/*cookie, err := r.Cookie("jwt")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "1Unauthenticated User", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
	detoken, err := base64.URLEncoding.DecodeString(cookie.Value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := tokenAuth.Decode(string(detoken))

	if err != nil {
		http.Error(w, "4Unauthenticated User", http.StatusBadRequest)
		return
	}
	jrender.JSON(w, 200, token)*/

}

func MakeToken(id string) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user": id})
	return tokenString
}

func Register(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm := r.FormValue("password_confirm")

	if password != confirm {
		http.Error(w, "Password not match", http.StatusInternalServerError)
		return
	}

	result := map[string]string{
		"name":     name,
		"email":    email,
		"password": password,
	}
	/*user := models.Users{
		Name:  name,
		Email: email,
	}

	user.SetPassword(password)*/

	//database.DB.Create(&user)
	//w.Header().Set("Content-Type", "application/json")

	//json.NewEncoder(w).Encode(user)*/

	//json.NewEncoder(w).Encode(result)
	jrender.JSON(w, 200, result)
}

func Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	//npassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	var user models.Users

	database.DB.Where("email = ?", email).First(&user)

	if user.Id == 0 {
		http.Error(w, "1Invalid Credentials", http.StatusInternalServerError)
		return
	}

	if err := user.ComparePassword(password); err != nil {
		http.Error(w, "2Invalid Credentials", http.StatusInternalServerError)
		return
	}

	//jrender.JSON(w, 200, token)

	token := MakeToken(strconv.Itoa(int(user.Id)))

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
		SameSite: http.SameSiteLaxMode,
		// Uncomment below for HTTPS:
		// Secure: true,
		Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value: base64.URLEncoding.EncodeToString([]byte(token)),
		//Value: token,
	})

	jrender.JSON(w, 200, "Successful")

}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(-time.Hour),
		SameSite: http.SameSiteLaxMode,
		// Uncomment below for HTTPS:
		// Secure: true,
		Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value: "",
	})
}
