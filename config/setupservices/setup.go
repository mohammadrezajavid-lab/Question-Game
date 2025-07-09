package setupservices

import (
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/adapter/publisher"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/adapter/subscriber"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/authhandler"
	"golang.project/go-fundamentals/gameapp/pkg/jwt"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/accesscontrolmysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/gamemysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/questionmysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/usermysql"
	"golang.project/go-fundamentals/gameapp/repository/redis/redismatching"
	"golang.project/go-fundamentals/gameapp/repository/redis/redisquiz"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/backofficeuserservice"
	"golang.project/go-fundamentals/gameapp/service/gameservice"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"golang.project/go-fundamentals/gameapp/service/quizservice"
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
	GameSvc           gameservice.Service
	QuizSvc           quizservice.Service
	MatchingValidator matchingvalidator.Validator
	PresenceClient    presenceclient.Client
	AuthHandler       authhandler.AuthHandler
}

func New(config httpservercfg.Config) *SetupServices {

	mysqlDB := mysql.NewDB(config.DataBaseCfg)

	authSvc := authenticationservice.NewService(jwt.NewJWT(config.JwtCfg))

	accessControlRepo := accesscontrolmysql.NewDataBase(mysqlDB)
	authorizationSvc := authorizationservice.NewService(accessControlRepo)

	userRepo := usermysql.NewDataBase(mysqlDB)
	userSvc := userservice.NewService(userRepo, &authSvc)
	userValidator := uservalidator.NewValidator(userRepo)

	backOfficeUserSvc := backofficeuserservice.NewService()

	redisAdapter := redis.New(config.RedisCfg)

	presenceClient, _ := presenceclient.NewClient(config.GrpcPresenceClientCfg)

	redisPublisher := publisher.NewPublisher(config.PublisherCfg, redisAdapter)
	redisSubscriber := subscriber.NewSubscriber(redisAdapter, config.SubscriberCfg)

	matchingSvc := matchingservice.NewService(
		config.MatchingCfg,
		redismatching.NewRedisDb(redisAdapter, config.MatchingRepoCfg),
		&presenceClient,
		redisPublisher,
	)
	matchingValidator := matchingvalidator.NewValidator()

	gameRepo := gamemysql.NewDataBase(mysqlDB)
	gameSvc := gameservice.New(redisAdapter, gameRepo, redisPublisher, redisSubscriber, config.GameServiceCfg)

	quizSvc := quizservice.New(config.QuizServiceCfg, redisquiz.NewRedisDb(redisAdapter), questionmysql.NewDataBase(mysqlDB))

	jwt := jwt.NewJWT(config.JwtCfg)
	authHandler := authhandler.New(authSvc, jwt)

	return &SetupServices{
		AuthSvc:           authSvc,
		AuthorizationSvc:  authorizationSvc,
		UserSvc:           userSvc,
		UserValidator:     userValidator,
		BackOfficeUserSvc: backOfficeUserSvc,
		MatchingSvc:       matchingSvc,
		GameSvc:           gameSvc,
		QuizSvc:           quizSvc,
		MatchingValidator: matchingValidator,
		PresenceClient:    presenceClient,
		AuthHandler:       authHandler,
	}
}
