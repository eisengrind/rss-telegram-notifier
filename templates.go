package main

import (
	"html/template"
	"log"
)

func init() {
	var err error
	parsedTemplateMessage, err = template.New("TemplateMessage").Parse(templateMessage)
	if err != nil {
		log.Fatal(err.Error())
	}
}

var (
	parsedTemplateMessage *template.Template
	templateMessage       = `
<b>{{.Title}}</b>
<i>{{.Date}} from {{.Author}}</i>

{{.Description}}

<a href="{{.Link}}">More Information</a>
`
)
