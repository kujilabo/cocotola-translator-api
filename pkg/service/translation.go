package service

import (
	"errors"
)

var ErrTranslationNotFound = errors.New("translation not found")
var ErrTranslationAlreadyExists = errors.New("custsomtranslation already exists")
