package iputil

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/alexandervantrijffel/goutil/errorcheck"
	"github.com/alexandervantrijffel/goutil/logging"
)

func GetIP(req *http.Request) (remoteIP string) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		if !strings.Contains(err.Error(), "missing port in address") {
			_ = errorcheck.CheckLogf(err, "SplitHostPort failed, remoteIP:", remoteIP)
		}
		ip = req.RemoteAddr
	}

	ip = removeReverseProxyIP(ip)
	if IPIsValidAndRemote(ip) {
		remoteIP = ip
	} else {
		remoteIP = fmt.Sprintf("??? '%s'", ip)
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	forward := req.Header.Get("X-Forwarded-For")
	forward = removeReverseProxyIP(forward)
	if IPIsValidAndRemote(forward) {
		remoteIP = forward
	}
	return
}

func IPIsValidAndRemote(ip string) bool {
	if len(ip) == 0 {
		return false
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		if ip == "localhost" {
			return false
		}
		logging.Errorf("Could not parse IP '%s'. req.RemoteAddr is not correct format", ip)
		return false
	}
	return ip != "127.0.0.1" && ip != "::1" && ip != "localhost"
}

func removeReverseProxyIP(remoteIP string) string {
	// Stackpath sets remote addr as `VISITORADDR, STACKPATHPROXYADDR`
	parts := strings.Split(remoteIP, ",")
	if len(parts) > 0 {
		return parts[0]
	}
	return remoteIP
}
