package webproxy

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

func (wp *WebProxy) authMiddleware(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	authHeader := req.Header.Get("Proxy-Authorization")
	if authHeader == "" {
		return req, proxyAuthRequired(req, "Proxy authentication required")
	}

	// Temporarily copying it to Authorization to use the standard BasicAuth() parser
	req.Header.Set("Authorization", authHeader)
	username, password, ok := req.BasicAuth()
	req.Header.Del("Authorization")

	if !ok {
		return req, proxyAuthRequired(req, "Invalid authentication format")
	}

	if !wp.authManager.CheckCredentials(username, password) {
		return req, proxyAuthRequired(req, "Invalid username or password")
	}

	return req, nil
}

func proxyAuthRequired(req *http.Request, msg string) *http.Response {
	return goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusProxyAuthRequired, msg)
}
