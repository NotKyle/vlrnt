package utils

import (
	"encoding/json"
	"errors"
	"strings"
)

func GetMethod(url string) (string, error) {
	if url == "" {
		return "", errors.New("url cannot be empty")
	}

	method := strings.Split(url, ":/")[0]

	return method, nil
}

func GetDomain(url string) (string, error) {
	if url == "" {
		return "", errors.New("url cannot be empty")
	}

	urlParts := strings.Split(url, "/")

	if len(urlParts) < 3 {
		return "", errors.New("url is not valid")
	}

	domain := urlParts[2]

	return domain, nil
}

func GetTLD(url string) (string, error) {
	if url == "" {
		return "", errors.New("url cannot be empty")
	}

	domain, err := GetDomain(url)

	if err != nil {
		return "", err
	}

	domainParts := strings.Split(domain, ".")

	if len(domainParts) < 2 {
		return "", errors.New("domain is not valid")
	}

	tld := domainParts[len(domainParts)-1]

	return tld, nil
}

func GetPath(url string) (string, error) {
	if url == "" {
		return "", errors.New("url cannot be empty")
	}

	urlParts := strings.Split(url, "/")

	if len(urlParts) < 3 {
		return "", errors.New("url is not valid")
	}

	path := strings.Join(urlParts[3:], "/")

	return path, nil
}

func GetSubDomain(url string) (string, error) {
	if url == "" {
		return "", errors.New("url cannot be empty")
	}

	domain, err := GetDomain(url)

	if err != nil {
		return "", err
	}

	domainParts := strings.Split(domain, ".")

	if len(domainParts) < 3 {
		return "", errors.New("domain is not valid")
	}

	subDomain := domainParts[0]

	return subDomain, nil
}

// Convert Struct to JSON
func AsJson(obj interface{}) (string, error) {
	json, err := json.Marshal(obj)

	if err != nil {
		return "", err
	}

	return string(json), nil
}
