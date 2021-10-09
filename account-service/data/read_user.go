package data

type IUserRead interface {
	GetUser(string) (*User, error)
}

func (ur *UserRepo) GetUser(username string) (*User, error) {
	user := User{}
	err := ur.db.Where("username = ?", username).First(&user).Error
	return &user, err
}
