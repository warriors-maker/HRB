package HRBAlgorithm

func ParseInt(i interface{}) (bool, int) {
	switch v := i.(type) {
	case int:
		return true,v
	}
	return false, 0
}

func ParseString(i interface{}) (bool, string) {
	switch v := i.(type) {
	case string:
		return true,v
	}
	return false, ""
}

