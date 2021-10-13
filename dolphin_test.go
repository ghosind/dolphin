package dolphin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetRequestQuery(t *testing.T) {
	app := Default()

	app.Use(func(c *Context) {
		if c.Method() != http.MethodGet {
			t.Errorf("Expect GET request, actual %s", c.Method())
			c.Fail("Method Should be GET")
		}

		name := c.Query("name")
		c.String(fmt.Sprintf("Hello %s", name), http.StatusOK)
	})

	go func() {
		app.Run()
	}()

	cli := http.Client{}

	resp, err := cli.Get("http://localhost:8080?name=dolphin")
	defer app.Shutdown()
	if err != nil {
		t.Errorf("GET request expect no error, actual got %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET request expect status code %d, actual got %d", http.StatusOK, resp.StatusCode)
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Read from response body expect no error, actual got %v", err)
		return
	}

	if string(res) != "Hello dolphin" {
		t.Errorf("Expect response body 'Hello dolphin', actual got %v", string(res))
	}
}
