package auth

import (
	"apiwithmysql/model"
	"apiwithmysql/repo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secreat_Key = []byte("This_is_secreat_Key")
)

type User struct {
	UserName string `json:"username,omitempty" gorm:"primarykey"`
	Password string `json:"password,omitempty"`
}

type Claims struct {
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	fatal(err)

	isPresent := repo.IsUserPresent(&user)
	if isPresent {
		w.WriteHeader(410)
		json.NewEncoder(w).Encode("user already present")
		return
	}

	claims := &Claims{
		UserName: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{time.Now().Add(5 * time.Minute)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secreat_Key)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("from tokenString : ", err)
		return
	}

	saved := repo.SaveUser(&user)
	if !saved {
		w.WriteHeader(411)
		json.NewEncoder(w).Encode("failed in saving user")
		return
	}

	mp := map[string]string{"token": tokenString}
	json.NewEncoder(w).Encode(&mp)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var ru model.User
	err := json.NewDecoder(r.Body).Decode(&ru)

	if err != nil {
		log.Fatal(err)
	}

	userFrmDb, isPresent := repo.GetUser(&ru)
	if !isPresent {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("user not present")
		return
	}

	if ru.Password != userFrmDb.Password {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("user password incorrect")
		return
	}

	claims := &Claims{
		UserName: ru.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{time.Now().Add(5 * time.Minute)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secreat_Key)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("from tokenString : ", err)
		return
	}

	mp := map[string]string{"token": tokenStr}
	json.NewEncoder(w).Encode(&mp)
}

type Token struct {
	TokenStr string `json:"tokenstring"`
}

func Authorise(w http.ResponseWriter, r *http.Request) bool {
	var rt Token
	err := json.NewDecoder(r.Body).Decode(&rt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	tokenStr := rt.TokenStr
	if tokenStr == "" {
		json.NewEncoder(w).Encode("token is empty")
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) { return secreat_Key, nil })

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("You are un authorised")
		return false
	}

	fmt.Println("Hi welcome : " + claims.UserName)
	return true
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
