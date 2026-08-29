package systeminfo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// helper to build a standard fixture of systems used across tests
func testSystems() map[string]SystemProperties {
	return map[string]SystemProperties{
		"sysA": {
			LD:     10.5,
			UD:     25.5,
			Power:  100,
			MaxVol: 500.0,
		},
		"sysB": {
			LD:     0,
			UD:     0,
			Power:  0,
			MaxVol: 0,
		},
		"sysC": {
			LD:     5.25,
			UD:     99.99,
			Power:  10,
			MaxVol: 1234.5,
		},
	}
}

func TestNewInfoProvider(t *testing.T) {
	tests := []struct {
		name    string
		systems map[string]SystemProperties
	}{
		{name: "nil map", systems: nil},
		{name: "empty map", systems: map[string]SystemProperties{}},
		{name: "populated map", systems: testSystems()},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip, err := NewInfoProvider(tc.systems)

			require.NoError(t, err)
			require.NotNil(t, ip)
			require.Equal(t, tc.systems, ip.systems)
		})
	}
}

func TestHasSystems(t *testing.T) {
	require := require.New(t)
	tests := []struct {
		name     string
		systems  map[string]SystemProperties
		expected bool
	}{
		{
			name:     "Empty",
			systems:  make(map[string]SystemProperties),
			expected: false,
		},
		{
			name: "Non empty",
			systems: map[string]SystemProperties{
				"system1": SystemProperties{},
			},
			expected: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip, err := NewInfoProvider(tc.systems)
			require.NoError(err)
			require.Equal(tc.expected, ip.HasSystems())
		})
	}
}

func TestGetLD(t *testing.T) {
	require := require.New(t)
	tests := []struct {
		name       string
		systemName string
		want       float32
		Error      bool
	}{
		{name: "positive value", systemName: "sysA", want: 10.5, Error: false},
		{name: "zero value", systemName: "sysB", want: 0, Error: false},
		{name: "negative value", systemName: "sysC", want: 5.25, Error: false},
		{name: "ne value", systemName: "sysD", Error: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip, err := NewInfoProvider(testSystems())
			require.NoError(err)
			got, err := ip.GetLD(tc.systemName)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.want, got)
			}
		})
	}
}

func TestGetUD(t *testing.T) {
	require := require.New(t)
	tests := []struct {
		name       string
		systemName string
		want       float32
		Error      bool
	}{
		{name: "positive value", systemName: "sysA", want: 25.5, Error: false},
		{name: "zero value", systemName: "sysB", want: 0, Error: false},
		{name: "large value", systemName: "sysC", want: 99.99, Error: false},
		{name: "ne value", systemName: "sysD", Error: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip, err := NewInfoProvider(testSystems())
			require.NoError(err)
			got, err := ip.GetUD(tc.systemName)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.want, got)
			}
		})
	}
}

func TestGetPower(t *testing.T) {
	require := require.New(t)
	tests := []struct {
		name       string
		systemName string
		want       int
		Error      bool
	}{
		{name: "positive value", systemName: "sysA", want: 100, Error: false},
		{name: "zero value", systemName: "sysB", want: 0, Error: false},
		{name: "negative value", systemName: "sysC", want: 10, Error: false},
		{name: "ne value", systemName: "sysD", Error: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip, err := NewInfoProvider(testSystems())
			require.NoError(err)
			got, err := ip.GetPower(tc.systemName)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.want, got)
			}
		})
	}
}

func TestGetMaxVol(t *testing.T) {
	require := require.New(t)
	tests := []struct {
		name       string
		systemName string
		want       float32
		Error      bool
	}{
		{name: "positive value", systemName: "sysA", want: 500.0, Error: false},
		{name: "zero value", systemName: "sysB", want: 0, Error: false},
		{name: "fractional value", systemName: "sysC", want: 1234.5, Error: false},
		{name: "ne value", systemName: "sysD", Error: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip, err := NewInfoProvider(testSystems())
			require.NoError(err)
			got, err := ip.GetMaxVol(tc.systemName)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.want, got)
			}
		})
	}
}

func TestInfoProvider_GetCurrentVol(t *testing.T) {
	require := require.New(t)
	tests := []struct {
		name       string
		systems    map[string]SystemProperties
		systemName string
		heightCM   float32
		wantVal    float32
		Error      bool
	}{
		{
			name: "unknown system returns error",
			systems: map[string]SystemProperties{
				"sysA": {LD: 20, UD: 20, MaxHeight: 30, MaxVol: 9.42478},
			},
			systemName: "does-not-exist",
			heightCM:   10,
			wantVal:    0,
			Error:      true,
		},
		{
			name: "cylinder (LD == UD) at half height matches pi*r^2*h",
			systems: map[string]SystemProperties{
				"cylinder": {LD: 20, UD: 20, MaxHeight: 30, MaxVol: 9.42478},
			},
			systemName: "cylinder",
			heightCM:   15,
			wantVal:    4.71239, // pi*10^2*15 = 4712.39 cm^3 = 4.71239 L
			Error:      false,
		},
		{
			name: "cylinder (LD == UD) at full height matches MaxVol",
			systems: map[string]SystemProperties{
				"cylinder": {LD: 20, UD: 20, MaxHeight: 30, MaxVol: 9.42478},
			},
			systemName: "cylinder",
			heightCM:   30,
			wantVal:    9.42478,
			Error:      false,
		},
		{
			name: "bucket-shaped frustum (LD < UD) at full height matches formula",
			systems: map[string]SystemProperties{
				"bucket": {LD: 20, UD: 30, MaxHeight: 30, MaxVol: 14.92257},
			},
			systemName: "bucket",
			heightCM:   30,
			wantVal:    14.92257,
			Error:      false,
		},
		{
			name: "bucket-shaped frustum (LD < UD) at partial height",
			systems: map[string]SystemProperties{
				"bucket": {LD: 20, UD: 30, MaxHeight: 30, MaxVol: 14.92257},
			},
			systemName: "bucket",
			heightCM:   15,
			wantVal:    5.98875,
			Error:      false,
		},
		{
			name: "inverted taper (LD > UD, e.g. pot narrower at rim) at partial height",
			systems: map[string]SystemProperties{
				"pot": {LD: 30, UD: 20, MaxHeight: 20, MaxVol: 15},
			},
			systemName: "pot",
			heightCM:   10,
			wantVal:    5.95113,
			Error:      false,
		},
		// Cross-validation cases below are taken from an independent
		// implementation of the same frustum formula:
		// https://rechneronline.de/litre/bucket.php
		{
			name: "rechneronline: bottom 30cm, top 33cm, height 25cm, full - ~19.5L",
			systems: map[string]SystemProperties{
				"rol-bucket-1": {LD: 30, UD: 33, MaxHeight: 25, MaxVol: 19.4977},
			},
			systemName: "rol-bucket-1",
			heightCM:   25,
			wantVal:    19.4977,
			Error:      false,
		},
		{
			name: "rechneronline: bottom 30cm, top 33cm, height 25cm, filled to 13cm - ~9.7L",
			systems: map[string]SystemProperties{
				"rol-bucket-1": {LD: 30, UD: 33, MaxHeight: 25, MaxVol: 19.4977},
			},
			systemName: "rol-bucket-1",
			heightCM:   13,
			wantVal:    9.67318, // fill width 31.56cm per source, matches our r(h) interpolation
			Error:      false,
		},
		{
			name: "rechneronline: bottom 19.5cm, top 25.8cm, height 25cm, full - ~10.1L",
			systems: map[string]SystemProperties{
				"rol-bucket-2": {LD: 19.5, UD: 25.8, MaxHeight: 25, MaxVol: 10.13812},
			},
			systemName: "rol-bucket-2",
			heightCM:   25,
			wantVal:    10.13812,
			Error:      false,
		},
		{
			name: "rechneronline: bottom 19.5cm, top 25.8cm, height 25cm, filled to 20cm - ~3/4 full",
			systems: map[string]SystemProperties{
				"rol-bucket-2": {LD: 19.5, UD: 25.8, MaxHeight: 25, MaxVol: 10.13812},
			},
			systemName: "rol-bucket-2",
			heightCM:   20,
			wantVal:    7.64975,
			Error:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip, err := NewInfoProvider(tc.systems)
			require.NoError(err)
			got, err := ip.GetCurrentVol(tc.systemName, tc.heightCM)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.InDelta(tc.wantVal, got, 0.01) // tolerance: 10 mL
			}
		})
	}
}
