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

func Test_multiError_withTwoErrors(t *testing.T) {
	expected := fmt.Errorf("2 errors: first, second")

	actual := multiError(
		fmt.Errorf("first"),
		fmt.Errorf("second"))

	assert.EqualValues(t, expected, actual)
}

func Test_multiError_withOneError(t *testing.T) {
	expected := fmt.Errorf("1 errors: second")

	actual := multiError(nil, fmt.Errorf("second"))

	assert.EqualValues(t, expected, actual)
}

func Test_multiError_withNoErrors(t *testing.T) {
	assert.NoError(t, multiError(nil, nil))
}

func Test_multiError_withNoErrorsOneArg(t *testing.T) {
	assert.NoError(t, multiError(nil))
}

func Test_multiError_withNoArgs(t *testing.T) {
	assert.NoError(t, multiError())
}
