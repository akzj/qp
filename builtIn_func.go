package qp

var (
	builtInFunctions = map[string]Function{
		"println": &println{},
	}
	arrayBuiltInFunctions = map[string]*Object{
		"append": &Object{
			inner: &appendArray{},
			label: "append",
		},
		"get":&Object{
			inner:   &getArray{},
			label:   "get",
		},
	}
	stringBuiltInFunctions = map[string]*Object{
		"to_lower": &Object{
			inner: &stringLowCase{},
			label: "to_lower",
		},
		"clone": &Object{
			inner: &StringObjectClone{},
			label: "clone",
		},
	}
)
