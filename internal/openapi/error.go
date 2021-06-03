package openapi

type InvalidDefinitionError struct {
	message string
}

func (i *InvalidDefinitionError) Error() string {
	return i.message
}
