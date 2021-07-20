package helper

type WebResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewWebResponse() WebResponse {
	return WebResponse{}
}
