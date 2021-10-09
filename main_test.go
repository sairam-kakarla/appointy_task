package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestUserGet(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:8080/users?id=616138bf1581728515bbc4a5", nil)
	writer := httptest.NewRecorder()
	userGETHandler(writer, req)
	t.Errorf("working")
	resp := writer.Result()
	defer resp.Body.Close()
	var user User
	err := json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	var emailCheck, passwordCheck, nameCheck bool

	emailCheck = (user.Email == "rdfoff@gmail.com")
	passwordCheck = (user.Password == "57aaefa3f13fcbdf27b6a06b21a9383bba4d03c6c2d1074c806513b8c8dd1fb0")
	nameCheck = (user.Name == "RDJ_uno")
	if !emailCheck {
		t.Errorf("expected rdfoff@gmail.com, got %v", user.Email)
	}
	if !passwordCheck {
		t.Errorf("expected 57aaefa3f13fcbdf27b6a06b21a9383bba4d03c6c2d1074c806513b8c8dd1fb0, got %v", user.Password)
	}
	if !nameCheck {
		t.Errorf("expected RDJ_uno, got %v", user.Name)
	}

}
