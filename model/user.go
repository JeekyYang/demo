package model

import (
	"database/sql"
)

type User struct {
	Id         int64   `sql:",pk"`
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