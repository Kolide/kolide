package datastore

import "github.com/kolide/kolide-ose/kolide"

func (orm gormDB) NewPasswordResetRequest(req *kolide.PasswordResetRequest) (*kolide.PasswordResetRequest, error) {
	err := orm.DB.Create(req).Error
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (orm gormDB) SavePasswordResetRequest(req *kolide.PasswordResetRequest) error {
	return orm.DB.Save(req).Error
}

func (orm gormDB) DeletePasswordResetRequest(req *kolide.PasswordResetRequest) error {
	err := orm.DB.Delete(req).Error
	return err
}

func (orm gormDB) FindPassswordResetByID(id uint) (*kolide.PasswordResetRequest, error) {
	reset := &kolide.PasswordResetRequest{
		ID: id,
	}
	err := orm.DB.Find(reset).First(reset).Error
	return reset, err
}

func (orm gormDB) FindPassswordResetsByUserID(id uint) (*kolide.PasswordResetRequest, error) {
	reset := &kolide.PasswordResetRequest{
		UserID: id,
	}
	err := orm.DB.Find(reset).First(reset).Error
	return reset, err
}

func (orm gormDB) FindPassswordResetByToken(token string) (*kolide.PasswordResetRequest, error) {
	reset := &kolide.PasswordResetRequest{
		Token: token,
	}
	err := orm.DB.Find(reset).First(reset).Error
	return reset, err
}

func (orm gormDB) FindPassswordResetByTokenAndUserID(token string, userID uint) (*kolide.PasswordResetRequest, error) {
	reset := &kolide.PasswordResetRequest{
		Token:  token,
		UserID: userID,
	}
	err := orm.DB.Find(reset).First(reset).Error
	return reset, err
}
