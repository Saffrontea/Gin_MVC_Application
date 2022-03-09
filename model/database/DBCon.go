/*
	database
	データベース処理
*/
package database

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//DB データベース本体
var DB *gorm.DB

/*
	DBConnection
	データベースに接続を開始
*/
func DBConnection() error {
	if DB == nil {
		var err error
		DB, err = getNewDB()

		if err != nil {
			return errors.New("Cannot connect database")
		}
	}
	return nil
}

//getNewDB 新規コネクションを獲得
func getNewDB() (*gorm.DB, error) {
	var db *gorm.DB
	db, err := gorm.Open(sqlite.Open("data.DB"), &gorm.Config{})
	//db, err := gorm.Open(mysql.Open("20jy0139:20jy0139@tcp(10.64.144.5:3306)/20jy0139?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return db, nil
}

//Migrator DBと構造体をリンクする
func Migrator(i []interface{}) error {
	for _, i3 := range i {
		err := DB.AutoMigrate(i3)
		if err != nil {
			return err
		}
	}
	return nil
}

//Transaction トランザクション処理の関数を入れる
func Transaction(f func(tx *gorm.DB) error) error {
	//db, _ := getNewDB()

	return DB.Transaction(f)
}
