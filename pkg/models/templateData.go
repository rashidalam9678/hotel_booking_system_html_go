package models

//TemplateData holds the data which has to be sent to templates
type TemplateData struct{
	StringMap map[string]string
	IntMap map[string]int
	FloatMap map[string]float32
	Data map[string]interface{}
	CSRFToken string
	Error string
	Warning string
	Flash string
}