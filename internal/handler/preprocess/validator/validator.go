package validator

type Validator interface {
	ValidateLine(line string) error
}
