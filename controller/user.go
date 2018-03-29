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


	//parse json data from request
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)

	name := params["name"]
	if name == "" {
		log.Printf("name is null")
		return
	}

	user := &model.User{
		Name: name,
	}
	err := db.InsertUser(user)
	if err != nil {
		log.Printf("failed to insert user, err: %+v", err)
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

	path := strings.Split(r.URL.Path, "/")
	if len(path) != 4 {
		log.Println("request url error", r.URL.Path)
	}

	userId, err := strconv.ParseInt(path[2], 10, 64)
	if err != nil {
		log.Printf("failed to parse user id, userId: %s", path[1])
		return
	}

	user, err := db.GetUserByUid(userId)
	if err != nil {
		log.Printf("failed to get user by uid, %d", userId)
		return
	}

	res := user.GetAllRelationShips()

	json.NewEncoder(w).Encode(res)
}

// update relationship
func UpdateRelationShip(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer log.Printf("request: %s, execution time: %+v", r.RequestURI, time.Since(start))

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
		return
	}

	suid, err := strconv.ParseInt(path[2], 10, 64)
	if err != nil {
		log.Printf("failed to parse user id, userId: %s", path[2])
		return
	}

	a, err := db.GetUserByUid(suid)
	if err != nil {
		log.Printf("failed to get source user by id: %d, err: %+v", suid, err)
		return
	}

	tuid, err := strconv.ParseInt(path[4], 10, 64)
	if err != nil {
		log.Printf("failed to parse user id, userId: %s", path[4])
		return
	}

	b, err := db.GetUserByUid(tuid)
	if err != nil {
		log.Printf("failed to get target user by id: %d, err: %+v", tuid, err)
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
	curA2B := a.GetRelationShipByUid(b.Id)

	if curA2B != util.LIKE && curA2B != util.MATCH { //should update
		if curA2B != "" { //if disliked earlier, clear it
			a.ClearDislike(b.Id)
		}

		curB2A := b.GetRelationShipByUid(a.Id)
		//update relationship

		if curB2A == util.LIKE { //like each other
			a.AddMatch(b.Id)

			b.ClearLike(a.Id) //clear original like relationship
			b.AddMatch(a.Id)

			//todo do transaction??
			err := db.UpdateUser(a)
			if err != nil {
				log.Printf("failed to update relationship to match, suid: %d, tuid: %d, err: %+v", a.Id, b.Id, err)
			}
			err = db.UpdateUser(b)
			if err != nil {
				log.Printf("failed to update relationship to match, suid: %d, tuid: %d, err: %+v", b.Id, a.Id, err)
			}
		} else { //only A like B
			a.AddLike(b.Id)
			err := db.UpdateUser(a)
			if err != nil {
				log.Printf("failed to update relationship to like, suid: %d, tuid: %d, err: %+v", a.Id, b.Id, err)
			}
		}
	}
	return map[string]interface{}{
		"user_id": b.Id,
		"state":   a.GetRelationShipByUid(b.Id),
		"type":    "relationship",
	}
}

func ADislikeB(a, b *model.User) map[string]interface{} {
	curA2B := a.GetRelationShipByUid(b.Id)

	if curA2B != util.DISLIKE { //stay the same
		//if there is no relationship before
		c := make(chan struct{})
		go func() { //start a new goroutine to update relationship
			switch curA2B {
			case util.LIKE:
				a.ClearLike(b.Id)
				a.AddDislike(b.Id)
				err := db.UpdateUser(a)
				if err != nil {
					log.Printf("failed to clear like, suid: %d, tuid: %d, err: %+v", a.Id, b.Id, err)
				}
			case util.MATCH:
				a.ClearMatch(b.Id)
				a.AddDislike(b.Id)
				//todo do transaction??
				err := db.UpdateUser(a)
				if err != nil {
					log.Printf("failed to clear match relationship, suid: %d, tuid: %d, err: %+v", a.Id, b.Id, err)
				}

				b.ClearMatch(a.Id)
				b.AddLike(a.Id)
				err = db.UpdateUser(b)
				if err != nil {
					log.Printf("failed to clear match relationship, suid: %d, tuid: %d, err: %+v", b.Id, a.Id, err)
				}
			default:
				a.Dislikeset = append(a.Dislikeset, b.Id)
			}
			c <- struct{}{}
		}()
		<-c
	}

	return map[string]interface{}{
		"user_id": b.Id,
		"state":   util.DISLIKE,
		"type":    "relationship",
	}
}
