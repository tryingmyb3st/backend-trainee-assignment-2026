package auth_repository

import "time"

type UserModel struct {
	ID        string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
}
