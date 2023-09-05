package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TODO Should be more readable
func TestValidateDate(t *testing.T) {
	now := time.Now().Add(time.Hour)
	daysUntilNextWednesday := 3 - int(now.Weekday())
	if daysUntilNextWednesday < 0 {
		daysUntilNextWednesday += 7
	}
	nextWednesday := now.AddDate(0, 0, daysUntilNextWednesday)
	err := validateDate(nextWednesday)
	assert.NoError(t, err, "a date in the future on wednesday must work")

	now = time.Now().Add(time.Hour)
	daysUntilNextMonday := 1 - int(now.Weekday())
	if daysUntilNextMonday < 0 {
		daysUntilNextMonday += 7
	}
	nextMonday := now.AddDate(0, 0, daysUntilNextMonday)
	err = validateDate(nextMonday)
	assert.NoError(t, err, "a date in the future on monday must work")

	now = time.Now().Add(time.Hour)
	daysUntilNextThursday := 4 - int(now.Weekday())
	if daysUntilNextThursday < 0 {
		daysUntilNextThursday += 7
	}
	nextThursday := now.AddDate(0, 0, daysUntilNextThursday)
	err = validateDate(nextThursday)
	assert.NoError(t, err, "a date in the future on thursday must work")

	// Test avec une date invalide (hier)
	yesterday := now.AddDate(0, 0, -1)
	err = validateDate(yesterday)
	assert.EqualError(t, err, ErrorDateCannotBeBeforeToday.Error(), "yesterday should throw an error")

	// Test avec une date invalide (mardi)
	now = time.Now().Add(time.Hour)
	daysUntilNextTuesday := 2 - int(now.Weekday())
	if daysUntilNextTuesday < 0 {
		daysUntilNextTuesday += 7
	}
	nextTuesday := now.AddDate(0, 0, daysUntilNextTuesday)
	err = validateDate(nextTuesday)
	assert.EqualError(t, err, ErrorDateIsNotMonWesThu.Error(), "next tuesday should throw an error")
}

func TestAbsenceValidate(t *testing.T) {
	absence := Absence{Name: "John", Date: time.Now()}
	err := absence.Validate()
	assert.NoError(t, err, "L'absence devrait être valide")

	absence = Absence{Name: "", Date: time.Now()}
	err = absence.Validate()
	assert.EqualError(t, err, ErrorNameCannotBeEmpty.Error(), "L'absence sans nom devrait générer une erreur")
}
