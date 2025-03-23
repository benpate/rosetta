package funcmap

func addLogicFuncs(target map[string]any) {

	target["iif"] = func(condition bool, trueValue any, falseValue any) any {
		if condition {
			return trueValue
		}
		return falseValue
	}

	target["and"] = func(values ...bool) bool {
		for _, value := range values {
			if !value {
				return false
			}
		}
		return true
	}

	target["or"] = func(values ...bool) bool {
		for _, value := range values {
			if value {
				return true
			}
		}
		return false
	}

}
