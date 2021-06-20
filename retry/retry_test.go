package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var count int
func EmulateTransientFailure(ctx context.Context) (string, error) {
	count++
	if count <= 3 {
		return "intentional fail", errors.New("error")
	} else {
		return "success", nil
	}
}

func TestRetry(t *testing.T) {
	ctx := context.Background()
	r := Retry(EmulateTransientFailure, 5, 2 * time.Second)
	_, err := r(ctx)
	if err != nil {
		t.Error("failed after 5 retries")
	}
}
