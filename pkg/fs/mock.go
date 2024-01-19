package fs

import "github.com/stretchr/testify/mock"

type MockFS struct {
	mock.Mock
}

func (m *MockFS) Delete(path string) error {
	ret := m.Called(path)

	if ret.Get(0) == nil {
		return nil
	}

	return ret.Error(0)
}

func (m *MockFS) Copy(src, dst string) error {
	ret := m.Called(src, dst)

	if ret.Get(0) == nil {
		return nil
	}

	return ret.Error(0)
}

func (m *MockFS) Exists(path string) bool {
	ret := m.Called(path)

	return ret.Get(0).(bool) //nolint:forcetypeassert
}

func (m *MockFS) Equal(src, dst string) bool {
	ret := m.Called(src, dst)

	return ret.Get(0).(bool) //nolint:forcetypeassert
}

type MockShadowScanner struct {
	mock.Mock
}

func (m *MockShadowScanner) Start(rootPath, dstRootPath string, fileCallback, dirCallback ScannerCallback) chan error {
	ret := m.Called(rootPath, dstRootPath, fileCallback, dirCallback)

	return ret.Get(0).(chan error) //nolint:forcetypeassert
}

func (m *MockShadowScanner) Stop() {
	m.Called()
}
