package vo

func StringPtr(v string) *string {
	return &v
}

func Float32Ptr(v float32) *float32 {
	return &v
}

func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
