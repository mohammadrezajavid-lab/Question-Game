package errormessage

const (
	ErrorMsgUnexpected             = "unexpected error: Try again later."
	ErrorMsgScanQuery              = "can't scan query result"
	ErrorMsgUserNotAllowed         = "user not allowed"
	ErrorMsgInvalidCategoryType    = "invalid category type"
	ErrorMsgInvalidCategory        = "invalid category"
	ErrorMsgInvalidRequest         = "invalid value request body"
	ErrorMsgInvalidPhoneType       = "invalid phone number type"
	ErrorMsgPhoneNotUniq           = "phone number is not uniq"
	ErrorMsgNotExistPhoneNumber    = "phone number does not exist"
	ErrorMsgFailedOpenMysqlConn    = "failed to open MySQL connection"
	ErrorMsgConnectionRefusedMysql = "connection refused db"
	ErrorMsgFailedExecuteQuery     = "can't execute query"
	ErrorMsgFailedStartServer      = "failed to start server error"
	ErrorMsgFailedCreateSch        = "failed to create scheduler"
	ErrorMsgHttpServerShutdown     = "http server shutdown error"
	ErrorMsgMetricsServerShutdown  = "metrics server shutdown error"
	ErrorMsgShutdownSch            = "not successfully shutdown scheduler"
)
