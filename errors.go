package main

import (
	"github.com/cenkalti/backoff/v4"
	"strings"
)

func isNetworkError(err error) bool {
	networkErrors := []string{
		"i/o timeout",
		"connection refused",
		"EOF",
	}

	for _, e := range networkErrors {
		if strings.Contains(err.Error(), e) {
			return true
		}
	}

	return false
}

func maybeRetry(err error) error {
	if isNetworkError(err) {
		return err
	}

	return backoff.Permanent(err)
}
