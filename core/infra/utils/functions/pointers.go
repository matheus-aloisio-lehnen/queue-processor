package functions

func PtrStr(s string) *string {
	return &s
}

func PtrBool(b bool) *bool {
	return &b
}
