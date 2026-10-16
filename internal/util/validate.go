package util

import (
	"fmt"
	"regexp"
)

/*
validates string is between 2-25 characters, alphanumerica/hyphen/underscore only
and doesn't start or end with hyphen/underscore. Label is used to name field in errors.
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
	isValid = pattern.MatchString(str)
	if !isValid {
		return fmt.Errorf("%s cannot start or end with hyphens or underscores.", label)
	}
	return nil
}
