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

const (
	RechargeNotifySubject = "ChatAlpha 积分充值成功通知"
	RechargeNotifyContent = "尊敬的 ChatAlpha 用户\n\n您好！感谢您使用 ChatAlpha 服务。\n\n我们很高兴地告诉您，您的账户已经成功充值了 %d 积分，这些积分可以用于对话等操作。充值成功后积分会即时到账，您可以随时查看积分余额。\n\n如有任何问题或疑问，请随时联系我们的客服团队，我们将尽快为您解决问题。\n\nChatAlpha 团队"
)
