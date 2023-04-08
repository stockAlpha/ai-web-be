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
	Prefix + SendVerificationCodeApi: Prefix + SendVerificationCodeApi,
	Prefix + RegisterApi:             Prefix + RegisterApi,
	Prefix + LoginApi:                Prefix + LoginApi,
}

const (
	Prefix                  = "/api/stock/web"
	SendVerificationCodeApi = "/v1/verify/send_code"
	RegisterApi             = "/v1/user/register"
	LoginApi                = "/v1/user/login"
)
