package sql

import (
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey"`
	Password string
}

func MigrateUser(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	return nil
}

func GetUserPassword(id string) (string, error) {
	var user User
	result := DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", result.Error
	}
	return user.Password, nil
}

func CreateUser(user *User) error {
	result := DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
