package auth

import "time"

type RefreshToken struct {
	Id        int       `db:"id"`
	TokenHash []byte    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
	UserId    int       `db:"user_id"`
	FamilyID  string	`db:"family_id"`
}
