package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/vedicsoft/vamps-core/models"
	"github.com/vedicsoft/vamps-core/routes"
)

var m *mux.Router
var req *http.Request
var respRec *httptest.ResponseRecorder

type JWTResponse struct {
	Token    string
	TenantId int
}

var jwtResponse JWTResponse

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

//initializing the server for api tests
func setup() {
	//constructing new routes
	m = routes.NewRouter()
	//The response recorder used to record HTTP responses
	respRec = httptest.NewRecorder()
}

func shutdown() {

}

func TestLogin(t *testing.T) {
	user := models.SystemUser{Username: "admin@super.com", Password: "admin"}
	b, err := json.Marshal(user)
	req, err = http.NewRequest("POST", "/login", strings.NewReader(string(b)))
	if err != nil {
		t.Fatal("Creating 'POST /login' request failed!")
	}
	//The response recorder used to record HTTP responses
	respRec = httptest.NewRecorder()
	m.ServeHTTP(respRec, req)
	if respRec.Code != http.StatusOK {
		//TestDeleteUser(t)
		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusBadRequest)
	}
	t.Log(respRec.Body)
	decoder := json.NewDecoder(respRec.Body)
	err = decoder.Decode(&jwtResponse)
	if err != nil {
		t.Error("Error while decoding JWT token responce")
	}
	respRec.Flush()
}

func TestLogout(t *testing.T) {
	user := models.SystemUser{Username: "admin@super.com", Password: "admin"}
	b, err := json.Marshal(user)
	req, err = http.NewRequest("POST", "/logout", strings.NewReader(string(b)))
	if err != nil {
		t.Fatal("Creating 'POST /logout' request failed!")
	}
	req.Header.Set("Authorization", "Bearer "+jwtResponse.Token)
	respRec = httptest.NewRecorder()
	m.ServeHTTP(respRec, req)
	if respRec.Code != http.StatusOK {
		//TestDeleteUser(t)
		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusBadRequest)
	}
	t.Log(respRec.Body)
}
