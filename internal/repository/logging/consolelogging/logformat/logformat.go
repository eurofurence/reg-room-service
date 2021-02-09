package logformat

import (
	"fmt"
)

func Logformat(severity string, requestId string, v ...interface{}) string {
	args := []interface{}{fmt.Sprintf("%-5s [%s] ", severity, requestId)}
	args = append(args, v...)
	return fmt.Sprint(args...)
}