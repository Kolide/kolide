// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import "github.com/kolide/kolide/server/kolide"

var _ kolide.OptionStore = (*OptionStore)(nil)

type SaveOptionsFunc func(opts []kolide.Option) error

type ListOptionsFunc func() ([]kolide.Option, error)

type OptionFunc func(id uint) (*kolide.Option, error)

type OptionByNameFunc func(name string) (*kolide.Option, error)

type GetOsqueryConfigOptionsFunc func() (map[string]interface{}, error)

type OptionStore struct {
	SaveOptionsFunc        SaveOptionsFunc
	SaveOptionsFuncInvoked bool

	ListOptionsFunc        ListOptionsFunc
	ListOptionsFuncInvoked bool

	OptionFunc        OptionFunc
	OptionFuncInvoked bool

	OptionByNameFunc        OptionByNameFunc
	OptionByNameFuncInvoked bool

	GetOsqueryConfigOptionsFunc        GetOsqueryConfigOptionsFunc
	GetOsqueryConfigOptionsFuncInvoked bool
}

func (s *OptionStore) SaveOptions(opts []kolide.Option) error {
	s.SaveOptionsFuncInvoked = true
	return s.SaveOptionsFunc(opts)
}

func (s *OptionStore) ListOptions() ([]kolide.Option, error) {
	s.ListOptionsFuncInvoked = true
	return s.ListOptionsFunc()
}

func (s *OptionStore) Option(id uint) (*kolide.Option, error) {
	s.OptionFuncInvoked = true
	return s.OptionFunc(id)
}

func (s *OptionStore) OptionByName(name string) (*kolide.Option, error) {
	s.OptionByNameFuncInvoked = true
	return s.OptionByNameFunc(name)
}

func (s *OptionStore) GetOsqueryConfigOptions() (map[string]interface{}, error) {
	s.GetOsqueryConfigOptionsFuncInvoked = true
	return s.GetOsqueryConfigOptionsFunc()
}
