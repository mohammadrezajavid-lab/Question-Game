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
	}

	res := NewResponse("", "health check OK")
	ResponseWrite(w, r, res)
}

type Response struct {
	Error   string `json:"error"`
	Message any    `json:"message"`
}

func NewResponse(error string, message any) *Response {
	return &Response{Error: error, Message: message}
}

func ResponseWrite(w http.ResponseWriter, r *http.Request, response *Response) {

	marshalResponse, mErr := json.Marshal(response)
	if mErr != nil {

	}

	if _, writeErr := fmt.Fprint(w, string(marshalResponse)); writeErr != nil {

		log.Println(writeErr.Error())
	}
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("recived one request from addres: %s, URL: %s", r.RemoteAddr, r.URL)

	if !strings.EqualFold(r.Method, http.MethodPost) {

		res := NewResponse("invalid method request", "")
		ResponseWrite(w, r, res)

		return
	}

	requestBody, readErr := io.ReadAll(r.Body)
	if readErr != nil {

		res := NewResponse(readErr.Error(), "")
		ResponseWrite(w, r, res)

		return
	}

	if strings.EqualFold(string(requestBody), "") {

		// when empty body request
		err := "invalid body request"
		res := NewResponse(err, "")
		ResponseWrite(w, r, res)

		return
	}

	var requestUser = user.NewRegisterRequest("", "")
	if uErr := json.Unmarshal(requestBody, requestUser); uErr != nil {

		res := NewResponse(uErr.Error(), "")
		ResponseWrite(w, r, res)

		return
	}

	var db = mysql.NewDB()
	var userService = user.NewService(db)

	registerResponse, registerErr := userService.Register(requestUser)
	if registerErr != nil {

		res := NewResponse(registerErr.Error(), "")
		ResponseWrite(w, r, res)

		return
	}

	res := NewResponse("", registerResponse)
	ResponseWrite(w, r, res)

}

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	http.Redirect(w, r, "/users/register", http.StatusTemporaryRedirect)
}

func main() {

	http.HandleFunc("/health-check", healthCheckHandler)
	http.HandleFunc("/users/register", registerUserHandler)
	http.HandleFunc("/users/login", loginUserHandler)

	fmt.Printf("server is ready in Address: %s\n", "127.0.0.1:8080")

	if listenErr := http.ListenAndServe("127.0.0.1:8080", nil); listenErr != nil {
		log.Fatal(listenErr)
	}

}
