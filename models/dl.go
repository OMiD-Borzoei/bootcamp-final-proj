package models

import "fmt"

const STRING_LIMIT = 64

func ValidateString(inp string, variable string) error {
	if len([]rune(inp)) > STRING_LIMIT {
		return fmt.Errorf("%s Must not have more than %v characters", variable, STRING_LIMIT)
	}

	if inp == "" {
		return fmt.Errorf("%s cannot be nil", variable)
	}

	return nil
}

type DL struct {
	Code    string `gorm:"primaryKey"`
	Title   string `gorm:"unique"`
	Version uint
}

func (dl *DL) ValidateDL() error {

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

func NewDL(code string, title string) (*DL, error) {
	newdl := DL{
		Code:    code,
		Title:   title,
		Version: 0,
	}

	if err := newdl.ValidateDL(); err != nil {
		return nil, err
	}

	return &newdl, nil
}
