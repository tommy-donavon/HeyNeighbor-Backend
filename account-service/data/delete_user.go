package data

type IUserDelete interface {
	DeleteUser(string) error
}

func (ur *UserRepo) DeleteUser(username string) error {
	user, err := ur.GetUser(username)
	if err != nil {
		return err
	}
	return ur.db.Delete(&user).Error
}
