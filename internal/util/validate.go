package util

import (
	"fmt"
	"regexp"
)

/*
ValidateServerName returns nil when srvName is 2-25 characters, contains only
alphanumeric/hyphens/underscores, and doesn't end in hyphens/underscores.
*/
func ValidateServerName(srvName string) error {
	if len(srvName) < 2 {
		return fmt.Errorf("server name must be at least 2 characters")
	}
	if len(srvName) > 25 {
		return fmt.Errorf("server name cannot be longer than 25 characters")
	}

	if !regexp.MustCompile(`^[\w-]+$`).MatchString(srvName) {
		return fmt.Errorf("server name must only contain aphanumeric, hyphens, and underscores")
	}

	if !regexp.MustCompile(`^[A-Za-z0-9](.*[A-Za-z0-9])?$`).MatchString(srvName) {
		return fmt.Errorf("server name cannot start or end with hyphens or underscores")
	}
	return nil
}

/*
ValidateJVMMemory returns nil if the string is a valid JVM memory
argument (1G, 500MB, etc.)
*/
func ValidateJVMMemory(str string) error {
	jvmMemRegex := regexp.MustCompile(`^(?i)\d+(?:\.\d+)?\s*[kmg]b?$`)
	if !jvmMemRegex.MatchString(str) {
		return fmt.Errorf("invalid memory amount")
	}
	return nil
}
