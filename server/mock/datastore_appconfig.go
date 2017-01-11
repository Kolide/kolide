// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import "github.com/kolide/kolide-ose/server/kolide"

var _ kolide.AppConfigStore = (*AppConfigStore)(nil)

type NewAppConfigFunc func(info *kolide.AppConfig) (*kolide.AppConfig, error)

type AppConfigFunc func() (*kolide.AppConfig, error)

type SaveAppConfigFunc func(info *kolide.AppConfig) error

type AppConfigStore struct {
	NewAppConfigFunc        NewAppConfigFunc
	NewAppConfigFuncInvoked bool

	AppConfigFunc        AppConfigFunc
	AppConfigFuncInvoked bool

	SaveAppConfigFunc        SaveAppConfigFunc
	SaveAppConfigFuncInvoked bool
}

func (s *AppConfigStore) NewAppConfig(info *kolide.AppConfig) (*kolide.AppConfig, error) {
	s.NewAppConfigFuncInvoked = true
	return s.NewAppConfigFunc(info)
}

func (s *AppConfigStore) AppConfig() (*kolide.AppConfig, error) {
	s.AppConfigFuncInvoked = true
	return s.AppConfigFunc()
}

func (s *AppConfigStore) SaveAppConfig(info *kolide.AppConfig) error {
	s.SaveAppConfigFuncInvoked = true
	return s.SaveAppConfigFunc(info)
}
