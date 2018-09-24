package config

import "github.com/gin-gonic/gin"

type user struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
}

type users []user

// AdminConfig - configs for admin
type AdminConfig struct {
	Admins users `yaml:"admins"`
}

// Accounts - convert to gin accounts
func (ac AdminConfig) Accounts() gin.Accounts {
	accounts := make(gin.Accounts)
	for _, user := range ac.Admins {
		accounts[user.Name] = user.Password
	}
	return accounts
}
