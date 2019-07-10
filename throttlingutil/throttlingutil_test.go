package throttlingutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterVisitWith16FailedVisitsShouldBanImmediately(t *testing.T) {
	theThrottling.clearAll()
	ip := IP("1.2.3.4")
	assert.Nil(t, registerVisit([]IP{ip}, 16*FAILEDVISITSCORE))
	assert.True(t, theThrottling.IsBanned(ip))
}

func TestRegisterVisitWith16FailedVisitsShouldNotBanLocalHost(t *testing.T) {
	theThrottling.clearAll()
	ip1 := IP("127.0.0.1")
	ip2 := IP("::1")
	ip3 := IP(" ")
	ips := []IP{ip1, ip2, ip3}
	err := registerVisit(ips, 16*FAILEDVISITSCORE)
	assert.NotNil(t, err)
	assert.False(t, theThrottling.IsBanned(ip1))
	assert.False(t, theThrottling.IsBanned(ip2))
	assert.False(t, theThrottling.IsBanned(ip3))
	assert.False(t, AreBanned([]string{string(ip1), string(ip2), string(ip3)}))
}

func TestRegisterVisitWith15FailedVisitsShouldNotBeBanned(t *testing.T) {
	theThrottling.clearAll()
	ip := []string{"1.2.3.4"}
	for i := 0; i < 15; i++ {
		RegisterFailedVisit(ip)
	}
	assert.False(t, theThrottling.IsBanned(IP(ip[0])))
}
func TestRegisterVisitWith16FailedVisitsShouldBanImmediately(t *testing.T) {
	theThrottling.clearAll()
	ip := IP("1.2.3.4")
	assert.Nil(t, registerVisit([]IP{ip}, 16*FAILEDVISITSCORE))
	assert.True(t, theThrottling.IsBanned(ip))
}

func TestRegisterVisitWith16FailedVisitsShouldNotBanLocalHost(t *testing.T) {
	theThrottling.clearAll()
	ip1 := IP("127.0.0.1")
	ip2 := IP("::1")
	ip3 := IP(" ")
	ips := []IP{ip1, ip2, ip3}
	err := registerVisit(ips, 16*FAILEDVISITSCORE)
	assert.NotNil(t, err)
	assert.False(t, theThrottling.IsBanned(ip1))
	assert.False(t, theThrottling.IsBanned(ip2))
	assert.False(t, theThrottling.IsBanned(ip3))
	assert.False(t, AreBanned([]string{string(ip1), string(ip2), string(ip3)}))
}

func TestRegisterVisitWith15FailedVisitsShouldNotBeBanned(t *testing.T) {
	theThrottling.clearAll()
	ip := []string{"1.2.3.4"}
	for i := 0; i < 15; i++ {
		RegisterFailedVisit(ip)
	}
	assert.False(t, theThrottling.IsBanned(IP(ip[0])))
}
