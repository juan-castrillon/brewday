package tools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefractometer(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name     string
		OG       float32
		MG       float32
		WCF      float32
		Expected float32
	}{
		{Name: "WCF 1 - High gravity", OG: 1.070, MG: 1.010, WCF: 1, Expected: 0.994},
		{Name: "WCF 1 - Low gravity", OG: 1.040, MG: 1.010, WCF: 1, Expected: 1.000},
		{Name: "WCF 1.04 - High gravity", OG: 1.070, MG: 1.010, WCF: 1.04, Expected: 0.995},
		{Name: "WCF 1.04 - Low gravity", OG: 1.040, MG: 1.010, WCF: 1, Expected: 1.000},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			require.InDelta(tc.Expected, CorrectRefractometerAlcohol(tc.OG, tc.MG, tc.WCF), 0.001)
		})
	}
}
