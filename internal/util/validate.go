package util

import (
	"fmt"
	"graft/internal/sys"
	"regexp"
)

/*
ValidateServerName returns `nil` if string is between 2-25 characters, alphanumerica/hyphen/underscore only,
doesn't start or end with hyphen/underscore, and isn't already in use, otherwie it returns an error.
Primarily used to validate server name for its folder/service.
*/
func ValidateServerName(srvName string) error {
	if len(srvName) < 2 {
		return fmt.Errorf("Server name must be at least 2 characters.")
	}
	if len(srvName) > 25 {
		return fmt.Errorf("Server name cannot be longer than 25 characters.")
	}

	pattern := regexp.MustCompile(`^[\w-]+$`)
	isValid := pattern.MatchString(srvName)
	if !isValid {
		return fmt.Errorf("Server name must only contain aphanumeric, hyphens, and underscores.")
	}

	pattern = regexp.MustCompile(`^[A-Za-z0-9](.*[A-Za-z0-9])?$`)
	if !pattern.MatchString(srvName) {
		return fmt.Errorf("Server name cannot start or end with hyphens or underscores.")
	}

	srvExists, err := sys.DoesServerPathExist(srvName)
	if err != nil {
		return fmt.Errorf("Unknown error occured checking path: %s", err)
	}
	if srvExists {
		return fmt.Errorf("Server name already in use.")
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
