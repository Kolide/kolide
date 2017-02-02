// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import "github.com/kolide/kolide/server/kolide"

var _ kolide.LicenseStore = (*LicenseStore)(nil)

type SaveLicenseFunc func(tokenString string, publicKey string) (*kolide.License, error)

type LicenseFunc func() (*kolide.License, error)

type PublicKeyFunc func(tokenString string) (string, error)

type LicenseStore struct {
	SaveLicenseFunc        SaveLicenseFunc
	SaveLicenseFuncInvoked bool

	LicenseFunc        LicenseFunc
	LicenseFuncInvoked bool

	PublicKeyFunc        PublicKeyFunc
	PublicKeyFuncInvoked bool
}

func (s *LicenseStore) SaveLicense(tokenString string, publicKey string) (*kolide.License, error) {
	s.SaveLicenseFuncInvoked = true
	return s.SaveLicenseFunc(tokenString, publicKey)
}

func (s *LicenseStore) License() (*kolide.License, error) {
	s.LicenseFuncInvoked = true
	return s.LicenseFunc()
}

func (s *LicenseStore) LicensePublicKey(tokenString string) (string, error) {
	s.PublicKeyFuncInvoked = true
	return s.PublicKeyFunc(tokenString)
}
