package password

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"net/http"
	"strings"
)

func CheckCode(w http.ResponseWriter, r *http.Request) {
	code, ok := r.URL.Query()["code"]
	phone, ok := r.URL.Query()["phone"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"No code or phone provided"}, 400)
		return
	}

	getcode, ok := Storage[strings.TrimSpace(phone[0])]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Wrong phone provided!"}, 400)
		return
	}
	if getcode != code[0] {
		w.WriteHeader(http.StatusUnauthorized)
		WriteAnswer(&w, nil, []string{"Wrong code provided"}, 401)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, "Code is correct", []string{}, 200)

}
