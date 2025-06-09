package main

import (
	"errors"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
)

func sessionsOperate(resp http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session-id")
	if errors.Is(err, http.ErrNoCookie) {
		sessionId, _ := uuid.NewV4()
		http.SetCookie(resp, &http.Cookie{
			Name:  "session-id",
			Value: sessionId.String(),
		})
	}
	fmt.Printf("fetch cookies is %s", cookie)
}
func main() {
	http.HandleFunc("/", sessionsOperate)
	http.ListenAndServe(":8080", nil)
}
