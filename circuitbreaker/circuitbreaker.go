package circuitbreaker

import (
	"context"
	"errors"
	"time"
)

type Circuit func(ctx context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold int) Circuit {
	var consecutiveFailure = 0
	var lastAttempt = time.Now()

	return func(ctx context.Context) (string, error) {
		d := consecutiveFailure - failureThreshold
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				return "", errors.New("service unreachable")
			}
		}

		lastAttempt = time.Now()
		response, err := circuit(ctx)
		if err != nil {
			consecutiveFailure++
			return response, err
		}
		consecutiveFailure = 0
		return response, nil
	}
}