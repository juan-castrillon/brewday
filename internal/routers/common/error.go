package common

import "errors"

var ErrNoRecipeLoaded = errors.New("no recipe loaded")

var ErrNoRecipeIDProvided = errors.New("no recipe id provided")
