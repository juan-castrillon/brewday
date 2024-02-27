package sql

import (
	"brewday/internal/recipe"
	"encoding/json"
)

type MarshalResult struct {
	StatusParams string
	MashingMalts string
	MashingRasts string
	HopHops      string
	HopAdd       string
	FermAdd      string
	Yeast        string
}

type UnmarshalResult struct {
	StatusParams []string
	MashingMalts []recipe.Malt
	MashingRasts []recipe.Rast
	HopHops      []recipe.Hops
	HopAdd       []recipe.AdditionalIngredient
	FermAdd      []recipe.AdditionalIngredient
	Yeast        recipe.Yeast
}

func (s *PersistentStore) marshalStructs(r *recipe.Recipe) (*MarshalResult, error) {
	_, statusParams := r.GetStatus()
	sP, err := s.marshalStatusParams(statusParams...)
	if err != nil {
		return nil, err
	}
	maltBytes, err := json.Marshal(r.Mashing.Malts)
	if err != nil {
		return nil, err
	}
	rastBytes, err := json.Marshal(r.Mashing.Rasts)
	if err != nil {
		return nil, err
	}
	hopBytes, err := json.Marshal(r.Hopping.Hops)
	if err != nil {
		return nil, err
	}
	hopAddBytes, err := json.Marshal(r.Hopping.AdditionalIngredients)
	if err != nil {
		return nil, err
	}
	fermAddBytes, err := json.Marshal(r.Fermentation.AdditionalIngredients)
	if err != nil {
		return nil, err
	}
	yeastBytes, err := json.Marshal(r.Fermentation.Yeast)
	if err != nil {
		return nil, err
	}
	return &MarshalResult{
		StatusParams: sP,
		MashingMalts: string(maltBytes),
		MashingRasts: string(rastBytes),
		HopHops:      string(hopBytes),
		HopAdd:       string(hopAddBytes),
		FermAdd:      string(fermAddBytes),
		Yeast:        string(yeastBytes),
	}, nil
}

func (s *PersistentStore) marshalStatusParams(statusParams ...string) (string, error) {
	statusParamsByte, err := json.Marshal(statusParams)
	if err != nil {
		return "", err
	}
	return string(statusParamsByte), nil
}

func (s *PersistentStore) unmarshalStructs(m *MarshalResult) (*UnmarshalResult, error) {
	var statusParams []string
	var mashingMalts []recipe.Malt
	var mashingRasts []recipe.Rast
	var hopHops []recipe.Hops
	var hopAdd []recipe.AdditionalIngredient
	var fermAdd []recipe.AdditionalIngredient
	var yeast recipe.Yeast
	err := json.Unmarshal([]byte(m.StatusParams), &statusParams)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(m.MashingMalts), &mashingMalts)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(m.MashingRasts), &mashingRasts)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(m.HopHops), &hopHops)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(m.HopAdd), &hopAdd)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(m.FermAdd), &fermAdd)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(m.Yeast), &yeast)
	if err != nil {
		return nil, err
	}
	return &UnmarshalResult{
		StatusParams: statusParams,
		MashingMalts: mashingMalts,
		MashingRasts: mashingRasts,
		HopHops:      hopHops,
		HopAdd:       hopAdd,
		FermAdd:      fermAdd,
		Yeast:        yeast,
	}, nil
}
