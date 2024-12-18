package models

import (
	"fmt"
)

const STRING_LIMIT = 64

// if string has more that STRING_LIMIT character return err
func ValidateSize(inp, variable string) error {
	if len([]rune(inp)) > STRING_LIMIT {
		return fmt.Errorf("%s Must not have more than %v characters", variable, STRING_LIMIT)
	}
	return nil
}

// if string is empty return err
func ValidateEmpty(inp, variable string) error {
	if inp == "" {
		return fmt.Errorf("%s cannot be nil", variable)
	}
	return nil
}

// check both emptiness and characters of a string
func ValidateString(inp string, variable string) error {
	if err := ValidateSize(inp, variable); err != nil {
		return err
	}
	if err := ValidateEmpty(inp, variable); err != nil {
		return err
	}
	return nil
}

type DL struct {
	ID      uint   `gorm:"primaryKey"`
	Code    string `gorm:"unique"` // Repetetive Code is handled in DB
	Title   string `gorm:"unique"` // Repetetive Title is handled in DB
	Version uint
}

// non-DB Validation:
func (dl *DL) Validate() error {

	if dl == nil {
		return fmt.Errorf("dl can't be nill")
	}

	if err := ValidateString(dl.Code, "Code"); err != nil {
		return err
	}

	if err := ValidateString(dl.Title, "Title"); err != nil {
		return err
	}

	return nil
}

// return a new valid DL, or err.
func NewDL(code string, title string) (*DL, error) {
	newdl := DL{
		Code:    code,
		Title:   title,
		Version: 0,
	}

	if err := newdl.Validate(); err != nil {
		return nil, err
	}

	return &newdl, nil
}
