package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/osamah22/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, string(body), "PONG")
}

func TestShowSnippet(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our
	// application for different URLs.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippets/view/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippets/view/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippets/view/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippets/view/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippets/view/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippets/view/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippets/view/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	// Create the application struct containing our mocked dependencies and set
	// up the test server for running and end-to-end test.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET /user/signup request and then extract the CSRF token from the
	// response body.
	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	t.Log(csrfToken)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Bob", "bob@foo.bar", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "bob@foo.bar", "validPa$$word", csrfToken, http.StatusUnprocessableEntity, []byte("This field cannot be blank")},
		{"Empty email", "Bob", "", "validPa$$word", csrfToken, http.StatusUnprocessableEntity, []byte("This field cannot be blank")},
		{"Empty password", "Bob", "bob@foo.bar", "", csrfToken, http.StatusUnprocessableEntity, []byte("This field cannot be blank")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", csrfToken, http.StatusUnprocessableEntity, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Bob", "bobfoo.bar", "validPa$$word", csrfToken, http.StatusUnprocessableEntity, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Bob", "@foo.bar", "validPa$$word", csrfToken, http.StatusUnprocessableEntity, []byte("This field is invalid")},
		{"Short password", "Bob", "bob@foo.bar", "pa$$", csrfToken, http.StatusUnprocessableEntity, []byte("This field is too short (minimum is 8 characters)")},
		{"Duplicate email", "Bob", "dupe@foo.bar", "validPa$$word", csrfToken, http.StatusUnprocessableEntity, []byte("Address is already in use")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.postForm(t, "/user/signup", form)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			// if !bytes.Contains(body, tt.wantBody) {
			// 	t.Errorf("want body %s to contain %q", body, tt.wantBody)
			// }
		})
	}
}
