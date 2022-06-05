package validations

type UamServiceValidations interface {
}

type uamServiceValidations struct {
}

func CreateUamServiceValidations() UamServiceValidations {
	return new(uamServiceValidations)
}
