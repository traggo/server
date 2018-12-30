package test_test

type fakeTesting struct {
	hasErrors bool
}

func (t *fakeTesting) Errorf(format string, args ...interface{}) {
	t.hasErrors = true
}
