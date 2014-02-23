package model

import (
	"database/sql"
	//"errors"
	"log"
)

type RoleInfo struct {
	Id         uint64
	Uid        uint64
	Role_Id    uint8
	Role_Attr  uint8
	Role_Level uint32
	Role_Name  string
}

func NewRoleInfo() *RoleInfo {
	return &RoleInfo{Role_Id: 1, Role_Name: "summer", Role_Level: 1, Role_Attr: 1}
}

//增加角色
func (r *RoleInfo) AddRole(uid uint64, role_id uint8, role_name string, role_level uint32, role_attr uint8) (err error) {
	result, err := db.Exec("INSERT INTO user_role (uid, role_id, role_name, role_level, role_attr) VALUES (?, ?, ?, ?, ?)", uid, role_id, role_name, role_level, role_attr)
	if err != nil {
		log.Println("sql error", err)
		return
	}
	insertid, err := result.LastInsertId()
	if err != nil {
		log.Println("增加角色失败")
	}
	r.Id = uint64(insertid)
	r.Uid = uid
	r.Role_Id = role_id
	r.Role_Name = role_name
	r.Role_Level = role_level
	r.Role_Attr = role_attr
	return
}

//获取角色信息
func (r *RoleInfo) GetRoleInfo(uid uint64, role_id uint8) (err error) {
	err = db.QueryRow("Select * FROM user_role WHERE uid = ? AND role_id = ?", uid, role_id).Scan(&r.Id, &r.Uid, &r.Role_Id, &r.Role_Name, &r.Role_Level, &r.Role_Attr)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return
	}
	return
}

//获取玩家所有的角色
func (r RoleInfo) GetRoleList(uid uint64) (urList []RoleInfo, err error) {
	rows, err := db.Query("Select * FROM user_role WHERE uid = ?", uid)

	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		ur := new(RoleInfo)
		err = rows.Scan(&ur.Id, &ur.Uid, &ur.Role_Id, &ur.Role_Name, &ur.Role_Level, &ur.Role_Attr)
		if err != nil {
			log.Println(err)
			return
		}
		urList = append(urList, *ur)
	}
	return
}

//角色升级
func (r *RoleInfo) RoleLevelUp(uid uint64, role_id uint8) (err error) {
	_, err = db.Query("UPDATE user_role SET role_level = role_level + 1 WHERE uid = ? AND role_id = ? ", uid, role_id)
	if err != nil {
		log.Println(err)
	}
	r.Role_Level = r.Role_Level + 1
	return
}
