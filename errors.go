package main

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"strings"
)

var retryableErrors = []string{
	"i/o timeout",
	"connection refused",
	"EOF",
	"No cluster leader",
}

func isRetryable(err error) bool {
	for _, e := range retryableErrors {
		if strings.Contains(err.Error(), e) {
			return true
		}
	}

	return false
}

func maybeRetry(err error) error {
	if isRetryable(err) {
		return err
	}

	return backoff.Permanent(err)
}

func collectErrorStrings(errors []error) []string {
	str := make([]string, 0)
	for _, e := range errors {
		if e != nil {
			str = append(str, e.Error())
		}
	}
	return str
}

func multiError(errors ...error) error {
	collected := collectErrorStrings(errors)

	if len(collected) > 0 {
		return fmt.Errorf("%d errors: %s", len(collected), strings.Join(collected, ", "))
	}

	return nil
}
