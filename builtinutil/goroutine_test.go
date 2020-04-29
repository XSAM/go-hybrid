package builtinutil

import (
	"context"
	"testing"
)

func TestWrappedGo(t *testing.T) {
	// Can handle panic
	done := make(chan struct{}, 1)
	go WrappedGo(context.Background(), func() {
		defer func() {
			done <- struct{}{}
		}()

		panic("testing")
	})

	<-done
}
