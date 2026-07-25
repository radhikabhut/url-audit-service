package validator

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}

func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}

func ValidateSafetyURL(urlStr string) error {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return errors.New("malformed URL")
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return errors.New("only http and https schemes are allowed")
	}

	hostname := u.Hostname()
	if hostname == "" {
		return errors.New("empty hostname")
	}

	// Reject localhost string representation directly
	if strings.ToLower(hostname) == "localhost" {
		return errors.New("connections to localhost are rejected")
	}

	// Resolve the domain's IP addresses
	ips, err := net.LookupIP(hostname)
	if err != nil {
		// If DNS resolution fails, we check if it is already a direct IP
		if ip := net.ParseIP(hostname); ip != nil {
			ips = []net.IP{ip}
		} else {
			return errors.New("could not resolve hostname")
		}
	}

	for _, ip := range ips {
		if ip.IsLoopback() {
			return errors.New("connections to loopback addresses are rejected")
		}
		if ip.IsPrivate() {
			return errors.New("connections to private IP addresses are rejected")
		}
		if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return errors.New("connections to link-local IP addresses are rejected")
		}
		if ip.IsUnspecified() {
			return errors.New("connections to unspecified IP addresses are rejected")
		}
	}

	return nil
}
