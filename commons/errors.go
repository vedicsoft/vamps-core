package commons

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Error   error
	Message string
	Code    int
}

type ErrorHandler func(http.ResponseWriter, *http.Request) *AppError

func (fn ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		fmt.Println(">>>>>>>>>>>>>>>>>")
		http.Error(w, e.Message, e.Code)
	}
}
