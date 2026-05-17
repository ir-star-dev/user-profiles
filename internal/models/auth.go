package models

import "time"

type RefreshToken struct {
	Id        int       `db:"id"`
	TokenHash []byte    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
	UserId    int       `db:"user_id"`
	FamilyID  string    `db:"family_id"`
}

type Logins struct {
	Id        int    `db:"id"`
	UserId    int    `db:"user_id"`
	Device    string `db:"device"`
	Ip        string `db:"ip"`
	UserAgent string `db:"user_agent"`
}

type LoginsResponse struct {
	Id        int       `db:"id"`
	UserId    int       `db:"user_id"`
	Device    string    `db:"device"`
	Ip        string    `db:"ip"`
	UserAgent string    `db:"user_agent"`
	Username  string    `db:"username"`
	LoginAt   time.Time `db:"login_at"`
}