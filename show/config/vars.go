package config

//swagger:model
type Slide struct {
	Position int `json:"-"`
	// Идентификатор слайда
	// Readonly: true
	// Min: 1
	// Example: 1
	ID int `json:"id"`
	// Идентификатор компо. Если не указано, то используется глобальный список слайдов.
	// Example: zxdemo
	Compo string `json:"compo"`
	// Имя шаблона, имеющегося в системе
	// Required: true
	// Example: ansi
	Template string `json:"template"`
	// Параметры, передаваемые в шаблон при отображении слайда
	// Example: { "title": "Super demo", "platform": "ZX Spectrum" }
	Params map[string]string `json:"params"`
}

//swagger:model
type Compo struct {
	// Идентификатор компо
	// Required: true
	// Example: zxdemo
	Alias string `json:"alias"`
	// Наименование компо
	// Required: true
	// Example: ZX Demo
	Title string `json:"title"`
}

type Show struct {
	/*
		EventsURL      string
		EventsUsername string
		EventsPassword string

		GeneralSlides []Slide
	*/

	Compos []Compo
}
