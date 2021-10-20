package dolphin

import (
	"bytes"
	"encoding/json"
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
			c.String("Method Should be GET", http.StatusBadRequest)
			return
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
		return
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

func TestGetRequestPostJSON(t *testing.T) {
	type PostJsonTestPayload struct {
		Name    *string `json:"name"`
		Message *string `json:"message"`
	}

	app := New(&Config{
		Port: 8081,
	})

	app.Use(func(c *Context) {
		if c.Method() != http.MethodPost {
			t.Errorf("Expect POST request, actual %s", c.Method())
			c.String("Method Should be POST", http.StatusBadRequest)
			return
		}

		var payload PostJsonTestPayload
		err := c.PostJSON(&payload)
		if err != nil {
			t.Log(err)
			c.String("Invalid parameter format", http.StatusBadRequest)
			return
		}

		c.JSON(O{
			"message": fmt.Sprintf("Hello %s", *payload.Name),
		})
	})

	go func() {
		app.Run()
	}()
	defer app.Shutdown()

	cli := http.Client{}

	body, err := json.Marshal(map[string]string{
		"name": "dolphin",
	})
	if err != nil {
		t.Errorf("json.Marshal expect ni error, actual %v", err)
		return
	}

	resp, err := cli.Post("http://localhost:8081", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("GET request expect no error, actual got %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET request expect status code %d, actual got %d", http.StatusOK, resp.StatusCode)
		return
	}

	var data PostJsonTestPayload
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		t.Errorf("Decode response body expect no error, actual got %v", err)
		return
	}

	if data.Message == nil || *data.Message != "Hello dolphin" {
		t.Errorf("Expect response body 'Hello dolphin', actual got %v", *data.Message)
	}
}
