package iputil

import (
	"net"

	"github.com/alexandervantrijffel/goutil/logging"
)

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
