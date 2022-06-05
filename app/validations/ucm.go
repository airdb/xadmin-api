package validations

type UcmServiceValidations interface {
}

type ucmServiceValidations struct {
}

func CreateUcmServiceValidations() UcmServiceValidations {
	return new(ucmServiceValidations)
}
