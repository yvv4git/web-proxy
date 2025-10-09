package webproxy

import "sync"

// Noop account manager

type NoopAuthManager struct{}

func NewNoopAuthManager() *NoopAuthManager {
	return &NoopAuthManager{}
}

func (am *NoopAuthManager) CheckCredentials(username, password string) bool {
	return true
}

// Basic account manager

type BasicAuthManager struct {
	mx       sync.RWMutex
	accounts map[string]string
}

func NewBasicAuthManager() *BasicAuthManager {
	return &BasicAuthManager{
		accounts: make(map[string]string),
	}
}

func (am *BasicAuthManager) AddAccount(username, password string) {
	am.mx.Lock()
	defer am.mx.Unlock()

	am.accounts[username] = password
}

func (am *BasicAuthManager) CheckCredentials(username, password string) bool {
	am.mx.RLock()
	defer am.mx.RUnlock()

	if am.accounts[username] == password {
		return true
	}

	return false
}
