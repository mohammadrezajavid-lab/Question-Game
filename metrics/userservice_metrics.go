package metrics

import "github.com/prometheus/client_golang/prometheus"

var FailedLoginCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_login_total",
		Help: "Total number of failed login",
	},
)

var FailedLoginIncorrectPhoneNumberOrPasswordCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_login_incorrect_phonenumber_password_total",
		Help: "Total number of failed login because incorrect phoneNumber of password",
	},
)

var FailedCreateAccessTokenCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_create_access_token_total",
		Help: "Total number of failed create access token",
	},
)

var FailedCreateRefreshTokenCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_create_refresh_token_total",
		Help: "Total number of failed create refresh token",
	},
)

var FailedGetUserByIDCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_getuserbyid_total",
		Help: "Total number of failed GetUserById",
	},
)

var FailedRegisterUserCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_registeruser_total",
		Help: "Total number of failed RegisterUser",
	},
)

func init() {
	Registry.MustRegister(
		FailedLoginCounter,
		FailedLoginIncorrectPhoneNumberOrPasswordCounter,
		FailedCreateAccessTokenCounter,
		FailedCreateRefreshTokenCounter,
		FailedGetUserByIDCounter,
		FailedRegisterUserCounter,
	)
}
