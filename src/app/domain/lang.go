package domain

import "golang.org/x/xerrors"

const Lang2Len = 2
const Lang3Len = 3
const Lang5Len = 5

type Lang2 interface {
	String() string
	ToLang3() string
}

type lang2 struct {
	value string
}

func NewLang2(lang string) (Lang2, error) {
	if len(lang) != Lang2Len {
		return nil, xerrors.Errorf("invalid parameter. Lang2: %s", lang)
	}

	return &lang2{
		value: lang,
	}, nil
}

func (l *lang2) String() string {
	return l.value
}

func (l *lang2) ToLang3() string {
	switch l.value {
	case "en":
		return "eng"
	case "ja":
		return "jpn"
	default:
		return "___"
	}
}

type Lang3 interface {
	String() string
}

type lang3 struct {
	value string
}

func NewLang3(lang string) (Lang3, error) {
	if len(lang) != Lang3Len {
		return nil, xerrors.Errorf("invalid parameter. Lang3: %s", lang)
	}

	return &lang3{
		value: lang,
	}, nil
}

func (l *lang3) String() string {
	return l.value
}

type Lang5 interface {
	String() string
	ToLang2() string
}

type lang5 struct {
	value string
}

func NewLang5(lang string) (Lang5, error) {
	if len(lang) != Lang5Len {
		return nil, xerrors.Errorf("invalid parameter. Lang5: %s", lang)
	}

	return &lang5{
		value: lang,
	}, nil
}

func (l *lang5) String() string {
	return l.value
}

func (l *lang5) ToLang2() string {
	switch l.value {
	case "eng":
		return "en"
	case "jpn":
		return "ja"
	default:
		return "__"
	}
}
