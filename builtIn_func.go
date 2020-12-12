package qp

var (
	builtInFunctions = map[string]*Object{
		"println": &Object{inner: printlnFunc{}},
		"now":     &Object{inner: NowFunc{}},
	}
	arrayBuiltInFunctions = map[string]*Object{
		"append": &Object{
			inner: &appendArray{},
			label: "append",
		},
		"get": &Object{
			inner: &getArray{},
			label: "get",
		},
		"size": {
			inner: getArraySize{},
			label: "size",
		},
	}
	stringBuiltInFunctions = map[string]*Object{
		"to_lower": &Object{
			inner: stringLowCase{},
			label: "to_lower",
		},
		"clone": &Object{
			inner: StringObjectClone{},
			label: "clone",
		},
	}
)
