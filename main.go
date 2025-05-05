package main

import (
	"encoding/json"
	"fmt"
	"gocasts.ir/go-fundamentals/gameapp/repository/mysql"
	"gocasts.ir/go-fundamentals/gameapp/service/authorize"
	"gocasts.ir/go-fundamentals/gameapp/service/user"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	db := mysql.NewDB()
	if err := db.MysqlConnection.Ping(); err != nil {

		log.Println(err)

		res := NewResponse("unexpected error: ping to database server failed", "")
		ResponseWrite(w, r, res, http.StatusInternalServerError, "application/json")

		return
	}

	res := NewResponse("", "health check OK")
	ResponseWrite(w, r, res, http.StatusOK, "application/json")
}

type Response struct {
	Error   string `json:"error"`
	Message any    `json:"message"`
}

func NewResponse(error string, message any) *Response {
	return &Response{Error: error, Message: message}
}

func ResponseWrite(w http.ResponseWriter, r *http.Request, response *Response, statusCode int, contentType string) {

	marshalResponse, mErr := json.Marshal(response)
	if mErr != nil {
		log.Println(mErr.Error())
	}

	// Set StatusCode
	w.WriteHeader(statusCode)

	// Set Content-Type
	w.Header().Add("Content-Type", contentType)

	if _, writeErr := fmt.Fprint(w, string(marshalResponse)); writeErr != nil {

		log.Println(writeErr.Error())
	}
}

type UserHandlers struct {
	UserService *user.Service
}

const (
	JWTSignKey            = `jwt_secret_key`
	AccessExpirationTime  = time.Hour * 24
	RefreshExpirationTime = time.Hour * 24 * 7
	AccessSubject         = `at`
	RefreshSubject        = `ar`
)

func NewUserHandlers() *UserHandlers {
	return &UserHandlers{
		UserService: user.NewService(
			mysql.NewDB(),
			authorize.NewService(
				[]byte(JWTSignKey),
				AccessExpirationTime,
				RefreshExpirationTime,
				AccessSubject,
				RefreshSubject,
			),
		),
	}
}

func (uh *UserHandlers) userRegisterHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("recived one request from addres: %s, URL: %s", r.RemoteAddr, r.URL)

	if !strings.EqualFold(r.Method, http.MethodPost) {

		res := NewResponse("invalid method request", "")
		ResponseWrite(w, r, res, http.StatusMethodNotAllowed, "application/json")

		return
	}

	requestBody, readErr := io.ReadAll(r.Body)
	if readErr != nil {

		log.Println(fmt.Errorf("can't read request body, error: %w", readErr))

		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest, "application/json")

		return
	}

	if strings.EqualFold(string(requestBody), "") {

		// when empty body request
		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest, "application/json")

		return
	}

	var requestUser = user.NewRegisterRequest("", "", "")
	if uErr := json.Unmarshal(requestBody, requestUser); uErr != nil {

		res := NewResponse(uErr.Error(), "")
		ResponseWrite(w, r, res, http.StatusBadRequest, "application/json")

		return
	}

	//var db = mysql.NewDB()
	//var userService = user.NewService(db)

	registerResponse, registerErr := uh.UserService.Register(requestUser)
	if registerErr != nil {

		res := NewResponse(registerErr.Error(), "")
		ResponseWrite(w, r, res, http.StatusOK, "application/json")

		return
	}

	res := NewResponse("", registerResponse)
	ResponseWrite(w, r, res, http.StatusOK, "application/json")

}

func (uh *UserHandlers) userLoginHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("recived one request from addres: %s, URL: %s", r.RemoteAddr, r.URL)

	if !strings.EqualFold(r.Method, http.MethodPost) {

		res := NewResponse("invalid method request", "")
		ResponseWrite(w, r, res, http.StatusMethodNotAllowed, "application/json")

		return
	}

	requestBody, readErr := io.ReadAll(r.Body)
	if readErr != nil {

		log.Println(fmt.Errorf("can't read request body, error: %w", readErr))

		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest, "application/json")

		return
	}

	if strings.EqualFold(string(requestBody), "") {

		// when empty body request
		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest, "application/json")

		return
	}

	var requestLogin = user.NewLoginRequest("", "")
	if uErr := json.Unmarshal(requestBody, requestLogin); uErr != nil {

		log.Println(fmt.Errorf("can't Unmarshal requestBody in Login procces, error: %w", uErr))

		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest, "application/json")

		return
	}

	// TODO -
	loginResponse, loginErr := uh.UserService.Login(requestLogin)
	if loginErr != nil {

		res := NewResponse(loginErr.Error(), "")
		ResponseWrite(w, r, res, http.StatusUnauthorized, "application/json")

		return
	}

	res := NewResponse("", loginResponse)
	ResponseWrite(w, r, res, http.StatusOK, "application/json")
}

func (uh *UserHandlers) userProfileHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("recived one request from addres: %s, URL: %s", r.RemoteAddr, r.URL)

	// TODO - we are sanitize userId in this handler after send userId to service layer

	if r.Method != http.MethodGet {
		res := NewResponse("invalid method request", "")
		ResponseWrite(w, r, res, http.StatusMethodNotAllowed, "application/json")

		return
	}

	requestBody, readErr := io.ReadAll(r.Body)
	if readErr != nil {

		log.Println(fmt.Errorf("can't read request body, error: %w", readErr))

		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest, "application/json")

		return
	}

	if !strings.EqualFold(string(requestBody), "") {

		// when empty body request
		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest, "application/json")

		return
	}

	authService := authorize.NewService(
		[]byte(JWTSignKey),
		AccessExpirationTime,
		RefreshExpirationTime,
		AccessSubject,
		RefreshSubject,
	)
	tokenAuth := r.Header.Get("Authorization")
	claims, pErr := authService.ParseJWT(tokenAuth)
	if pErr != nil {
		log.Println(pErr)
	}

	if claims == nil {
		res := NewResponse("Unauthorized User", "")
		ResponseWrite(w, r, res, http.StatusUnauthorized, "application/json")

	} else {

		profile, pErr := uh.UserService.Profile(user.NewProfileRequest(claims.UserId))
		if pErr != nil {

			res := NewResponse(pErr.Error(), "")
			ResponseWrite(w, r, res, http.StatusOK, "application/json")

			return
		}

		res := NewResponse("", profile)
		ResponseWrite(w, r, res, http.StatusOK, "application/json")
	}

}

func main() {

	userHandlers := NewUserHandlers()

	http.HandleFunc("/health-check", healthCheckHandler)
	http.HandleFunc("/users/register", userHandlers.userRegisterHandler)
	http.HandleFunc("/users/login", userHandlers.userLoginHandler)
	http.HandleFunc("/users/profile", userHandlers.userProfileHandler)

	fmt.Printf("server is ready in Address: %s\n", "127.0.0.1:8080")

	if listenErr := http.ListenAndServe("127.0.0.1:8080", nil); listenErr != nil {
		log.Fatal(listenErr)
	}
}
