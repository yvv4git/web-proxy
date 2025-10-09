package webproxy

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

func (wp *WebProxy) authMiddleware(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// todo: implement authentication logic here
	return req, nil
}
