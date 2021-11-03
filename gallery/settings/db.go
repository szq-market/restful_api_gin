package settings

import (
	"gin_project/gallery/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(cfg *Config) (*gorm.DB,error) {
	var(
		db *gorm.DB
		err error
	)
	db,err = gorm.Open(mysql.Open(cfg.DBConfig.DataSourceName),&gorm.Config{})

	if err != nil{
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.UserInfoModel{})
	return db,nil
}
