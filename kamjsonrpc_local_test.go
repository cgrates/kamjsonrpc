package kamjsonrpc

import (
	"encoding/json"
	"flag"
	"reflect"
	"testing"
)

var testLocal = flag.Bool("local", false, "Perform the tests only on local test environment, not by default.") // This flag will be passed here via "go test -local" args
var kamAddr = flag.String("kam_addr", "http://127.0.0.1:5060", "Address where to reach kamailio http server")

var kamRpc *KamailioJsonRpc

func TestKamJsonRpcConn(t *testing.T) {
	if !*testLocal {
		return
	}
	var err error
	if kamRpc, err = NewKamailioJsonRpc(*kamAddr, true); err != nil {
		t.Fatal("Cannot connect to kamailio:", err)
	}
}

func TestKamJsonRpcCall(t *testing.T) {
	if !*testLocal {
		return
	}
	var reply json.RawMessage
	if err := kamRpc.Call("core.psx", []string{}, &reply); err != nil {
		t.Error(err)
	} else if reflect.DeepEqual(reply, json.RawMessage{}) {
		t.Error("Empty reply")
	}
}

func TestKamJsonRpcUacRegEnable(t *testing.T) {
	if !*testLocal {
		return
	}
	eReply := "OK"
	var reply string
	if err := kamRpc.UacRegEnable([]string{"l_uuid", "unknown"}, &reply); err != nil {
		t.Error(err)
	} else if eReply != reply {
		t.Errorf("Expecting: %s, received: %s", eReply, reply)
	}
}

func TestKamJsonRpcUacRegDisable(t *testing.T) {
	if !*testLocal {
		return
	}
	eReply := "OK"
	var reply string
	if err := kamRpc.UacRegDisable([]string{"l_uuid", "unknown"}, &reply); err != nil {
		t.Error(err)
	} else if eReply != reply {
		t.Errorf("Expecting: %s, received: %s", eReply, reply)
	}
}

func TestKamJsonRpcUacRegReload(t *testing.T) {
	if !*testLocal {
		return
	}
	eReply := "OK"
	var reply string
	if err := kamRpc.UacRegReload([]string{}, &reply); err != nil {
		t.Error(err)
	} else if eReply != reply {
		t.Errorf("Expecting: %s, received: %s", eReply, reply)
	}
}

func TestKamJsonRpcUacRegRefresh(t *testing.T) {
	if !*testLocal {
		return
	}
	eReply := "OK"
	var reply string
	if err := kamRpc.UacRegRefresh([]string{"l_uuid", "unknown"}, &reply); err != nil {
		t.Error(err)
	} else if eReply != reply {
		t.Errorf("Expecting: %s, received: %s", eReply, reply)
	}
}

func TestKamJsonRpcUacRegInfo(t *testing.T) {
	if !*testLocal {
		return
	}
	var eReply RegistrationInfo
	var reply RegistrationInfo
	if err := kamRpc.UacRegInfo([]string{"l_uuid", "unknown"}, &reply); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(reply, eReply) {
		t.Errorf("Expecting: %+v, received: %+v", eReply, reply)
	}
}
