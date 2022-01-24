package models

import (
	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Uid           int64     `gorm:"comment:用户唯一ID;autoIncrement;primaryKey;not null"`
	LastLoginTime time.Time `gorm:"comment:最近一次登录的时间;not null"`
	Name          string    `gorm:"type:varchar(20);comment:用户名;not null; default:''"`
	Gender        string    `gorm:"type:varchar(6);comment:性别[male|famale];not null;default:male"`
	Email         string    `gorm:"type:varchar(60);comment:邮件;not null;default:''"`

	Balance string `gorm:"type:decimal(18,2);comment:账户余额;not null;default:0.00"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (l User) TableName() string {
	return "simple_users"
}

func Mock() *User {
	return &User{
		LastLoginTime: time.Now(),
		Name:          gofakeit.Name(),
		Gender:        gofakeit.Gender(),
		Email:         gofakeit.Email(),
	}
}
