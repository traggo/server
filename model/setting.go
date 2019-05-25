package model

// Setting a setting for a user.
type Setting struct {
	UserID    int
	Namespace string
	Key       string
	Value     string
}
