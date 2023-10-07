package logformat

import (
	"fmt"
)

func Logformat(severity, requestID string, v ...interface{}) string {
	args := []interface{}{fmt.Sprintf("%-5s [%s] ", severity, requestID)}
	args = append(args, v...)
	return fmt.Sprint(args...)
}
