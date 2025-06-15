// Package user contains entities and operations related to users
package user

import "gorm.io/gorm"

// User represents a user entity in the system
// @Description User account information
type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(250)" json:"name"`
}
