package consts

// Log tag
const (
	SLTagUndefined           = " _undef"
	SLTagRequestIn           = " _com_request_in"
	SLTagRequestOut          = " _com_request_out"
	SLTagHTTPFailed          = " _http_fail"
	SLTagSpaceHTTPFailed     = " _space_http_fail"
	SLTagHTTPSuccess         = " _http_success"
	SLTagMysqlFail           = " _mysql_fail"
	SLTagMysqlRecordNotFound = " _mysql_record_not_found"
	SLTagRedisSuccess        = " _redis_success"
	SLTagRedisFail           = " _redis_fail"
	SLTagPanic               = " _panic"
	SLTagSeverStart          = " _server_start"
	SLTagSeverFail           = " _server_fail"
	SLTagPprofFail           = " _pprof_fail"
	SLtagJsonMarshal         = " _json_marshal"
	SLTagNotificationFail    = " _notification_fail"
	SLTagBrokePipe           = " _broke_pipe"
	SLTagConnRST             = " _connection_reset_by_peer"
)

// for middleware
const (
	AUTH_HEADER        = "Authorization"
	TimeMilliRequestIn = "TimeMilliRequestIn"

	SHARED_LIB_TOKEN_KEY = "Token"
)

// Header
const (
	ContentType = "Content-Type"
	Location    = "Location"
)

var NotAuthApisMap = map[string]string{
	UserPrefix + SendVerificationCodeApi: UserPrefix + SendVerificationCodeApi,
	UserPrefix + RegisterApi:             UserPrefix + RegisterApi,
	UserPrefix + LoginApi:                UserPrefix + LoginApi,
	"/swagger/*any":                      "/swagger/*any",
}

var CanGenerateRechargeKeyUserMap = map[string]string{
	"stalary@163.com": "stalary@163.com",
}

const (
	UserPrefix              = "/api/v1/user"
	SendVerificationCodeApi = "/verify/send_code"
	RegisterApi             = "/register"
	LoginApi                = "/login"
	ProfileApi              = "/profile"

	IntegralPrefix         = "/api/v1/integral"
	GenerateRechargeKeyApi = "/generate_key"
	RechargeApi            = "/recharge"

	OpenaiPrefix         = "/api/v1/openai"
	OpenaiCompletionsApi = "/v1/chat/completions"
)

// env
const (
	Env = "ENV"
)
