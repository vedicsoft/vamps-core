package commons

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type AppError struct {
	Error   error
	Message string
	Code    int
}

func (z AppError) String() string {
	b, _ := json.Marshal(z)
	return string(b)
}

type AppHandler func(http.ResponseWriter, *http.Request) *AppError

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Error(e.String())
		http.Error(w, e.Message, e.Code)
	}
}
