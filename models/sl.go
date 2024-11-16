package models

type SL struct {
	DL
	HasSL bool
}

func (sl *SL) Validate() error {
	return sl.DL.Validate()
}

func NewSL(code string, title string, hassl bool) (*SL, error) {
	newsl := SL{
		DL: DL{
			Code:    code,
			Title:   title,
			Version: 0,
		},
		HasSL: hassl,
	}

	if err := newsl.Validate(); err != nil {
		return nil, err
	}

	return &newsl, nil
}
