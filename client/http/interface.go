package http

import (
	"go-svc-tpl/logs"

	"github.com/asmcos/requests"
	"github.com/imroc/req/v3"
)

type IHttpService interface {
	HttpGetRequest(url string, headers requests.Header) string
	HttpPostRequest(url string, headers map[string]string, body string) (string, error)
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

func (h *HttpServiceImpl) HttpPostRequest(url string, headers map[string]string, body string) (string, error) {
	resp, err := req.DevMode().R().SetHeaders(headers).SetBody(body).Post(url)
	if err != nil {
		logs.Error("[HttpPostRequest] %v", err)
		return "", err
	}
	return resp.ToString()
}
