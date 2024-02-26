package cmd

func rightpad(s string, length int) string {
	for len(s) < length {
		s = s + " "
	}
	return s
}

func leftpad(s string, length int) string {
	for len(s) < length {
		s = " " + s
	}
	return s
}
