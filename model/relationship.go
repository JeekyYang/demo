package model

import "database/sql"

type Relationship struct {
	Id int64
	Suid int64
	Tuid int64
	State string
}

func(m *Model) InsertRelationship(r *Relationship) error {
	return m.DB.Insert(r)
}

func(m *Model) UpdateRelationship(r *Relationship) error {
	return m.DB.Update(r)
}

func (m *Model) GetAllRelationshipBySuid(suid int64) ([]Relationship, error) {

	var rs []Relationship

	_, err := m.DB.Query(&rs, `select * from relationships where suid = ?`, suid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return rs, nil
}

func(m *Model) GetRelationshipBySuidAndTuid(suid, tuid int64) (*Relationship, error) {
	var r Relationship
	_, err := m.DB.QueryOne(&r, `select * from relationships where suid = ? and tuid = ?`, suid, tuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &r, nil
}
