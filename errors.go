package main

import (
	"fmt"
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

func collectErrorStrings(errors []error) []string {
	str := make([]string, 0)
	for _, e := range errors {
		str = append(str, e.Error())
	}
	return str
}

func multiError(errors ...error) error {
	collected := collectErrorStrings(errors)

	return fmt.Errorf("%d errors: %s", len(collected), strings.Join(collected, ", "))
}
