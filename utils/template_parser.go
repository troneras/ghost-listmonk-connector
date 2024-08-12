package utils

import (
	"bytes"
	"html/template"
)

func ParseTemplate(templateString string, postData map[string]interface{}) (string, error) {
	tmpl, err := template.New("email").Parse(templateString)
	if err != nil {
		return "", err
	}

	// Convert the HTML content to template.HTML type
	if html, ok := postData["Html"].(string); ok {
		postData["Html"] = template.HTML(html)
	}

	data := map[string]interface{}{
		"Post": postData,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
