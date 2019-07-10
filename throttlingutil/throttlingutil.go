package throttlingutil

import (
	"time"
)

type throttling struct {
	RegisteredVisits []registeredVisit
	BannedIPs        []IP
}

var theThrottling throttling

func init() {
	theThrottling = throttling{make([]registeredVisit, 0), make([]IP, 0)}
}

type Score int

type IP string

type registeredVisit struct {
	IP           IP
	Score        Score
	RegisteredAt time.Time
}

func (t *throttling) IsBanned(ip IP) bool {
	for _, bannedIp := range theThrottling.BannedIPs {
		if bannedIp == ip {
			return true
		}
	}
	return false
}

func (t *throttling) ban(ip IP) {
	t.BannedIPs = append(t.BannedIPs, ip)
}

func (t *throttling) getScore(ip IP, lookbackDuration time.Duration) (totalScore Score) {
	from := time.Now().UTC().Add(-lookbackDuration)
	for _, visit := range t.RegisteredVisits {
		if visit.RegisteredAt.After(from) {
			totalScore += visit.Score
		}
	}
	return
}

// nolint: vet
func (t *throttling) clearAll() {
	t.BannedIPs = []IP{}
	t.RegisteredVisits = []registeredVisit{}
}
