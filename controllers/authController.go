package controllers

import (
	"gochitest/database"
	"gochitest/middlewares"
	"gochitest/models"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/unrolled/render"
)

// var tokenAuth *jwtauth.JWTAuth
var jrender *render.Render

const Secret = "jwt-secret>"

func init() {
	jrender = render.New()
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world Go1"))
}

func Users(w http.ResponseWriter, r *http.Request) {
	id, _ := middlewares.GetUserId(w, r)

	var user models.Users
	database.DB.Where("id = ?", id).First(&user)

	jrender.JSON(w, 200, user)

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

	/*result := map[string]string{
		"name":     name,
		"email":    email,
		"password": password,
	}*/
	user := models.Users{
		Name:  name,
		Email: email,
	}

	user.SetPassword(password)
	database.DB.Create(&user)
	jrender.JSON(w, 200, "Successful")
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

	payload := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(Secret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
		SameSite: http.SameSiteLaxMode,
		// Uncomment below for HTTPS:
		// Secure: true,
		Name: "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		//Value: base64.URLEncoding.EncodeToString([]byte(token)),
		Value: token,
	})

	jrender.JSON(w, 200, "Successful")

}

func Logout(w http.ResponseWriter, r *http.Request) {
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

func UpdateInfo(w http.ResponseWriter, r *http.Request) {
	id, _ := middlewares.GetUserId(w, r)
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm := r.FormValue("password_confirm")

	if password != confirm {
		http.Error(w, "Password not match", http.StatusInternalServerError)
		return
	}

	user := models.Users{
		Id:    id,
		Name:  name,
		Email: email,
	}

	user.SetPassword(password)

	database.DB.Model(&user).Updates(&user)

	jrender.JSON(w, 200, user)

}
