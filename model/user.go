package model

import (
	"database/sql"
	"demo/util"
)

type User struct {
	Id         int64   `sql:",pk"`
	Likeset    []int64 `pg:",array"`
	Dislikeset []int64 `pg:",array"`
	Match      []int64 `pg:",array"`
	Name       string
}

//todo: add unit test
//get all users
func (m *Model) GetAllUsers() ([]User, error) {
	var users []User
	err := m.DB.Model(&users).Select()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return users, nil
}

//select user by uid
func (m *Model) GetUserByUid(uid int64) (*User, error) {
	user := &User{Id: uid}
	err := m.DB.Select(user)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return user, nil
}

//insert one user
func (m *Model) InsertUser(user *User) error {
	return m.DB.Insert(user)
}

//update user info
func (m *Model) UpdateUser(user *User) error {
	return m.DB.Update(user)
}

func (u *User) GetRelationShipByUid(uid int64) string {
	for _, id := range u.Likeset {
		if id == uid {
			return util.LIKE
		}
	}

	for _, id := range u.Dislikeset {
		if id == uid {
			return util.DISLIKE
		}
	}

	for _, id := range u.Match {
		if id == uid {
			return util.MATCH
		}
	}

	return ""
}

func (u *User) AddLike(uid int64) {
	u.Likeset = append(u.Likeset, uid)
}

func (u *User) ClearLike(uid int64) {
	newLikedSet := make([]int64, 0) //len(u.LikedSet)?
	for _, likedId := range u.Likeset {
		if likedId != uid {
			newLikedSet = append(newLikedSet, likedId)
		}
	}
	u.Likeset = newLikedSet
}

func (u *User) AddMatch(uid int64) {
	u.Match = append(u.Match, uid)
}

func (u *User) ClearMatch(uid int64) {
	newMatchSet := make([]int64, 0)
	for _, matchId := range u.Match {
		if matchId != uid {
			newMatchSet = append(newMatchSet, matchId)
		}
	}
	u.Match = newMatchSet
}

func (u *User) AddDislike(uid int64) {
	u.Dislikeset = append(u.Dislikeset, uid)
}

func (u *User) ClearDislike(uid int64) {
	newDislikeSet := make([]int64, 0)
	for _, dislikeId := range u.Dislikeset {
		if dislikeId != uid {
			newDislikeSet = append(newDislikeSet, dislikeId)
		}
	}
	u.Dislikeset = newDislikeSet
}

func (u *User) GetAllRelationShips() []map[string]interface{} {
	res := make([]map[string]interface{}, 0)

	for _, likeId := range u.Likeset {
		res = append(res, map[string]interface{}{
			"user_id": likeId,
			"state":   util.LIKE,
			"type":    "relationship",
		})
	}

	for _, dislikeId := range u.Dislikeset {
		res = append(res, map[string]interface{}{
			"user_id": dislikeId,
			"state":   util.DISLIKE,
			"type":    "relationship",
		})
	}

	for _, matchedId := range u.Match {
		res = append(res, map[string]interface{}{
			"user_id": matchedId,
			"state":   util.MATCH,
			"type":    "relationship",
		})
	}
	return res
}
