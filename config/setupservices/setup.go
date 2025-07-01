package setupservices

import (
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/adapter/publisher"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/authhandler"
	"golang.project/go-fundamentals/gameapp/pkg/jwt"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/accesscontrolmysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/usermysql"
	"golang.project/go-fundamentals/gameapp/repository/redis/redismatching"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/backofficeuserservice"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"golang.project/go-fundamentals/gameapp/service/userservice"
	"golang.project/go-fundamentals/gameapp/validator/matchingvalidator"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
)

type SetupServices struct {
	AuthSvc           authenticationservice.Service
	AuthorizationSvc  authorizationservice.Service
	UserSvc           userservice.Service
	UserValidator     uservalidator.Validator
	BackOfficeUserSvc backofficeuserservice.Service
	MatchingSvc       matchingservice.Service
	MatchingValidator matchingvalidator.Validator
	PresenceClient    presenceclient.Client
	AuthHandler       authhandler.AuthHandler
}

func New(config httpservercfg.Config) *SetupServices {

	mysqlRepo := mysql.NewDB(config.DataBaseCfg)

	authSvc := authenticationservice.NewService(jwt.NewJWT(config.JwtCfg))

	mysqlAccessControl := accesscontrolmysql.NewDataBase(mysqlRepo)
	authorizationSvc := authorizationservice.NewService(mysqlAccessControl)

	mysqlUser := usermysql.NewDataBase(mysqlRepo)
	userSvc := userservice.NewService(mysqlUser, &authSvc)
	userValidator := uservalidator.NewValidator(mysqlUser)

	backOfficeUserSvc := backofficeuserservice.NewService()

	redisAdapter := redis.New(config.RedisCfg)

	presenceClient, _ := presenceclient.NewClient(config.GrpcPresenceClientCfg)
	redisPublisher := publisher.NewPublish(config.PublisherCfg, redisAdapter)

	matchingSvc := matchingservice.NewService(
		config.MatchingCfg,
		redismatching.NewRedisDb(redisAdapter, config.MatchingRepoCfg),
		&presenceClient,
		redisPublisher,
	)
	matchingValidator := matchingvalidator.NewValidator()

	jwt := jwt.NewJWT(config.JwtCfg)
	authHandler := authhandler.New(authSvc, jwt)

	return &SetupServices{
		AuthSvc:           authSvc,
		AuthorizationSvc:  authorizationSvc,
		UserSvc:           userSvc,
		UserValidator:     userValidator,
		BackOfficeUserSvc: backOfficeUserSvc,
		MatchingSvc:       matchingSvc,
		MatchingValidator: matchingValidator,
		PresenceClient:    presenceClient,
		AuthHandler:       authHandler,
	}
}
