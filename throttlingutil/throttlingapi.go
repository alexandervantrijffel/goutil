package throttlingutil

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/alexandervantrijffel/goutil/errorcheck"
	"github.com/alexandervantrijffel/goutil/logging"
)

const FAILEDVISITSCORE = -5

func IsBanned(ip string) bool {
	return theThrottling.IsBanned(IP(ip))
}
func AreBanned(ips []string) bool {
	filteredIPs := excludeLocalhost(toIPs(ips))
	for _, ip := range filteredIPs {
		if theThrottling.IsBanned(ip) {
			return true
		}
	}
	return false
}

func excludeLocalhost(ips []IP) (filteredIPs []IP) {
	for _, ip := range ips {
		if len(strings.TrimSpace(string(ip))) > 0 && ip != "::1" && ip != "127.0.0.1" && ip != "?" {
			filteredIPs = append(filteredIPs, ip)
		}
	}
	return
}

func RegisterFailedVisit(ips []string) {
	_ = errorcheck.CheckLogf(registerVisit(toIPs(ips), FAILEDVISITSCORE), "RegisterVisit")
}

func toIPs(stringIPs []string) (ips []IP) {
	ips = make([]IP, unsafe.Sizeof(stringIPs))
	for i := range stringIPs {
		ips[i] = IP(stringIPs[i])
	}
	return
}

func registerVisit(ips []IP, score Score) error {
	filteredIPs := excludeLocalhost(ips)
	if len(filteredIPs) == 0 {
		return fmt.Errorf("No non-localhost ip's are in the list so it is ignored. ips: %+v", ips)
	}
	for _, ip := range filteredIPs {
		addNewRegisteredVisit(ip, score)
	}
	return nil
}

func addNewRegisteredVisit(ip IP, score Score) {
	if theThrottling.IsBanned(ip) {
		return
	}
	theThrottling.RegisteredVisits = append(theThrottling.RegisteredVisits,
		registeredVisit{ip, score, time.Now().UTC()})

	scoreLastTwoMinutes := theThrottling.getScore(ip, time.Minute*4)
	if scoreLastTwoMinutes <= Score(16*FAILEDVISITSCORE) {
		if string(ip) != "1.2.3.4" {
			logging.Errorf("Banning ip %s because it tried to open 16 invalid URLs within the last 4 minutes", ip)
		}
		theThrottling.ban(ip)
	}

	if theThrottling.getScore(ip, time.Hour*24) <= Score(40*FAILEDVISITSCORE) {
		if string(ip) != "1.2.3.4" {
			logging.Errorf("Banning ip %s because it tried to open 40 invalid URLs within the last day", ip)
		}
		theThrottling.ban(ip)
	}
}
