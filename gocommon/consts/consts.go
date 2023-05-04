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
	SLTagJsonMarshal         = " _json_marshal"
	SLTagNotificationFail    = " _notification_fail"
	SLTagBrokePipe           = " _broke_pipe"
	SLTagConnRST             = " _connection_reset_by_peer"
	SyncStop                 = " _async_ready_stop"
	SLTagAlipaySuccess       = " _alipay_success"
	MailSyncFailRetry        = " _mail_sync_fail_retry"
	SLTagAliOssSuccess       = " _ali_oss_success"
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
	UserPrefix + SendVerificationCodeApi:         UserPrefix + SendVerificationCodeApi,
	UserPrefix + RegisterApi:                     UserPrefix + RegisterApi,
	UserPrefix + LoginApi:                        UserPrefix + LoginApi,
	UserPrefix + SendPasswordVerificationCodeApi: UserPrefix + SendPasswordVerificationCodeApi,
	UserPrefix + ChangePasswordApi:               UserPrefix + ChangePasswordApi,
	IntegralPrefix + ManualRechargeApi:           IntegralPrefix + ManualRechargeApi,
	AlipayPrefix + NotifyApi:                     AlipayPrefix + NotifyApi,
	"/public/*filepath":                          "/public/*filepath",
	"/favicon.ico":                               "/favicon.ico",
	"/swagger/*any":                              "/swagger/*any",
	"/swagger/index.html":                        "/swagger/index.html",
}

var CanGenerateRechargeKeyUserMap = map[string]string{
	"stalary@163.com":    "stalary@163.com",
	"vinson_neo@163.com": "vinson_neo@163.com",
}

const (
	UserPrefix                      = "/api/v1/user"
	SendVerificationCodeApi         = "/verify/send_code"
	RegisterApi                     = "/register"
	LoginApi                        = "/login"
	LogoutApi                       = "/logout"
	ProfileApi                      = "/profile"
	FeedbackApi                     = "/feedback"
	ChangePasswordApi               = "/change_password"
	SendPasswordVerificationCodeApi = "/change_password/verify/code"
	MenuApi                         = "/menu"

	IntegralPrefix         = "/api/v1/integral"
	GenerateRechargeKeyApi = "/generate_key"
	RechargeApi            = "/recharge"
	ManualRechargeApi      = "/manual/recharge"

	OpenaiPrefix         = "/api/v1/openai"
	OpenaiCompletionsApi = "/v1/chat/completions"
	ImageApi             = "/v1/image"
	AudioApi             = "/v1/audio"

	AlipayPrefix = "/api/v1/alipay"
	NotifyApi    = "/notify"

	PayPrefix    = "/api/v1/pay"
	PreCreateApi = "/pre_create"
	StatusApi    = "/status"
)

// env
const (
	Env = "ENV"
)

const (
	RechargeNotifySubject = "ChatAlpha 积分充值成功通知"
	RechargeNotifyContent = "尊敬的 ChatAlpha 用户\n\n您好！感谢您使用 ChatAlpha 服务。\n\n我们很高兴地告诉您，您的账户已经成功充值了 %d 积分，这些积分可以用于对话等操作。充值成功后积分会即时到账，您可以随时查看积分余额。\n\n如有任何疑问或需要帮助，请随时邮件联系我们。\n\n谢谢！\n\nChatAlpha 团队"

	SendCodeSubject = "ChatAlpha 账户安全代码"
	SendCodeContent = "尊敬的 ChatAlpha 用户\n\n您好！感谢您使用 ChatAlpha 服务。\n\n为了确保您的账户安全，我们已向您发送了一封验证码邮件，请勿将验证码泄露给他人。验证码用于验证您的身份，并防止恶意攻击。\n\n验证码：【%s】\n\n如果您没有进行任何操作，或者不希望继续使用 ChatAlpha 服务，请忽略此邮件。\n\n如有任何疑问或需要帮助，请随时邮件联系我们。\n\n谢谢！\n\nChatAlpha 团队"

	InviteSubject = "ChatAlpha 邀请注册成功通知"
	InviteContent = "尊敬的 ChatAlpha 用户\n\n您好！感谢您使用 ChatAlpha 服务，并成功邀请了新用户加入我们的平台。\n\n我们很高兴地告诉您，您邀请的用户【%s】已经成功进行了注册，我们已经为您的账户充值了 %d 积分作为奖励，这些积分可以用于对话等操作。充值成功后积分会即时到账，您可以随时查看积分余额。\n\n如有任何疑问或需要帮助，请随时邮件联系我们。\n\n谢谢！\n\nChatAlpha 团队"
)
