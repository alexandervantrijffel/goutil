package webserviceutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/alexandervantrijffel/goutil/errorcheck"
	"github.com/alexandervantrijffel/goutil/iputil"
	"github.com/alexandervantrijffel/goutil/logging"
	"gitlab.com/avtnl/ps-737migration-be/workspace/throttling"
	"gitlab.com/avtnl/ps-737migration-be/workspace/util"
)

func ReturnText(w http.ResponseWriter, data string) {
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Accept, X-Access-Token, X-Application-Name, X-Request-Sent-Time, Accept-Encoding, X-Compress")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	_, err := w.Write([]byte(data))
	_ = errorcheck.CheckLogf(err, "Failed to return http messaage %s", data)
}
func Return(rw http.ResponseWriter, statusCode int, data map[string]interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(statusCode)
	_ = json.NewEncoder(rw).Encode(&data)
}

type StatusMessageResponse struct {
	Description string `json:"description"`
}

func ReturnStatusMessageResponse(rw http.ResponseWriter, statusCode int, description string) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(statusCode)
	_ = json.NewEncoder(rw).Encode(&StatusMessageResponse{description})
}
func ReturnFile(w http.ResponseWriter, r *http.Request, file []byte, filename string) {
	w.Header().Add("Content-Disposition", "inline")
	w.Header().Add("filename", filename)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Accept, X-Access-Token, X-Application-Name, X-Request-Sent-Time")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(file))
}
func GetIP(req *http.Request) (remoteIP string) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		_ = errorcheck.CheckLogf(err, "SplitHostPort failed, remoteIP:", remoteIP)
		ip = req.RemoteAddr
	}

	if iputil.IPIsValidAndRemote(ip) {
		remoteIP = ip
	} else {
		remoteIP = fmt.Sprintf("??? '%s'", ip)
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	forward := req.Header.Get("X-Forwarded-For")
	if util.IPIsValid(forward) {
		remoteIP = forward
	}
	return
}

func VerifyRemoteIPIsBannedAndReject(r *http.Request, rw http.ResponseWriter) (remoteAddress string, isBanned bool) {
	remoteAddress = GetIP(r)
	if throttling.IsBanned(remoteAddress) {
		ReplyUnauthorized(rw)
		isBanned = true
	}
	return
}

func ReplyUnauthorized(rw http.ResponseWriter) {
	ReturnStatusMessageResponse(rw, http.StatusUnauthorized, "You are not authorized for this action. Please check your link.")
}

func HandleUnauthorized(remoteAddress string, token string, method string, r *http.Request, rw http.ResponseWriter) {
	logging.Warningf("Could not authenticate user with token %s, ip %+v", token, GetIP(r))
	ReplyUnauthorized(rw)
	throttling.RegisterFailedVisit([]string{remoteAddress})
}

func toFullString(r *http.Request) string {
	return fmt.Sprintf("%s %s %s %s", r.RemoteAddr, r.Method, r.Host, r.RequestURI)
}

func ReplyInternalServerError(rw http.ResponseWriter, recoveredPanic interface{}, r *http.Request, description string) {
	logging.Error(description, recoveredPanic, toFullString(r), fmt.Sprintf("%s", debug.Stack()))
	ReturnStatusMessageResponse(rw, http.StatusInternalServerError, "Something went wrong on our end. "+description)
}

func HandlePanic(w http.ResponseWriter, r interface{}, handlerName string) {
	if r == nil {
		return
	}
	logging.Errorf("Unexpected failure in %s %+v", handlerName, r)
	ReturnStatusMessageResponse(w,
		http.StatusInternalServerError,
		fmt.Sprintf("Unexpected failure %+v", r))
}
