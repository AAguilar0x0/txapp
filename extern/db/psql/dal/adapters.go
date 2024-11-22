package dal

import (
	"time"

	"github.com/AAguilar0x0/txapp/core/pkg/assert"
)

// User

func (d *User) GetID() string {
	return d.ID
}
func (d *User) GetPassword() string {
	return d.Password
}
func (d *User) GetRole() string {
	return d.Role
}

func (d *UserGetForAuthRow) GetID() string {
	return d.ID
}
func (d *UserGetForAuthRow) GetPassword() string {
	return d.Password
}
func (d *UserGetForAuthRow) GetRole() string {
	return d.Role
}

// Token

func (d *RefreshToken) GetID() string {
	return d.ID
}
func (d *RefreshToken) GetSub() string {
	return d.UserID
}
func (d *RefreshToken) GetIss() string {
	assert.Never("Not implemented", "fault", "dal.RefreshToken")
	return ""
}
func (d *RefreshToken) GetExp() time.Time {
	return d.ExpiresAt
}
