package main

import (
	"bytes"
	"encoding/json"
	"gocasts.ir/go-fundamentals/gameapp/service/user"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type data struct {
	testName        string
	httpMethod      string
	registerRequest *user.RegisterRequest
	response        *Response
	statusCode      int
}

func TestRegisterUserHttpServer(t *testing.T) {
	testCases := []data{
		{
			"test httpMethod",
			http.MethodGet,
			user.NewRegisterRequest("ali", "09196881929"),
			NewResponse("invalid method request", ""),
			http.StatusOK,
		},
		{
			"invalid name",
			http.MethodPost,
			user.NewRegisterRequest("a", "09196881929"),
			NewResponse("name length should be greater than 3", ""),
			http.StatusOK,
		},
		{
			"test empty body request",
			http.MethodPost,
			nil,
			NewResponse("invalid body request", ""),
			http.StatusOK,
		},
		{
			"register new user",
			http.MethodPost,
			user.NewRegisterRequest("ali", "09126551927"),
			NewResponse("", `{"user": {"name": "ali","phone_number": "09126551927"}}`),
			http.StatusOK,
		},
		{
			"phone number is invalid",
			http.MethodPost,
			user.NewRegisterRequest("ali", "0912"),
			NewResponse("phone number is invalid", ""),
			http.StatusOK,
		},
	}

	for index, input := range testCases {
		t.Log(input.testName)

		reader := bytes.NewReader([]byte(""))

		if testCases[index].registerRequest != nil {
			marshalUser, _ := json.Marshal(testCases[index].registerRequest)
			reader = bytes.NewReader(marshalUser)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(input.httpMethod, "/users/register", reader)
		registerUserHandler(w, r)

		resRegisterUser := w.Result()
		resRegisterUserBody, _ := io.ReadAll(resRegisterUser.Body)

		res := NewResponse("", "")
		_ = json.Unmarshal(resRegisterUserBody, res)

		if !strings.EqualFold(res.Error, testCases[index].response.Error) {

			log.Fatalf("mismatch response error, expected: %s\n got: %s", input.response.Error, res.Error)
		}
	}

}
