package core

//User user
type User struct {
	Name         string `gorm:"not null;index"`
	Email        string `gorm:"not null;unique"`
	ProviderID   string `gorm:"not null;inde"`
	ProviderType string `gorm:"not null;index"`
	ContactID    uint
	Contact      Contact
}

//Contact contact
type Contact struct {
}

//Role role
type Role struct {
	Name string
}

//Permission permission
type Permission struct {
	RoleID uint
}

//Log log
type Log struct {
}
