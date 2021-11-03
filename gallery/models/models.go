package models

import "gorm.io/gorm"

type (
	//定义原始的数据库字段
	UserInfoModel struct {
		gorm.Model
		Name string `json:"name"`
		Sex string `json:"sex"`
		Phone int `json:"phone"`
		City string `json:"city"`
	}
	//处理返回的字段
	TransformedUserInfo struct {
		ID uint `json:"id"`
		Name string `json:"name"`
		Sex string `json:"sex"`
		City string `json:"city"`
		Phone int `json:"phone"`
	}

)
//封装单独的这个model
type UserDB struct {
	Db *gorm.DB
}

//实例化
func NewUsersDB(db *gorm.DB) *UserDB {
	return &UserDB{
		Db: db,
	}
}