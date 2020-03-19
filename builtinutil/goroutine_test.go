package builtinutil

import "testing"

func TestWrappedGo(t *testing.T) {
	// Can handle panic
	done := make(chan struct{}, 1)
	go WrappedGo(nil, func() {
		defer func() {
			done <- struct{}{}
		}()

		panic("testing")
	})

	<-done
}
