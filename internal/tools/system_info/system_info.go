package systeminfo

import "errors"

type InfoProvider struct {
	systems map[string]*SystemProperties
}

func NewInfoProvider(i map[string]*SystemProperties) (*InfoProvider, error) {
	return &InfoProvider{systems: i}, nil

}

func (ip *InfoProvider) GetLD(systemName string) float32 {
	return ip.systems[systemName].LD
}

func (ip *InfoProvider) GetUD(systemName string) float32 {
	return ip.systems[systemName].UD
}

func (ip *InfoProvider) GetPower(systemName string) int {
	return ip.systems[systemName].Power
}

func (ip *InfoProvider) GetMaxVol(systemName string) float32 {
	return ip.systems[systemName].MaxVol
}

func (ip *InfoProvider) GetCurrentVol(systemName string, heightCM float32) (float32, error) {
	return 0, errors.New("Not implemented")
}
