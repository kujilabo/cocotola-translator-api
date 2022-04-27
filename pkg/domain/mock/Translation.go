// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/kujilabo/cocotola-translator-api/pkg/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"

	time "time"
)

// Translation is an autogenerated mock type for the Translation type
type Translation struct {
	mock.Mock
}

// GetCreatedAt provides a mock function with given fields:
func (_m *Translation) GetCreatedAt() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// GetLang2 provides a mock function with given fields:
func (_m *Translation) GetLang2() domain.Lang2 {
	ret := _m.Called()

	var r0 domain.Lang2
	if rf, ok := ret.Get(0).(func() domain.Lang2); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Lang2)
		}
	}

	return r0
}

// GetPos provides a mock function with given fields:
func (_m *Translation) GetPos() domain.WordPos {
	ret := _m.Called()

	var r0 domain.WordPos
	if rf, ok := ret.Get(0).(func() domain.WordPos); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.WordPos)
	}

	return r0
}

// GetProvider provides a mock function with given fields:
func (_m *Translation) GetProvider() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetText provides a mock function with given fields:
func (_m *Translation) GetText() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetTranslated provides a mock function with given fields:
func (_m *Translation) GetTranslated() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetUpdatedAt provides a mock function with given fields:
func (_m *Translation) GetUpdatedAt() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// GetVersion provides a mock function with given fields:
func (_m *Translation) GetVersion() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// NewTranslation creates a new instance of Translation. It also registers a cleanup function to assert the mocks expectations.
func NewTranslation(t testing.TB) *Translation {
	mock := &Translation{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
