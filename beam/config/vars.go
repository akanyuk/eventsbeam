package config

//swagger:model
type Template struct {
	// Имя шаблона (имя каталога)
	// Required: true
	// Example: ascii
	Name string `json:"name"`
	// Параметры шаблона
	Params map[string]TemplateParam `json:"params"`
}

type TemplateParam struct {
	// Заголовок параметра
	// Example: varX
	Title string `json:"title"`
}
