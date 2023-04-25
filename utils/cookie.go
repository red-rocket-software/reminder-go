package utils

import "net/http"

func CreateCookie(name, value, path, domain string, maxAge int, secure, httpOnly bool) http.Cookie {
	cookie := http.Cookie{}
	cookie.Name = name
	cookie.Value = value
	cookie.Path = path
	cookie.Domain = domain
	cookie.MaxAge = maxAge
	cookie.Secure = secure
	cookie.HttpOnly = httpOnly

	return cookie
}

func CreateExpiredCookie(name, value, path, domain string, secure, httpOnly bool) http.Cookie {
	cookie := http.Cookie{}
	cookie.Name = name
	cookie.Value = value
	cookie.Path = path
	cookie.Domain = domain
	cookie.MaxAge = -1
	cookie.Secure = secure
	cookie.HttpOnly = httpOnly

	return cookie
}
