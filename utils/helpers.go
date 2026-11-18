package utils

import "net/url"

func IsValidURL(endpoint string) bool {
	url, err := url.Parse(endpoint)
	if err != nil || url == nil || url.Host == "" || url.Scheme == "" {
		return false
	}
	return true
}
