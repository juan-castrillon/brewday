package tools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnyToString(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name      string
		Input     any
		InputType string
		Expected  string
	}{
		{Name: "int", Input: 1, InputType: "int", Expected: "1"},
		{Name: "int64", Input: int64(2), InputType: "int64", Expected: "2"},
		{Name: "int32", Input: int32(3), InputType: "int32", Expected: "3"},
		{Name: "float32", Input: float32(2.5), InputType: "float32", Expected: "2.500"},
		{Name: "float64", Input: float64(2.5), InputType: "float64", Expected: "2.500"},
		{Name: "float32_trim", Input: float32(2.5789), InputType: "float32", Expected: "2.579"},
		{Name: "float64_trim", Input: float64(2.5789), InputType: "float64", Expected: "2.579"},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var result string
			switch tc.InputType {
			case "int":
				v, ok := tc.Input.(int)
				require.True(ok)
				result = AnyToString(v)
			case "int32":
				v, ok := tc.Input.(int32)
				require.True(ok)
				result = AnyToString(v)
			case "int64":
				v, ok := tc.Input.(int64)
				require.True(ok)
				result = AnyToString(v)
			case "float32":
				v, ok := tc.Input.(float32)
				require.True(ok)
				result = AnyToString(v)
			case "float64":
				v, ok := tc.Input.(float64)
				require.True(ok)
				result = AnyToString(v)
			default:
				t.Fail()
			}
			require.Equal(tc.Expected, result)
		})
	}
}
