package models

type SL struct {
	DL
	HasDL bool
}

func (sl *SL) Validate() error {
	return sl.DL.Validate()
}

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
