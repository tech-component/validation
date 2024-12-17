package interfaces

type Validator interface {
	ValidateStruct(s interface{}) error
}
