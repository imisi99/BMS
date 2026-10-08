// Package core
package core

import (
	"fmt"
)

type CustomError struct {
	Code    int
	Message string
}

func (cm *CustomError) Error() string {
	return fmt.Sprintf("[CM] An error occured -> %s", cm.Message)
}
