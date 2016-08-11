package api_test

import (
	"os"
	"testing"
	"github.com/vedicsoft/vamps-core/commons"
	"net/http"
	"strings"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"github.com/vedicsoft/vamps-core/routes"
	"encoding/json"
)

var m *mux.Router
var req *http.Request
var respRec *httptest.ResponseRecorder

//initializing the server for api tests
func setup() {
	//os.Setenv("SERVER_HOME","/home/anuruddha/workspace/go/src/github.com/vamps-console/server")
	os.Chdir(commons.ServerConfigurations.Home)
	commons.ConstructConnectionPool(commons.ServerConfigurations.DBConfigMap)
	m = routes.NewRouter()
	//The response recorder used to record HTTP responses
	respRec = httptest.NewRecorder()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	//shutdown()
	os.Exit(code)
}

func TestLogin(t *testing.T) {
	user := commons.SystemUser{Username:"admin", Password:"admin", TenantDomain:"super.com"}
	b, err := json.Marshal(user)
	req, err = http.NewRequest("POST", "/api/login", strings.NewReader(string(b)))
	if err != nil {
		t.Fatal("Creating 'POST /questions/1/SC' request failed!")
	}

	m.ServeHTTP(respRec, req)
	if respRec.Code != http.StatusOK {
		//TestDeleteUser(t)
		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusBadRequest)
	}
	t.Log(respRec.Body)
}