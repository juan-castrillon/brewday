package tools

import (
	"fmt"
)

func AnyToString(value any) string {
	switch v := value.(type) {
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%.3f", v)
	default:
		return ""
	}
}
