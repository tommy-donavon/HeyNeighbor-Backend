package data

type IUserCreate interface {
	CreateUser(*User) error
}

func (ur *UserRepo) CreateUser(user *User) error {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash
	return ur.db.Create(user).Error
}
