package transactor

func setDefaultIfNotDefined(prop interface{}, def interface{}) interface{} {
	if prop == nil {
		return def
	}
	return prop
}