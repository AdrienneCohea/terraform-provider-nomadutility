package main

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_maybeRetry_withPermanentError(t *testing.T) {
	err := fmt.Errorf("some permanent error")
	expected := backoff.Permanent(err)

	assert.IsType(t, expected, maybeRetry(err))
}

func Test_maybeRetry_withRetryableError(t *testing.T) {
	err := fmt.Errorf("i/o timeout")
	assert.IsType(t, err, maybeRetry(err))
}

func Test_isRetryable_positiveCases(t *testing.T) {
	for _, e := range retryableErrors {
		assert.True(t, isRetryable(fmt.Errorf(e)))
	}
}

func Test_isRetryable_negativeCases(t *testing.T) {
	negativeCases := []error{
		fmt.Errorf("ACLs disabled"),
		fmt.Errorf("ACLs already bootstrapped"),
	}
	for _, e := range negativeCases {
		assert.False(t, isRetryable(e))
	}
}
