package stockapi

type EchoReq struct {
	Msg string `json:"msg" binding:"required"`
}

type EchoResponse struct {
	Text string `json:"text" binding:"required"`
}
