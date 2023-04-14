package consts

// Log tag
const (
	SLTagUndefined           = " _undef"
	SLTagRequest             = " _request"
	SLTagHTTPFailed          = " _http_fail"
	SLTagSpaceHTTPFailed     = " _space_http_fail"
	SLTagHTTPSuccess         = " _http_success"
	SLTagMysqlFail           = " _mysql_fail"
	SLTagMysqlRecordNotFound = " _mysql_record_not_found"
	SLTagRedisSuccess        = " _redis_success"
	SLTagRedisFail           = " _redis_fail"
	SLTagPanic               = " _panic"
	SLTagSeverStart          = " _server_start"
	SLTagSeverStop           = " _server_stop"
	SLTagSeverFail           = " _server_fail"
	SLTagPprofFail           = " _pprof_fail"
	SLtagJsonMarshal         = " _json_marshal"
	SLTagNotificationFail    = " _notification_fail"
	SLTagBrokePipe           = " _broke_pipe"
	SLTagConnRST             = " _connection_reset_by_peer"
	SyncStop                 = " _async_ready_stop"
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
	IntegralPrefix + ManualRechargeApi:   IntegralPrefix + ManualRechargeApi,
	"/swagger/*any":                      "/swagger/*any",
}

var CanGenerateRechargeKeyUserMap = map[string]string{
	"stalary@163.com":    "stalary@163.com",
	"vinson_neo@163.com": "vinson_neo@163.com",
}

const (
	UserPrefix              = "/api/v1/user"
	SendVerificationCodeApi = "/verify/send_code"
	RegisterApi             = "/register"
	LoginApi                = "/login"
	ProfileApi              = "/profile"
	FeedbackApi             = "/feedback"

	IntegralPrefix         = "/api/v1/integral"
	GenerateRechargeKeyApi = "/generate_key"
	RechargeApi            = "/recharge"
	ManualRechargeApi      = "/manual/recharge"
	RecordApi              = "/record"

	OpenaiPrefix         = "/api/v1/openai"
	OpenaiCompletionsApi = "/v1/chat/completions"
	ImageApi             = "/v1/image"
	AudioApi             = "/v1/audio"
)

// env
const (
	Env = "ENV"
)
