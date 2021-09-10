package dolphin

import (
	"net/http"
	"testing"
)

func TestSetStatus(t *testing.T) {
	var resp *Response = &Response{}
	resp.reset()

	err := resp.SetStatusCode(-1)
	if err == nil {
		t.Errorf("Set status code as -1 expect return error, actual return nil")
	}
	if resp.statusCode != 200 {
		t.Errorf("Set status code as -1 expect resp.statusCode is 200, actual %d", resp.statusCode)
	}

	err = resp.SetStatusCode(0)
	if err == nil {
		t.Errorf("Set status code as 0 expect return error, actual return nil")
	}
	if resp.statusCode != 200 {
		t.Errorf("Set status code as 0 expect resp.statusCode is 200, actual %d", resp.statusCode)
	}

	err = resp.SetStatusCode(1000)
	if err == nil {
		t.Errorf("Set status code as 1000 expect return error, actual return nil")
	}
	if resp.statusCode != 200 {
		t.Errorf("Set status code as 1000 expect resp.statusCode is 200, actual %d", resp.statusCode)
	}

	err = resp.SetStatusCode(http.StatusCreated)
	if err != nil {
		t.Errorf("Set status code as %d expect return nil, actual return %v", http.StatusCreated, err)
	}
	if resp.statusCode != http.StatusCreated {
		t.Errorf("Set status code as %d expect resp.statusCode is 200, actual %d", http.StatusCreated, resp.statusCode)
	}
}
