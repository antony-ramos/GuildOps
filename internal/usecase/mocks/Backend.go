// Code generated by mockery v2.33.3. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/antony-ramos/guildops/internal/entity"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Backend is an autogenerated mock type for the Backend type
type Backend struct {
	mock.Mock
}

// CreateAbsence provides a mock function with given fields: ctx, absence
func (_m *Backend) CreateAbsence(ctx context.Context, absence entity.Absence) (entity.Absence, error) {
	ret := _m.Called(ctx, absence)

	var r0 entity.Absence
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Absence) (entity.Absence, error)); ok {
		return rf(ctx, absence)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.Absence) entity.Absence); ok {
		r0 = rf(ctx, absence)
	} else {
		r0 = ret.Get(0).(entity.Absence)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.Absence) error); ok {
		r1 = rf(ctx, absence)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateLoot provides a mock function with given fields: ctx, loot
func (_m *Backend) CreateLoot(ctx context.Context, loot entity.Loot) (entity.Loot, error) {
	ret := _m.Called(ctx, loot)

	var r0 entity.Loot
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Loot) (entity.Loot, error)); ok {
		return rf(ctx, loot)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.Loot) entity.Loot); ok {
		r0 = rf(ctx, loot)
	} else {
		r0 = ret.Get(0).(entity.Loot)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.Loot) error); ok {
		r1 = rf(ctx, loot)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePlayer provides a mock function with given fields: ctx, player
func (_m *Backend) CreatePlayer(ctx context.Context, player entity.Player) (entity.Player, error) {
	ret := _m.Called(ctx, player)

	var r0 entity.Player
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Player) (entity.Player, error)); ok {
		return rf(ctx, player)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.Player) entity.Player); ok {
		r0 = rf(ctx, player)
	} else {
		r0 = ret.Get(0).(entity.Player)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.Player) error); ok {
		r1 = rf(ctx, player)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateRaid provides a mock function with given fields: ctx, raid
func (_m *Backend) CreateRaid(ctx context.Context, raid entity.Raid) (entity.Raid, error) {
	ret := _m.Called(ctx, raid)

	var r0 entity.Raid
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Raid) (entity.Raid, error)); ok {
		return rf(ctx, raid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.Raid) entity.Raid); ok {
		r0 = rf(ctx, raid)
	} else {
		r0 = ret.Get(0).(entity.Raid)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.Raid) error); ok {
		r1 = rf(ctx, raid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateStrike provides a mock function with given fields: ctx, strike, player
func (_m *Backend) CreateStrike(ctx context.Context, strike entity.Strike, player entity.Player) error {
	ret := _m.Called(ctx, strike, player)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Strike, entity.Player) error); ok {
		r0 = rf(ctx, strike, player)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAbsence provides a mock function with given fields: ctx, absenceID
func (_m *Backend) DeleteAbsence(ctx context.Context, absenceID int) error {
	ret := _m.Called(ctx, absenceID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, absenceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteLoot provides a mock function with given fields: ctx, lootID
func (_m *Backend) DeleteLoot(ctx context.Context, lootID int) error {
	ret := _m.Called(ctx, lootID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, lootID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePlayer provides a mock function with given fields: ctx, playerID
func (_m *Backend) DeletePlayer(ctx context.Context, playerID int) error {
	ret := _m.Called(ctx, playerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, playerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteRaid provides a mock function with given fields: ctx, raidID
func (_m *Backend) DeleteRaid(ctx context.Context, raidID int) error {
	ret := _m.Called(ctx, raidID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, raidID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteStrike provides a mock function with given fields: ctx, strikeID
func (_m *Backend) DeleteStrike(ctx context.Context, strikeID int) error {
	ret := _m.Called(ctx, strikeID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, strikeID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadAbsence provides a mock function with given fields: ctx, absenceID
func (_m *Backend) ReadAbsence(ctx context.Context, absenceID int) (entity.Absence, error) {
	ret := _m.Called(ctx, absenceID)

	var r0 entity.Absence
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (entity.Absence, error)); ok {
		return rf(ctx, absenceID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) entity.Absence); ok {
		r0 = rf(ctx, absenceID)
	} else {
		r0 = ret.Get(0).(entity.Absence)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, absenceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadLoot provides a mock function with given fields: ctx, lootID
func (_m *Backend) ReadLoot(ctx context.Context, lootID int) (entity.Loot, error) {
	ret := _m.Called(ctx, lootID)

	var r0 entity.Loot
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (entity.Loot, error)); ok {
		return rf(ctx, lootID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) entity.Loot); ok {
		r0 = rf(ctx, lootID)
	} else {
		r0 = ret.Get(0).(entity.Loot)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, lootID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadPlayer provides a mock function with given fields: ctx, playerID
func (_m *Backend) ReadPlayer(ctx context.Context, playerID int) (entity.Player, error) {
	ret := _m.Called(ctx, playerID)

	var r0 entity.Player
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (entity.Player, error)); ok {
		return rf(ctx, playerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) entity.Player); ok {
		r0 = rf(ctx, playerID)
	} else {
		r0 = ret.Get(0).(entity.Player)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, playerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadRaid provides a mock function with given fields: ctx, raidID
func (_m *Backend) ReadRaid(ctx context.Context, raidID int) (entity.Raid, error) {
	ret := _m.Called(ctx, raidID)

	var r0 entity.Raid
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (entity.Raid, error)); ok {
		return rf(ctx, raidID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) entity.Raid); ok {
		r0 = rf(ctx, raidID)
	} else {
		r0 = ret.Get(0).(entity.Raid)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, raidID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadStrike provides a mock function with given fields: ctx, strikeID
func (_m *Backend) ReadStrike(ctx context.Context, strikeID int) (entity.Strike, error) {
	ret := _m.Called(ctx, strikeID)

	var r0 entity.Strike
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (entity.Strike, error)); ok {
		return rf(ctx, strikeID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) entity.Strike); ok {
		r0 = rf(ctx, strikeID)
	} else {
		r0 = ret.Get(0).(entity.Strike)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, strikeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchAbsence provides a mock function with given fields: ctx, playerName, playerID, date
func (_m *Backend) SearchAbsence(ctx context.Context, playerName string, playerID int, date time.Time) ([]entity.Absence, error) {
	ret := _m.Called(ctx, playerName, playerID, date)

	var r0 []entity.Absence
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, time.Time) ([]entity.Absence, error)); ok {
		return rf(ctx, playerName, playerID, date)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, time.Time) []entity.Absence); ok {
		r0 = rf(ctx, playerName, playerID, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Absence)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, time.Time) error); ok {
		r1 = rf(ctx, playerName, playerID, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchLoot provides a mock function with given fields: ctx, name, date, difficulty
func (_m *Backend) SearchLoot(ctx context.Context, name string, date time.Time, difficulty string) ([]entity.Loot, error) {
	ret := _m.Called(ctx, name, date, difficulty)

	var r0 []entity.Loot
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, string) ([]entity.Loot, error)); ok {
		return rf(ctx, name, date, difficulty)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, string) []entity.Loot); ok {
		r0 = rf(ctx, name, date, difficulty)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Loot)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time, string) error); ok {
		r1 = rf(ctx, name, date, difficulty)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchPlayer provides a mock function with given fields: ctx, id, name
func (_m *Backend) SearchPlayer(ctx context.Context, id int, name string) ([]entity.Player, error) {
	ret := _m.Called(ctx, id, name)

	var r0 []entity.Player
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) ([]entity.Player, error)); ok {
		return rf(ctx, id, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, string) []entity.Player); ok {
		r0 = rf(ctx, id, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Player)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, id, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchRaid provides a mock function with given fields: ctx, raidName, date, difficulty
func (_m *Backend) SearchRaid(ctx context.Context, raidName string, date time.Time, difficulty string) ([]entity.Raid, error) {
	ret := _m.Called(ctx, raidName, date, difficulty)

	var r0 []entity.Raid
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, string) ([]entity.Raid, error)); ok {
		return rf(ctx, raidName, date, difficulty)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, string) []entity.Raid); ok {
		r0 = rf(ctx, raidName, date, difficulty)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Raid)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time, string) error); ok {
		r1 = rf(ctx, raidName, date, difficulty)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchStrike provides a mock function with given fields: ctx, playerID, Date, Season, Reason
func (_m *Backend) SearchStrike(ctx context.Context, playerID int, Date time.Time, Season string, Reason string) ([]entity.Strike, error) {
	ret := _m.Called(ctx, playerID, Date, Season, Reason)

	var r0 []entity.Strike
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, time.Time, string, string) ([]entity.Strike, error)); ok {
		return rf(ctx, playerID, Date, Season, Reason)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, time.Time, string, string) []entity.Strike); ok {
		r0 = rf(ctx, playerID, Date, Season, Reason)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Strike)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, time.Time, string, string) error); ok {
		r1 = rf(ctx, playerID, Date, Season, Reason)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAbsence provides a mock function with given fields: ctx, absence
func (_m *Backend) UpdateAbsence(ctx context.Context, absence entity.Absence) error {
	ret := _m.Called(ctx, absence)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Absence) error); ok {
		r0 = rf(ctx, absence)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLoot provides a mock function with given fields: ctx, loot
func (_m *Backend) UpdateLoot(ctx context.Context, loot entity.Loot) error {
	ret := _m.Called(ctx, loot)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Loot) error); ok {
		r0 = rf(ctx, loot)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePlayer provides a mock function with given fields: ctx, player
func (_m *Backend) UpdatePlayer(ctx context.Context, player entity.Player) error {
	ret := _m.Called(ctx, player)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Player) error); ok {
		r0 = rf(ctx, player)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateRaid provides a mock function with given fields: ctx, raid
func (_m *Backend) UpdateRaid(ctx context.Context, raid entity.Raid) error {
	ret := _m.Called(ctx, raid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Raid) error); ok {
		r0 = rf(ctx, raid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateStrike provides a mock function with given fields: ctx, strike
func (_m *Backend) UpdateStrike(ctx context.Context, strike entity.Strike) error {
	ret := _m.Called(ctx, strike)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Strike) error); ok {
		r0 = rf(ctx, strike)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewBackend creates a new instance of Backend. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBackend(t interface {
	mock.TestingT
	Cleanup(func())
}) *Backend {
	mock := &Backend{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}