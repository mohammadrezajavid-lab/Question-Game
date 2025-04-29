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

func helloHandler(w http.ResponseWriter, r *http.Request) {

	if _, writeErr := fmt.Fprint(w, `{"message":"welcome to this game"}`); writeErr != nil {

		log.Fatal(writeErr.Error())
	}

	log.Printf("recived one request from addres: %s, URL: %s", r.RemoteAddr, r.URL)
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

	//marshalRegisterRes, mErr := json.Marshal(registerResponse)
	//if mErr != nil {
	//
	//	resJson := NewResponse(mErr.Error(), "")
	//	ResponseWrite(w, r, resJson)
	//
	//	return
	//
	//}

	//resJson := NewResponse("", marshalRegisterRes)

	//ResponseWrite(w, r, resJson)

	fmt.Println(registerResponse.User.String())
}

func main() {

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/users/register", registerUserHandler)

	fmt.Printf("server is ready in Address: %s\n", "127.0.0.1:8080")

	if listenErr := http.ListenAndServe("127.0.0.1:8080", nil); listenErr != nil {
		log.Fatal(listenErr)
	}

}
