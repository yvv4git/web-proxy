package webproxy

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

func (wp *WebProxy) authMiddleware(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	authHeader := req.Header.Get("Proxy-Authorization")
	if authHeader == "" {
		return req, goproxy.NewResponse(req,
			goproxy.ContentTypeText, http.StatusProxyAuthRequired,
			"Proxy authentication required")
	}

	username, password, ok := wp.parseBasicAuth(authHeader)
	if !ok {
		return req, goproxy.NewResponse(req,
			goproxy.ContentTypeText, http.StatusProxyAuthRequired,
			"Invalid authentication format")
	}

	if !wp.authManager.CheckCredentials(username, password) {
		return req, goproxy.NewResponse(req,
			goproxy.ContentTypeText, http.StatusProxyAuthRequired,
			"Invalid username or password")
	}

	return req, nil
}

func (p *WebProxy) parseBasicAuth(auth string) (username, password string, ok bool) {
	if !strings.HasPrefix(auth, "Basic ") {
		return "", "", false
	}

	payload := strings.TrimPrefix(auth, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", "", false
	}

	pair := strings.SplitN(string(decoded), ":", 2)
	if len(pair) != 2 {
		return "", "", false
	}

	return pair[0], pair[1], true
}
