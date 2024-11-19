package models

type SL struct {
	DL    // Repetetive Code and Title is handled like DL
	HasDL bool
}

// non-DB Validation just like DL validation:
func (sl *SL) Validate() error {
	return sl.DL.Validate()
}

// Returns a new valid SL, or err.
func NewSL(code string, title string, hasdl bool) (*SL, error) {
	newsl := SL{
		DL: DL{
			Code:    code,
			Title:   title,
			Version: 0,
		},
		HasDL: hasdl,
	}

	if err := newsl.Validate(); err != nil {
		return nil, err
	}

	return &newsl, nil
}
