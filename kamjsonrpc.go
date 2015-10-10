/*
Released under MIT License <http://www.opensource.org/licenses/mit-license.php
Copyright (C) ITsysCOM GmbH. All Rights Reserved.

Provides simple Kamailio JSON-RPC over HTTP communication.
*/

package kamjsonrpc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type KamJsonRpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      uint64        `json:"id"`
}

type KamJsonRpcResponse struct {
	Jsonrpc string           `json:"jsonrpc"`
	Id      uint64           `json:"id"`
	Result  *json.RawMessage `json:"result"`
	Error   error            `json:"error"`
}

func NewKamailioJsonRpc(url string, skipTlsVerify bool) (*KamailioJsonRpc, error) {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTlsVerify}}}
	return &KamailioJsonRpc{url: url, client: client, mutex: new(sync.Mutex)}, nil
}

type KamailioJsonRpc struct {
	url    string
	client *http.Client
	id     uint64
	mutex  *sync.Mutex
}

// Generic function to remotely call a method and pass the parameters
func (self *KamailioJsonRpc) Call(serviceMethod string, args interface{}, reply *json.RawMessage) error {
	self.mutex.Lock()
	reqId := self.id
	self.id += 1
	self.mutex.Unlock()
	req := &KamJsonRpcRequest{Jsonrpc: "2.0", Method: serviceMethod, Id: reqId}
	if argSlice, isSlice := args.([]string); isSlice {
		req.Params = make([]interface{}, len(argSlice))
		for idx, val := range argSlice {
			req.Params[idx] = val
		}
	} else {
		req.Params = []interface{}{args}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := self.client.Post(self.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var kamResponse KamJsonRpcResponse
	if err = json.Unmarshal(respBody, &kamResponse); err != nil {
		return err
	}
	if kamResponse.Error != nil {
		return kamResponse.Error
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("Unexpected status code received: %d", resp.StatusCode)
	}
	if kamResponse.Id != reqId {
		return fmt.Errorf("Unsynchronized request, had: %d, received: %d", reqId, kamResponse.Id)
	}
	*reply = *kamResponse.Result
	return nil
}

// Add inidividual methods over the generic one
type RegistrationInfo struct {
	LocalUuid      string `json:"l_uuid"`
	LocalUsername  string `json:"l_username"`
	LocalDomain    string `json:"l_domain"`
	RemoteUsername string `json:"r_username"`
	RemoteDomain   string `json:"r_domain"`
	Realm          string `json:"realm"`
	AuthUsername   string `json:"auth_username"`
	AuthPassword   string `json:"auth_password"`
	AuthProxy      string `json:"auth_proxy"`
	Expires        int64  `json:"expires"`
	Flags          int64  `json:"flags"`
	DiffExpires    int64  `json:"diff_expires"`
	TimerExpires   int64  `json:"timer_expires"`
}

func (self *KamailioJsonRpc) UacRegInfo(params []string, reply *RegistrationInfo) error {
	var regRaw json.RawMessage
	if err := self.Call("uac.reg_info", params, &regRaw); err != nil {
		return err
	}
	return json.Unmarshal(regRaw, reply)
}
