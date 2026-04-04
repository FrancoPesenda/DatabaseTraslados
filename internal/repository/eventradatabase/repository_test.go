package eventradatabase

import "testing"

func TestRepository_New(t *testing.T) {
	_ = NewRepository(nil)
}
