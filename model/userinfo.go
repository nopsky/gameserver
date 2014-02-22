package model

import (
	"database/sql"
	"errors"
	"log"
)

type UserInfo struct {
	Uid      uint64 // 用户id
	Username string // 用户名
	Role_Id  uint8  // 角色ID
	Cash     int32
	Diamond  int32
	State    uint8 // 玩家状态，0:空闲，1:游戏中，2:准备中 3:游戏中并掉线了
}

func NewUserInfo() *UserInfo {
	return &UserInfo{Uid: 0, Role_Id: 1, Cash: 1000, Diamond: 10, State: 0}
}

func (u *UserInfo) AddUser(username string, password string) (err error) {
	err = db.QueryRow("SELECT uid FROM users WHERE username = ?", username).Scan(&u.Uid)
	switch {
	case err == sql.ErrNoRows:
		result, err := db.Exec("INSERT INTO users(username, password) VALUES (?, ?)", username, password)
		if err != nil {
			log.Println("sql error", err)
			return err
		}
		insertid, err := result.LastInsertId()
		u.Uid = uint64(insertid)
		u.Username = username
		return err
	case err != nil:
		err = errors.New("sql error:" + username)
		return err
	default:
		err = errors.New("isExists:" + username)
		return
	}
}

func (u *UserInfo) GetUserInfo(uid uint64) (err error) {
	err = db.QueryRow("SELECT uid, username, role_id, cash, diamond FROM users WHERE uid = ?", uid).Scan(&u.Uid, &u.Username, &u.Role_Id, &u.Cash, &u.Diamond)

	if err != nil {
		log.Println("uid :", uid, " is no exits")
	}

	return
}

func (u *UserInfo) CheckLogin(username string, password string) (err error) {
	rows := db.QueryRow("SELECT uid, username, role_id, cash, diamond FROM users WHERE username = ? and password = ?", username, password)

	err = rows.Scan(&u.Uid, &u.Username, &u.Role_Id, &u.Cash, &u.Diamond)
	log.Println("err:", err, "userinfo:", u)
	if err != nil {
		log.Println("username or password is error")
	}
	return
}

func (u *UserInfo) ChangeRole(uid uint64, role_id uint8) (err error) {
	_, err = db.Query("UPDATE users SET role_id = ? WHERE uid = ? ", role_id, uid)
	if err != nil {
		log.Println(err)
	}
	u.Role_Id = role_id
	return
}

func (u *UserInfo) ChangeCash(uid uint64, cash int32) (err error) {
	_, err = db.Query("UPDATE users SET cash = cash + ? WHERE uid = ? ", cash, uid)
	if err != nil {
		log.Println(err)
	}

	u.Cash = u.Cash + cash
	return
}

func (u *UserInfo) ChangeDiamond(uid uint64, diamond int32) (err error) {
	_, err = db.Query("UPDATE users SET diamond = diamond + ? WHERE uid = ? ", diamond, uid)
	if err != nil {
		log.Println(err)
	}

	u.Diamond = u.Diamond + diamond
	return
}
