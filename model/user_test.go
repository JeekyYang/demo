package model

import (
	"testing"
)

//add test for model User
func TestModel_GetAllUsers(t *testing.T) {
	m := NewDB()
	u0 := &User{Id: 1}
	u1 := &User{Id: 2}

	err := m.InsertUser(u0)
	if err != nil {
		t.Errorf("failed to insert user, %+v, %+v", u0, err)
	}

	err = m.InsertUser(u1)
	if err != nil {
		t.Errorf("failed to insert user, %+v, %+v", u1, err)
	}

	users, err := m.GetAllUsers()
	if len(users) != 2 {
		t.Errorf("total count required: 2, in fact: %d", len(users))
	}
}

func TestModel_GetUserByUid(t *testing.T) {
	m := NewDB()
	u0 := &User{Id: 1}

	err := m.InsertUser(u0)
	if err != nil {
		t.Errorf("failed to insert user, %+v, %+v", u0, err)
	}

	user, err := m.GetUserByUid(1)
	if err != nil || user.Id != 1 {
		t.Errorf("failed to get user by uid: %d", 1)
	}
}

func TestModel_InsertUser(t *testing.T) {
	m := NewDB()
	u0 := &User{Id: 1}
	u1 := &User{Id: 2}
	u2 := &User{Name: "test123"}

	err := m.InsertUser(u0)
	if err != nil {
		t.Errorf("failed to insert user, %+v, %+v", u0, err)
	}

	err = m.InsertUser(u1)
	if err != nil {
		t.Errorf("failed to insert user, %+v, %+v", u1, err)
	}

	//duplicated insert test
	err = m.InsertUser(u1)
	if err == nil {
		t.Errorf("duplicated insert should not be successful, %+v", err)
	}

	err = m.InsertUser(u2)
	if err != nil {
		t.Errorf("failed to insert user with name, %+v", err)
	}
}

func TestModel_UpdateUser(t *testing.T) {

	m := NewDB()

	u0 := &User{Id: 1, Name: "test"}

	err := m.InsertUser(u0)
	if err != nil {
		t.Errorf("failed to insert user, %+v, %+v", u0, err)
	}

	u0.Name = "test2"
	err = m.UpdateUser(u0)
	if err != nil {
		t.Errorf("failed to update user info, %+v", err)
	}

	u, err := m.GetUserByUid(1)
	if err != nil {
		t.Errorf("failed to get user by uid: %d, err: %+v", 1, err)
	}

	if u.Name != "test2" {
		t.Error("user info update fail")
	}
}
