// Code generated by mockery v2.33.3. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/antony-ramos/guildops/internal/entity"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// AbsenceUseCase is an autogenerated mock type for the AbsenceUseCase type
type AbsenceUseCase struct {
	mock.Mock
}

// CreateAbsence provides a mock function with given fields: ctx, playerName, date
func (_m *AbsenceUseCase) CreateAbsence(ctx context.Context, playerName string, date time.Time) error {
	ret := _m.Called(ctx, playerName, date)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) error); ok {
		r0 = rf(ctx, playerName, date)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAbsence provides a mock function with given fields: ctx, playerName, date
func (_m *AbsenceUseCase) DeleteAbsence(ctx context.Context, playerName string, date time.Time) error {
	ret := _m.Called(ctx, playerName, date)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) error); ok {
		r0 = rf(ctx, playerName, date)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListAbsence provides a mock function with given fields: ctx, date
func (_m *AbsenceUseCase) ListAbsence(ctx context.Context, date time.Time) ([]entity.Absence, error) {
	ret := _m.Called(ctx, date)

	var r0 []entity.Absence
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time) ([]entity.Absence, error)); ok {
		return rf(ctx, date)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Time) []entity.Absence); ok {
		r0 = rf(ctx, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Absence)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Time) error); ok {
		r1 = rf(ctx, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAbsenceUseCase creates a new instance of AbsenceUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAbsenceUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *AbsenceUseCase {
	mock := &AbsenceUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
