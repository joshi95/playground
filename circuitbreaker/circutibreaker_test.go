package circuitbreaker

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func failAfter(threshold int) Circuit {
	count := 0
	return func(ctx context.Context) (string, error) {
		count += 1
		if count > threshold {
			return "", errors.New("intentional fail")
		}
		return "Success", nil
	}
}

func TestCircuitBreaker(t *testing.T) {
	circuit := failAfter(3)
	breaker := Breaker(circuit, 1)

	ctx := context.Background()

	circuitOpen := false
	doesCircuitOpen := false
	doesCircuitReClose := false
	count := 0

	for range time.NewTicker(time.Second).C {
		_, err := breaker(ctx)
		if err != nil {
			if strings.HasPrefix(err.Error(), "service unreachable") {
				if !circuitOpen {
					circuitOpen = true
					doesCircuitOpen = true
					t.Log("circuit became open")
				}
			} else {
				if circuitOpen {
					circuitOpen = false
					doesCircuitReClose = true
					t.Log("circuit has automatically closed")
				}
			}
		} else {
			t.Log("circuit close and operational")
		}
		count += 1
		if count > 10 {
			break
		}
	}

	if !doesCircuitOpen {
		t.Error("circuit didn't open")
	}
	if !doesCircuitReClose {
		t.Error("circuit didn't re-close")
	}

}
