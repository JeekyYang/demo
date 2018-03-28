package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Msg  string
	Code string
}

//get all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get all user")
	res := &Response{
		"success to get all users",
		"10000",
	}
	json.NewEncoder(w).Encode(res)
}

//create new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	res := &Response{
		"success to create user",
		"10000",
	}
	json.NewEncoder(w).Encode(res)
}

// get all relationship by uid
func GetAllRelationShipByUid(w http.ResponseWriter, r *http.Request) {
	res := &Response{
		"get all relationship",
		"10000",
	}
	json.NewEncoder(w).Encode(res)
}

// update relationship
func UpdateRelationShip(w http.ResponseWriter, r *http.Request) {
	//todo transaction
	res := &Response{
		"success to update relationship",
		"10000",
	}
	json.NewEncoder(w).Encode(res)

}
