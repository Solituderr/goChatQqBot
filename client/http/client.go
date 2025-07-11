package http

var HttpService IHttpService

func InitHttpClient() {
	HttpService = NewHttpService()
}
