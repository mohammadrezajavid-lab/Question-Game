package main

import (
	"encoding/json"
	"fmt"
	"gocasts.ir/go-fundamentals/gameapp/repository/mysql"
	"gocasts.ir/go-fundamentals/gameapp/service/user"
	"io"
	"log"
	"net/http"
	"strings"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	db := mysql.NewDB()
	if err := db.MysqlConnection.Ping(); err != nil {

		log.Println(err)

		res := NewResponse("unexpected error: ping to database server failed", "")
		ResponseWrite(w, r, res, http.StatusInternalServerError)

		return
	}

	res := NewResponse("", "health check OK")
	ResponseWrite(w, r, res, http.StatusOK)
}

type Response struct {
	Error   string `json:"error"`
	Message any    `json:"message"`
}

func NewResponse(error string, message any) *Response {
	return &Response{Error: error, Message: message}
}

func ResponseWrite(w http.ResponseWriter, r *http.Request, response *Response, statusCode int) {

	marshalResponse, mErr := json.Marshal(response)
	if mErr != nil {
		log.Println(mErr.Error())
	}

	w.WriteHeader(statusCode)

	if _, writeErr := fmt.Fprint(w, string(marshalResponse)); writeErr != nil {

		log.Println(writeErr.Error())
	}
}

type UserHandlers struct {
	UserService *user.Service
}

var (
	SignKey = []byte(`jwt_secret`)
)

func NewUserHandlers() *UserHandlers {
	return &UserHandlers{UserService: user.NewService(mysql.NewDB(), SignKey)}
}

func (uh *UserHandlers) userRegisterHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("recived one request from addres: %s, URL: %s", r.RemoteAddr, r.URL)

	if !strings.EqualFold(r.Method, http.MethodPost) {

		res := NewResponse("invalid method request", "")
		ResponseWrite(w, r, res, http.StatusMethodNotAllowed)

		return
	}

	requestBody, readErr := io.ReadAll(r.Body)
	if readErr != nil {

		log.Println(fmt.Errorf("can't read request body, error: %w", readErr))

		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest)

		return
	}

	if strings.EqualFold(string(requestBody), "") {

		// when empty body request
		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest)

		return
	}

	var requestUser = user.NewRegisterRequest("", "", "")
	if uErr := json.Unmarshal(requestBody, requestUser); uErr != nil {

		res := NewResponse(uErr.Error(), "")
		ResponseWrite(w, r, res, http.StatusBadRequest)

		return
	}

	//var db = mysql.NewDB()
	//var userService = user.NewService(db)

	registerResponse, registerErr := uh.UserService.Register(requestUser)
	if registerErr != nil {

		res := NewResponse(registerErr.Error(), "")
		ResponseWrite(w, r, res, http.StatusOK)

		return
	}

	res := NewResponse("", registerResponse)
	ResponseWrite(w, r, res, http.StatusOK)

}

func (uh *UserHandlers) userLoginHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("recived one request from addres: %s, URL: %s", r.RemoteAddr, r.URL)

	if !strings.EqualFold(r.Method, http.MethodPost) {

		res := NewResponse("invalid method request", "")
		ResponseWrite(w, r, res, http.StatusMethodNotAllowed)

		return
	}

	requestBody, readErr := io.ReadAll(r.Body)
	if readErr != nil {

		log.Println(fmt.Errorf("can't read request body, error: %w", readErr))

		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest)

		return
	}

	if strings.EqualFold(string(requestBody), "") {

		// when empty body request
		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest)

		return
	}

	var requestLogin = user.NewLoginRequest("", "")
	if uErr := json.Unmarshal(requestBody, requestLogin); uErr != nil {

		log.Println(fmt.Errorf("can't Unmarshal requestBody in Login procces, error: %w", uErr))

		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest)

		return
	}

	// TODO -
	loginResponse, loginErr := uh.UserService.Login(requestLogin)
	if loginErr != nil {

		res := NewResponse(loginErr.Error(), "")
		ResponseWrite(w, r, res, http.StatusUnauthorized)

		return
	}

	res := NewResponse("", loginResponse)
	ResponseWrite(w, r, res, http.StatusOK)
}

func (uh *UserHandlers) userProfileHandler(w http.ResponseWriter, r *http.Request) {

	// TODO - we are sanitize userId in this handler after send userId to service layer

	if r.Method != http.MethodGet {
		res := NewResponse("invalid method request", "")
		ResponseWrite(w, r, res, http.StatusMethodNotAllowed)

		return
	}

	requestBody, readErr := io.ReadAll(r.Body)
	if readErr != nil {

		log.Println(fmt.Errorf("can't read request body, error: %w", readErr))

		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest)

		return
	}

	if strings.EqualFold(string(requestBody), "") {

		// when empty body request
		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest)

		return
	}

	var requestProfile = user.NewProfileRequest(0)
	if uErr := json.Unmarshal(requestBody, requestProfile); uErr != nil {

		log.Printf("can't Unmarshal requestProfile data: %v", uErr)

		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res, http.StatusBadRequest)

		return
	}

	profile, pErr := uh.UserService.Profile(requestProfile)
	if pErr != nil {

		res := NewResponse(pErr.Error(), "")
		ResponseWrite(w, r, res, http.StatusOK)

		return
	}

	res := NewResponse("", profile)
	ResponseWrite(w, r, res, http.StatusOK)
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
