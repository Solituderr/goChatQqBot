package http

import (
	"go-svc-tpl/logs"

	"github.com/asmcos/requests"
)

type IHttpService interface {
	HttpGetRequest(url string, headers requests.Header) string
}

type HttpServiceImpl struct {
}

func NewHttpService() IHttpService {
	return &HttpServiceImpl{}
}

func (h *HttpServiceImpl) HttpGetRequest(url string, headers requests.Header) string {
	resp, err := requests.Get(url, headers)
	if err != nil {
		logs.Error("[HttpGetRequest] %v", err)
		return ""
	}
	return resp.Text()
}
