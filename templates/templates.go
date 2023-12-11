/*
Package templates initializes all the templates for use in services.
*/
package templates

import "html/template"

var Tpl *template.Template

func init() {
	Tpl = template.Must(template.ParseGlob("../../templates/*"))
}
