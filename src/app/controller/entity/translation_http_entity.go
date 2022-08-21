package entity

type TranslationFindParameterHTTPEntity struct {
	Letter string `json:"letter"`
}

type TranslationHTTPEntity struct {
	Lang2      string `json:"lang2"`
	Text       string `json:"text"`
	Pos        int    `json:"pos"`
	Translated string `json:"translated"`
	Provider   string `json:"provider"`
}

type TranslationFindResponseHTTPEntity struct {
	Results []TranslationHTTPEntity `json:"results"`
}

type TranslationAddParameterHTTPEntity struct {
	Lang2      string `json:"lang2" binding:"required"`
	Text       string `json:"text" binding:"required"`
	Pos        int    `json:"pos" binding:"required"`
	Translated string `json:"translated" binding:"required"`
}

type TranslationUpdateParameterHTTPEntity struct {
	Translated string `json:"translated" binding:"required"`
}
