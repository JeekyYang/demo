package model

import "log"

type User struct {
	Id         int64
	Likeset    []int64
	Dislikeset []int64
	Match      []int64
	Name       string
}

//todo: add unit test
//get all users
func (m *Model) selectAll() (error, []*User) {
	return nil, nil
}

//insert one user
func (m *Model) insert(user *User) error {
	err := m.DB.Insert(user)
	if err != nil {
		log.Fatalf("failed to insert User, %+v", user)
	}
	return err
}

//update user info
func (m *Model) update(user *User) error {
	return nil
}
