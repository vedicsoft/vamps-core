package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/vedicsoft/vamps-core/commons"
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

var redisProcess *exec.Cmd
var jwtResponse JWTResponse

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

//initializing the server for api tests
func setup() {
	serverConfigs := commons.InitConfigurations(commons.GetServerHome() + "/resources/.test/config.default.yaml")
	commons.ConstructConnectionPool(serverConfigs.DBConfigMap)

	//create the database in sqlite
	constructTestDB(serverConfigs.Home)
	redisProcess = startRedis(serverConfigs.Home)
	//constructing new routes
	m = routes.NewRouter()
	//The response recorder used to record HTTP responses
	respRec = httptest.NewRecorder()
}

func shutdown() {
	//deleting sqlite database
	err := os.Remove(commons.GetServerHome() + "/resources/.test/vampstest.db")
	if err != nil {
		fmt.Println("Unable to remove the test databsae stack:" + err.Error())
	}

	// Stopping redis process by reading the pid
	redisPid, err := ioutil.ReadFile(commons.GetServerHome() + "/resources/.test/redis.pid")
	err = redisProcess.Process.Kill()
	redisProcess.Wait()

	redisPid2, err := strconv.Atoi(string(redisPid)[:len(string(redisPid))-1])
	if err != nil {
		println(err.Error())
	}
	redisProcess2, _ := os.FindProcess(redisPid2)
	redisProcess2.Kill()
	redisProcess2.Wait()
	if err != nil {
		fmt.Println(err.Error())
	}
}

//takes the sqlite database descriptor and create a new one
func constructTestDB(serverHome string) {
	os.Chdir(serverHome + "/resources/.test")
	c1 := exec.Command("cat", "sqlite_serverdb.sql")
	c2 := exec.Command("./sqlite3", "vampstest.db")
	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()
	io.Copy(os.Stdout, &b2)
}

func startRedis(serverHome string) *exec.Cmd {
	os.Chdir(serverHome + "/resources/.test")
	// Starting caddy server to server static files
	args := []string{"redis.default.conf"}
	cmd := exec.Command("./redis-server", args...)
	if err := cmd.Start(); err != nil {
		fmt.Println("Error occourred while starting redis server : ", err.Error())
		os.Exit(1)
	}
	return cmd
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
	user := models.SystemUser{Username: "admin", Password: "admin", TenantDomain: "super.com"}
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
