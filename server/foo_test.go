package server

import (
	"encoding/json"
	"testing"

	"github.com/eremetic-framework/eremetic/api"
)

var s = `{"image":"alpine","command":"ls","cpu":1,"mem":100,"callback_uri":"","fetch":[{"uri":"file:///etc/mesos/docker.cfg"}]}`

func TestReqToTask(t *testing.T) {
	var req api.RequestV1
	err := json.Unmarshal([]byte(s), &req)
	if err != nil {
		t.Fatal(err)
	}
	request := api.RequestFromV1(req)

	if len(request.URIs) != 1 {
		t.Errorf("No URIs found, got  %v", request)
	}
}
