package middlewares_test

import (
	"fmt"
	"mxshop_api/user-web/middlewares"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	jwt := middlewares.NewTestJWT()
	token, err := jwt.RefreshToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6OSwiTmlja05hbWUiOiJseW5uOCIsIkF1dGhvcml0eUlkIjoxLCJleHAiOjE2MzQ4MTU3NTIsImlzcyI6Ikx5bm4iLCJuYmYiOjE2MzQ4MTIxNTJ9.4LNXZsLAvDDZ7ICrUfrkgyqVYGhrrGueAmBpY9aXHDs")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)
}
