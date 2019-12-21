package config

//swagger:model
type Template struct {
	// Имя шаблона (имя каталога)
	// Required: true
	// Example: ascii
	Name string `json:"name"`
}
