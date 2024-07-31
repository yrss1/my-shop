package helpers

func GetStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
