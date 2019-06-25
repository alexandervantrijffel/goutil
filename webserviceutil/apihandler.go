package webserviceutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/alexandervantrijffel/goutil/errorcheck"
	"github.com/alexandervantrijffel/goutil/logging"
	"github.com/julienschmidt/httprouter"
)

type RequestHandler func() (responseData interface{}, err error)
type ApiHandler func(r *http.Request) (requestData interface{}, requestHandler RequestHandler)

func OuterApiHandler(getInnerHandler ApiHandler, handlerName string) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			HandlePanic(w, recover(), handlerName)
		}()

		requestData, requestHandler := getInnerHandler(r)
		err, orgErr := decodeJSON(requestData, r.Body)
		if orgErr != io.EOF && err != nil { // ignore error for empty request
			logging.Errorf("DecodeJSON failed for handler name %s with error %s. Data: %s",
				handlerName, err.Error(), readBodyPrefix(r.Body))
			return
		}
		results, err := requestHandler()
		errorcheck.CheckPanic(err)
		err = dumpJSON(w, results)
		errorcheck.CheckPanic(err)
	}
}
func readBodyPrefix(body io.ReadCloser) string {
	if s, err := ioutil.ReadAll(body); err != nil {
		return ""
	} else if len(s) > 512 {
		return string(s[0:512])
	} else {
		return string(s)
	}
}

func decodeJSON(destinationObject interface{}, body io.ReadCloser) (error, error) {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(destinationObject); err != nil {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(body)
		_ = errorcheck.CheckLogf(err, "Failed to read HTTP request body as string")
		bodySz := buf.String()
		if len(bodySz) > 512 {
			bodySz = bodySz[:512]
		}
		return fmt.Errorf("Failed to decode object from JSON %+v. Request body: %s",
			err, bodySz), err
	}
	return nil, nil
}

func dumpJSON(w http.ResponseWriter, zeData interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	// stream the contents of the file to the response
	return json.NewEncoder(w).Encode(zeData)
}
