package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anyongjitiger/photo-backup-server/utils"
	"github.com/anyongjitiger/photo-backup-server/web/common"
	"github.com/anyongjitiger/photo-backup-server/web/core/render"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	mux "github.com/julienschmidt/httprouter"
)

const (
	SecretKey = "welcome to wangshubo's blog"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type UserCredentials struct {
	Username string `json:"username"`
    Password string `json:"password"`
    PIN string `json:"pin"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)
}

func ProtectedHandler2(w http.ResponseWriter, r *http.Request) {
	// get httprouter.Params from request context
	if ps, ok := common.ParamsFromContext(r.Context()); ok {
		response := Response{"Hello," + ps.ByName("name")}
		JsonResponse(response, w)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request, pm mux.Params) {

	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	/* if strings.ToLower(user.Username) != "someone" || user.Password != "p@ssword" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprint(w, "Invalid credentials")
		return
    } */
    // log.Println(user.PIN)
    // log.Println(utils.GetCurrentPIN())
    if user.PIN != utils.GetCurrentPIN() {
        w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprint(w, "Invalid credentials")
		return
    }

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		fatal(err)
	}
    pin := utils.GetCurrentPIN()
    // tokenString, err := token.SignedString([]byte(SecretKey))
    tokenString, err := token.SignedString([]byte(pin))
    
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	response := Token{tokenString}
	// JsonResponse(response, w)
	render.RenderJson(w, response)
}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			// return []byte(SecretKey), nil
			return []byte(utils.GetCurrentPIN()), nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}

}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
