// Package systeminfo provides access to static and derived properties of
// physical brewing systems such as their diameters, power rating, and volume.
package systeminfo

import (
	"errors"
	"math"
)

// InfoProvider exposes lookup methods for system properties, keyed by
// system name.
type InfoProvider struct {
	// systems maps a system name to its properties. A nil or missing
	// entry means the system is unknown to this provider.
	systems map[string]SystemProperties
}

// NewInfoProvider creates an InfoProvider backed by the given map of
// system properties, keyed by system name. The map is used directly
// (not copied), so mutating it after construction will affect the
// InfoProvider's view of the data.
//
// The returned error is currently always nil; it is reserved for future
// validation of the input map.
func NewInfoProvider(i map[string]SystemProperties) (*InfoProvider, error) {
	return &InfoProvider{systems: i}, nil

}

// HasSystems can be used as a check before calling methods to verify this info provider
// has at least one brewing system configured
func (ip *InfoProvider) HasSystems() bool {
	return len(ip.systems) > 0
}

// GetLD returns the lower diameter (LD) of the named system, in the same
// units as stored on SystemProperties (e.g. cm).
//
// It returns an error if systemName is not known to the provider.
func (ip *InfoProvider) GetLD(systemName string) (float32, error) {
	val, ok := ip.systems[systemName]
	if !ok {
		return 0, errors.New("System not found")
	}
	return val.LD, nil
}

// GetUD returns the upper diameter (UD) of the named system, in the same
// units as stored on SystemProperties (e.g. cm).
//
// It returns an error if systemName is not known to the provider.
func (ip *InfoProvider) GetUD(systemName string) (float32, error) {
	val, ok := ip.systems[systemName]
	if !ok {
		return 0, errors.New("System not found")
	}
	return val.UD, nil
}

// GetPower returns the power rating of the named system.
//
// It returns an error if systemName is not known to the provider.
func (ip *InfoProvider) GetPower(systemName string) (int, error) {
	val, ok := ip.systems[systemName]
	if !ok {
		return 0, errors.New("System not found")
	}
	return val.Power, nil
}

// GetMaxVol returns the maximum volume capacity of the named system.
//
// It returns an error if systemName is not known to the provider.
func (ip *InfoProvider) GetMaxVol(systemName string) (float32, error) {
	val, ok := ip.systems[systemName]
	if !ok {
		return 0, errors.New("System not found")
	}
	return val.MaxVol, nil
}

// GetCurrentVol calculates the volume of liquid, in liters, held in the
// named system at a given fill height (heightCM). The vessel is modeled
// as a frustum (truncated cone) defined by its lower diameter (LD),
// upper diameter (UD), and full height (MaxHeightCM), all in cm. A
// cylindrical vessel is handled automatically, since it is the special
// case LD == UD.
//
// The formula follows the standard conical frustum volume derivation:
// radius is interpolated linearly with height, r(z) = R1 + (R2-R1)*(z/h),
// and integrated over the circular cross-section, giving
// V = (pi*h/3) * (R1^2 + R1*R2 + R2^2) for the full frustum, and the same
// form with R2 replaced by r(heightCM) for a partial fill. The result is
// in cm^3, so it is divided by 1000 to convert to liters.
// Source: Weisstein, Eric W., "Conical Frustum," CRC Standard
// Mathematical Tables, 28th ed. (Beyer 1987, pp. 129-130, 133):
// https://archive.lib.msu.edu/crcmath/math/math/c/c581.htm
//
// heightCM is assumed to be a valid, non-zero fill height that does not
// exceed the vessel's MaxHeightCM; callers are responsible for enforcing
// that invariant (e.g. at config-parsing time) before calling this
// method.
//
// It returns an error if systemName is not known to the provider.
func (ip *InfoProvider) GetCurrentVol(systemName string, heightCM float32) (float32, error) {
	val, ok := ip.systems[systemName]
	if !ok {
		return 0, errors.New("System not found")
	}

	h := float64(heightCM)
	maxH := float64(val.MaxHeight)
	r1 := float64(val.LD) / 2
	r2 := float64(val.UD) / 2

	// Radius at the fill height, linearly interpolated between the base
	// radius (r1) and the top radius (r2) based on how far up the
	// vessel's full height heightCM reaches.
	rh := r1 + (r2-r1)*(h/maxH)

	// Volume of a frustum from 0 to heightCM, in cm^3.
	volCM3 := (math.Pi * h / 3) * (r1*r1 + r1*rh + rh*rh)

	// Convert cm^3 to liters (1 L = 1000 cm^3).
	const cm3PerLiter = 1000
	return float32(volCM3 / cm3PerLiter), nil
}
