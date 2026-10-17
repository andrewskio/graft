package util

import (
	"fmt"
	"regexp"
)

/*
ValidateSlug returns `nil` if string is between 2-25 characters, alphanumerica/hyphen/underscore only,
and doesn't start or end with hyphen/underscore, otherwie it returns an error. Label is used to name the offending field in errors.
Primarily used to validate server name for its folder/service.
*/
func ValidateSlug(str string, label string) error {
	if len(str) < 2 {
		return fmt.Errorf("%s must be at least 2 characters.", label)
	}
	if len(str) > 25 {
		return fmt.Errorf("%s cannot be longer than 25 characters.", label)
	}

	// only contains alphanumeric, hyphen, underscore
	pattern := regexp.MustCompile(`^[\w-]+$`)
	isValid := pattern.MatchString(str)
	if !isValid {
		return fmt.Errorf("%s must only contain aphanumeric, hyphens, and underscores.", label)
	}

	pattern = regexp.MustCompile(`^[A-Za-z0-9](.*[A-Za-z0-9])?$`)
	if !pattern.MatchString(str) {
		return fmt.Errorf("%s cannot start or end with hyphens or underscores.", label)
	}
	return nil
}

/*
ValidateJVMMemory takes a string and returns `nil` if the string is a valid JVM memory
argument (1G, 500MB, etc.), otherwise returns an error.
*/
func ValidateJVMMemory(str string) error {
	jvmMemRegex := regexp.MustCompile(`^(?i)\d+(?:\.\d+)?\s*[kmg]b?$`)
	if !jvmMemRegex.MatchString(str) {
		return fmt.Errorf("Invalid memory amount.")
	}
	return nil
}
