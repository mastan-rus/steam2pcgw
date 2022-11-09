package main

type Game map[string]GameValue

type GameValue struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

type Data struct {
	Type                string          `json:"type"`
	Name                string          `json:"name"`
	SteamAppid          int64           `json:"steam_appid"`
	IsFree              bool            `json:"is_free"`
	ControllerSupport   *string         `json:"controller_support,omitempty"`
	Dlc                 []int64         `json:"dlc,omitempty"`
	DetailedDescription string          `json:"detailed_description"`
	AboutTheGame        string          `json:"about_the_game"`
	ShortDescription    string          `json:"short_description"`
	SupportedLanguages  string          `json:"supported_languages"`
	Website             *string         `json:"website"`
	PCRequirements      Requirement     `json:"pc_requirements,omitempty"`
	MACRequirements     Requirement     `json:"mac_requirements,omitempty"`
	LinuxRequirements   Requirement     `json:"linux_requirements,omitempty"`
	Developers          []string        `json:"developers"`
	Publishers          []string        `json:"publishers"`
	Packages            []int64         `json:"packages"`
	Platforms           Platforms       `json:"platforms"`
	Metacritic          *Metacritic     `json:"metacritic,omitempty"`
	Categories          []Category      `json:"categories"`
	Genres              []Genre         `json:"genres"`
	Recommendations     Recommendations `json:"recommendations"`
	ReleaseDate         ReleaseDate     `json:"release_date"`
	SupportInfo         SupportInfo     `json:"support_info"`
}

type Category struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

type Genre struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type Requirement map[string]interface{}

type Metacritic struct {
	Score int64  `json:"score"`
	URL   string `json:"url"`
}

type Sub struct {
	Packageid                int64  `json:"packageid"`
	PercentSavingsText       string `json:"percent_savings_text"`
	PercentSavings           int64  `json:"percent_savings"`
	OptionText               string `json:"option_text"`
	OptionDescription        string `json:"option_description"`
	CanGetFreeLicense        string `json:"can_get_free_license"`
	IsFreeLicense            bool   `json:"is_free_license"`
	PriceInCentsWithDiscount int64  `json:"price_in_cents_with_discount"`
}

type Platforms struct {
	Windows bool `json:"windows"`
	MAC     bool `json:"mac"`
	Linux   bool `json:"linux"`
}

type Recommendations struct {
	Total int64 `json:"total"`
}

type ReleaseDate struct {
	ComingSoon bool   `json:"coming_soon"`
	Date       string `json:"date"`
}

type SupportInfo struct {
	URL   string `json:"url"`
	Email string `json:"email"`
}

type Language map[string]LanguageValue

type LanguageValue struct {
	UI        bool
	Audio     bool
	Subtitles bool
}
