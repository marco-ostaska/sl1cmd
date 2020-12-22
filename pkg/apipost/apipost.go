// Package apipost post into sl1api
package apipost

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/marco-ostaska/sl1cli-tools/pkg/apicryptcfg"
	"github.com/marco-ostaska/sl1cli-tools/pkg/sl1generics"
)

// APIData an abstraction to API
type APIData struct {
	API     string          // API section as : /api/account
	ARGS    string          // any parameter need to be sent to api
	Payload *strings.Reader //payload for posting
	Result  []byte          // result from call
}

// APIPost make the http calls
func (a *APIData) APIPost() error {

	var uCFG apicryptcfg.UserInfo
	if err := uCFG.ReadCryptFile(); err != nil {
		return err
	}

	url := uCFG.URL + a.API + a.ARGS
	method := "POST"
	if err := sl1generics.IsReachable(uCFG.URL); err != nil {
		return fmt.Errorf("%s is unreachable", uCFG.URL)
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, a.Payload)

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+uCFG.B64)

	res, err := client.Do(req)
	defer func() {
		cerr := res.Body.Close()
		if err == nil {
			err = cerr
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	a.Result = body
	return nil
}

// NewPost make new post to sl1 API
func (a *APIData) NewPost(v interface{}, as ...string) error {
	a.API = as[0]
	if len(as) == 2 {
		a.ARGS = as[1]
	}

	if len(as) > 2 {
		a.Payload = strings.NewReader(as[2])
	}

	if err := a.APIPost(); err != nil {
		return err
	}

	return json.Unmarshal(a.Result, &v)

}

// APIDelete delete entry in api
func (a *APIData) APIDelete() error {

	var uCFG apicryptcfg.UserInfo
	if err := uCFG.ReadCryptFile(); err != nil {
		return err
	}

	url := uCFG.URL + a.API + a.ARGS
	method := "DELETE"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+uCFG.B64)

	res, err := client.Do(req)
	defer func() {
		cerr := res.Body.Close()
		if err == nil {
			err = cerr
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	a.Result = body
	return nil
}
