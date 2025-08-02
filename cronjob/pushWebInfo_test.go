package cronjob

import (
	"go-svc-tpl/client"
	"testing"
)

func Test_extractUrl(t *testing.T) {
	client.Init()
	extractUrl()
}
