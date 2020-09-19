package welcome_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mail-sending/handlers/welcome"
	"mail-sending/helpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

//SendWelcomeMailService is a mock Mail Service Interface
type fakeWelcomeMailService struct {
	SendWelcomeMailFn func(cred *helpers.WelcomeMail) (bool, error)
}

func (u *fakeWelcomeMailService) SendWelcomeMail(cred *helpers.WelcomeMail) (bool, error) {
	return u.SendWelcomeMailFn(cred)
}

var (
	fakeWelcomeMail fakeWelcomeMailService

	w = welcome.NewWelcome(&fakeWelcomeMail)
)

type MailResponse struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

//We dont need to mock the email layer, because we will never get there.
func TestNewWelcome_WrongInput(t *testing.T) {

	inputJSON := `{"name": "Victor Steven", "email": "wrongInput"}`

	req, err := http.NewRequest(http.MethodPost, "/welcome", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	r := gin.Default()
	r.POST("/welcome", w.WelcomeMail)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var response = MailResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}

	assert.EqualValues(t, rr.Code, 400)
	assert.EqualValues(t, response.Status, 400)
	assert.EqualValues(t, response.Body, "The email should be a valid email")
}


func TestNewWelcome_Success(t *testing.T) {

	fakeWelcomeMail.SendWelcomeMailFn = func(cred *helpers.WelcomeMail) (bool, error) {
		return true, nil
	}

	inputJSON := `{"name": "Victor Steven", "email": "victor@example.com"}`

	req, err := http.NewRequest(http.MethodPost, "/welcome", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	r := gin.Default()
	r.POST("/welcome", w.WelcomeMail)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var response = MailResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}

	assert.EqualValues(t, rr.Code, 200)
	assert.EqualValues(t, response.Status, 200)
	assert.EqualValues(t, response.Body, "Please check your mail")
}

func TestNewWelcome_Failure(t *testing.T) {

	fakeWelcomeMail.SendWelcomeMailFn = func(cred *helpers.WelcomeMail) (bool, error) {
		return false, errors.New("something went wrong sending mail")
	}

	inputJSON := `{"name": "Victor Steven", "email": "victor@example.com"}`

	req, err := http.NewRequest(http.MethodPost, "/welcome", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	r := gin.Default()
	r.POST("/welcome", w.WelcomeMail)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var response = MailResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}

	assert.EqualValues(t, rr.Code, 422)
	assert.EqualValues(t, response.Status, 422)
	assert.EqualValues(t, response.Body, "something went wrong sending mail")
}

