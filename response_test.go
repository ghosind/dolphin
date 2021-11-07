package dolphin

import (
	"net/http"
	"testing"
)

func TestReset(t *testing.T) {
	resp := &Response{}
	resp.reset()

	if len(resp.body.Bytes()) != 0 {
		t.Errorf("Response body expect 0, actual %d", len(resp.body.Bytes()))
	}

	if len(resp.cookies) != 0 {
		t.Errorf("Length of response cookies expect 0, actual %d", len(resp.cookies))
	}

	if len(resp.header) != 0 {
		t.Errorf("Length of response header expect 0, actual %d", len(resp.header))
	}

	if resp.statusCode != http.StatusOK {
		t.Errorf("Response status code expect %d, actual %d",
			http.StatusOK,
			resp.statusCode,
		)
	}
}

func TestSetStatus(t *testing.T) {
	resp := &Response{}
	resp.reset()

	err := resp.SetStatusCode(-1)
	if err == nil {
		t.Errorf("Set status code as -1 expect return error, actual return nil")
	}
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Set status code as -1 expect resp.statusCode is 200%d, actual %d",
			http.StatusOK,
			resp.statusCode,
		)
	}

	err = resp.SetStatusCode(0)
	if err == nil {
		t.Errorf("Set status code as 0 expect return error, actual return nil")
	}
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Set status code as 0 expect resp.statusCode is %d, actual %d",
			http.StatusOK,
			resp.statusCode,
		)
	}

	err = resp.SetStatusCode(1000)
	if err == nil {
		t.Errorf("Set status code as 1000 expect return error, actual return nil")
	}
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Set status code as 1000 expect resp.statusCode is %d, actual %d",
			http.StatusOK,
			resp.statusCode,
		)
	}

	err = resp.SetStatusCode(http.StatusCreated)
	if err != nil {
		t.Errorf("Set status code as %d expect return nil, actual return %v",
			http.StatusCreated,
			err,
		)
	}
	if resp.StatusCode() != http.StatusCreated {
		t.Errorf("Set status code as %d expect resp.statusCode is %d, actual %d",
			http.StatusCreated,
			http.StatusCreated,
			resp.statusCode,
		)
	}
}
