/*
Package templates initializes all the templates for use in services.
*/
package templates

import (
	"html/template"

	"attendance.com/src/logger"
)

var Tpl *template.Template

func init() {
	logger.Println("Initializing templates...")
	Tpl = template.Must(template.ParseGlob("./templates/*.gohtml"))
	logger.Println("Templates ready!")

}
