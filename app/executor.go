package main


var dict = map[string]string{}

func executeCommandSet(k, v string) error {
	dict[k] = v
	return nil
}

func executeCommandGet(k string) (string, error) {
	v, ok := dict[k]
	if !ok {
		return "-1", nil
	}
	return v, nil
}
