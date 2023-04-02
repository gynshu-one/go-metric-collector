package tools

func Contains(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}
func Int64Ptr(i int64) *int64 {
	return &i
}
func Float64Ptr(f float64) *float64 {
	return &f
}
