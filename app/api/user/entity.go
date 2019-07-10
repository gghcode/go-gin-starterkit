package user

var emptyUser = User{}

// User is user data model
type User struct {
	ID           int64  `gorm:"primary_key;"`
	UserName     string `gorm:"unique;not null;"`
	PasswordHash []byte `gorm:"not null;"`
	CreatedAt    int64  `gorm:"not null;"`
}
