package setupservices

import (
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/accesscontrolmysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/usermysql"
	"golang.project/go-fundamentals/gameapp/repository/redis/redismatching"
	"golang.project/go-fundamentals/gameapp/repository/redis/redispresence"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/backofficeuserservice"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"golang.project/go-fundamentals/gameapp/service/presenceservice"
	"golang.project/go-fundamentals/gameapp/service/userservice"
	"golang.project/go-fundamentals/gameapp/validator/matchingvalidator"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
)

type SetupServices struct {
	AuthSvc           *authenticationservice.Service
	AuthorizationSvc  *authorizationservice.Service
	UserSvc           *userservice.Service
	UserValidator     *uservalidator.Validator
	BackOfficeUserSvc *backofficeuserservice.Service
	MatchingSvc       *matchingservice.Service
	MatchingValidator *matchingvalidator.Validator
	PresenceSvc       *presenceservice.Service
}

func New(config httpservercfg.Config) *SetupServices {

	mysqlRepo := mysql.NewDB(config.DataBaseCfg)

	authSvc := authenticationservice.NewService(
		authenticationservice.NewConfig(
			config.AuthCfg.SignKey,
			config.AuthCfg.AccessExpirationTime,
			config.AuthCfg.RefreshExpirationTime,
			config.AuthCfg.AccessSubject,
			config.AuthCfg.RefreshSubject),
	)

	mysqlAccessControl := accesscontrolmysql.NewDataBase(mysqlRepo)
	authorizationSvc := authorizationservice.NewService(mysqlAccessControl)

	mysqlUser := usermysql.NewDataBase(mysqlRepo)
	userSvc := userservice.NewService(mysqlUser, authSvc)
	userValidator := uservalidator.NewValidator(mysqlUser)

	backOfficeUserSvc := backofficeuserservice.NewService()

	redisAdapter := redis.New(config.RedisCfg)

	matchingSvc := matchingservice.NewService(config.MatchingCfg, redismatching.NewRedisDb(redisAdapter))
	matchingValidator := matchingvalidator.NewValidator()

	presenceSvc := presenceservice.New(redispresence.NewRedisDb(redisAdapter), config.PresenceCfg)

	return &SetupServices{
		AuthSvc:           authSvc,
		AuthorizationSvc:  authorizationSvc,
		UserSvc:           userSvc,
		UserValidator:     userValidator,
		BackOfficeUserSvc: backOfficeUserSvc,
		MatchingSvc:       matchingSvc,
		MatchingValidator: matchingValidator,
		PresenceSvc:       presenceSvc,
	}
}
