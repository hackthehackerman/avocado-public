package dao

import (
	"database/sql"
	"time"
)

func (d *Dao) GetUserById(userId string, db DBX) (user *User, err error) {
	result := []User{}
	if err = db.Select(&result, "select * from user where id = ?", userId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if len(result) == 0 {
		return nil, err
	}

	return &result[0], nil
}

func (d *Dao) GetUserByEmail(email string, db DBX) (user *User, err error) {
	result := []User{}
	if err = db.Select(&result, "select * from user where email = ?", email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if len(result) == 0 {
		return nil, err
	}

	return &result[0], nil
}

func (d *Dao) UpdateUser(user *User, db DBX) (err error) {
	updatedAt := time.Now().Unix()
	if _, err := db.Exec("update user set first_name=?, last_name=?, email=?, created_at=?, updated_at=? where id =?",
		user.FirstName, user.LastName, user.Email, user.CreatedAt, updatedAt, user.Id); err != nil {
		return err
	}
	return
}

func (d *Dao) InsertUser(user *User, db DBX) (err error) {
	updatedAt := time.Now().Unix()
	if _, err := db.Exec("INSERT INTO user VALUES(?,?,?,?,?,?)",
		user.Id, user.FirstName, user.LastName, user.Email, user.CreatedAt, updatedAt); err != nil {
		return err
	}
	return
}
