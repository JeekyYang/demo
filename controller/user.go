package controller

import (
	"demo/model"
	"demo/util"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Msg  string
	Code string
}

var db *model.Model

func init() {
	db = model.NewDB()
}

//get all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer log.Printf("request: %s, execution time: %+v", r.RequestURI, time.Since(start))

	users, err := db.GetAllUsers()
	if err != nil {
		log.Printf("failed to get users")
	}

	res := make([]map[string]interface{}, 0)

	for _, user := range users {
		res = append(res, map[string]interface{}{
			"id":   user.Id,
			"name": user.Name,
			"type": "user",
		})
	}

	json.NewEncoder(w).Encode(res)
}

//create new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer log.Printf("request: %s, execution time: %+v", r.RequestURI, time.Since(start))

	errRes := &Response{"internal error", "10000"}

	//parse json data from request
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)

	name := params["name"]
	if name == "" {
		log.Printf("name is null")
		json.NewEncoder(w).Encode(errRes)
		return
	}

	user := &model.User{
		Name: name,
	}
	err := db.InsertUser(user)
	if err != nil {
		log.Printf("failed to insert user, err: %+v", err)
		json.NewEncoder(w).Encode(errRes)
		return
	}

	res := map[string]interface{}{
		"id":   user.Id,
		"name": user.Name,
		"type": "user",
	}

	json.NewEncoder(w).Encode(res)

}

// get all relationship by uid
func GetAllRelationShipByUid(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer log.Printf("request: %s, execution time: %+v", r.RequestURI, time.Since(start))

	errRes := &Response{"internal error", "10000"}

	path := strings.Split(r.URL.Path, "/")
	if len(path) != 4 {
		log.Println("request url error", r.URL.Path)
	}

	userId, err := strconv.ParseInt(path[2], 10, 64)
	if err != nil {
		log.Printf("failed to parse user id, userId: %s", path[1])
		json.NewEncoder(w).Encode(errRes)
		return
	}

	user, err := db.GetUserByUid(userId)
	if err != nil {
		log.Printf("failed to get user by uid, %d", userId)
		json.NewEncoder(w).Encode(errRes)
		return
	}


	rs, err := db.GetAllRelationshipBySuid(user.Id)
	if err != nil {
		log.Printf("failed to get all relationship by suid: %d, err: %+v", user.Id, err)
	}
	res := make([]map[string]interface{}, len(rs))
	for _, r := range rs {
		res = append(res, map[string]interface{}{
			"user_id": r.Tuid,
			"state": r.State,
			"type": "relationship",
		})
	}

	json.NewEncoder(w).Encode(res)
}

// update relationship
func UpdateRelationShip(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer log.Printf("request: %s, execution time: %+v", r.RequestURI, time.Since(start))

	errRes := &Response{"internal error", "10000"}
	//parse urI
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 5 {
		log.Fatalln("request url error", r.URL.Path)
	}

	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)

	state := params["state"]
	if state != util.LIKE && state != util.DISLIKE {
		log.Printf("state error, state: %s", state)
		json.NewEncoder(w).Encode(errRes)
		return
	}

	suid, err := strconv.ParseInt(path[2], 10, 64)
	if err != nil {
		log.Printf("failed to parse user id, userId: %s", path[2])
		json.NewEncoder(w).Encode(errRes)
		return
	}

	a, err := db.GetUserByUid(suid)
	if err != nil {
		log.Printf("failed to get source user by id: %d, err: %+v", suid, err)
		json.NewEncoder(w).Encode(errRes)
		return
	}

	tuid, err := strconv.ParseInt(path[4], 10, 64)
	if err != nil {
		log.Printf("failed to parse user id, userId: %s", path[4])
		json.NewEncoder(w).Encode(errRes)
		return
	}

	b, err := db.GetUserByUid(tuid)
	if err != nil {
		log.Printf("failed to get target user by id: %d, err: %+v", tuid, err)
		json.NewEncoder(w).Encode(errRes)
		return
	}

	res := make(map[string]interface{})
	switch state {
	case util.LIKE:
		res = ALikeB(a, b) // a like b
	case util.DISLIKE:
		res = ADislikeB(a, b) // a dislike b
	}

	json.NewEncoder(w).Encode(res)
}

func ALikeB(a, b *model.User) map[string]interface{} {
	//get relationship
	curA2B, err := db.GetRelationshipBySuidAndTuid(a.Id, b.Id)
	if err != nil {
		log.Printf("failed to get relationship by suid: %d, tuid: %d, err: %+v", a.Id, b.Id, err)
		curA2B = &model.Relationship{Suid: a.Id, Tuid: b.Id, State: util.LIKE}
	}
	curA2B.State = util.LIKE

	if curA2B.State == util.LIKE || curA2B.State == util.MATCH {
		if curA2B.Id <= 0 {
			if err := db.InsertRelationship(curA2B); err != nil {
				log.Printf("failed to insert new relationship when like, r: %+v, err: %+v", curA2B, err)
			}
		}
	}
	curB2A, err := db.GetRelationshipBySuidAndTuid(b.Id, a.Id)
	if err != nil || curB2A == nil {
		log.Printf("failed to get relationship by suid: %d, tuid: %d, err: %+v", b.Id, a.Id, err)
	}

	if curB2A.State == util.LIKE {
		curB2A.State = util.MATCH
		curA2B.State = util.MATCH

		db.UpdateRelationship(curB2A)
	}

	log.Printf("cur relationship: %s", curA2B.State)
	db.UpdateRelationship(curA2B)
	return map[string]interface{}{
		"user_id": b.Id,
		"state":   curA2B.State,
		"type":    "relationship",
	}
}

func ADislikeB(a, b *model.User) map[string]interface{} {
	curA2B, err := db.GetRelationshipBySuidAndTuid(a.Id, b.Id)
	if err != nil {
		log.Printf("failed to get relationship by suid: %d, tuid: %d, err: %+v", a.Id, b.Id, err)
		curA2B = &model.Relationship{Suid: a.Id, Tuid: b.Id, State: util.DISLIKE}
	}

	switch curA2B.State {
	case util.DISLIKE:
		if curA2B.Id <= 0 {
			if err := db.InsertRelationship(curA2B); err != nil {
				log.Printf("failed to insert new relationship when dislike, r: %+v, err: %+v", curA2B, err)
			}
		}
	case util.LIKE, util.MATCH:
		curA2B.State = util.DISLIKE
		db.UpdateRelationship(curA2B)

		curB2A, err := db.GetRelationshipBySuidAndTuid(b.Id, a.Id)
		if err != nil {
			log.Printf("failed to get relationship by suid: %d, tuid: %d, err: %+v", b.Id, a.Id, err)
			break
		}
		if curB2A.State == util.MATCH {
			curB2A.State = util.LIKE
			db.UpdateRelationship(curB2A)
		}
	}

	return map[string]interface{}{
		"user_id": b.Id,
		"state":   util.DISLIKE,
		"type":    "relationship",
	}
}
