package models

type SL struct {
	DL
	HasSL bool
}

func (sl *SL) ValidateSL() error {

	if err := sl.DL.ValidateDL(); err != nil {
		return err
	}

	// Handle UnUpdatability Here:

	return nil
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

	if err := newsl.ValidateSL(); err != nil {
		return nil, err
	}

	return &newsl, nil
}
