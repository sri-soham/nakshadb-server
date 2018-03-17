package lib

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

var timestamp = "20170505"

func ParseTemplates(template_root string) (*template.Template, error) {
	tmpl := template.New("")
	func_map := getTemplateFuncs()
	tmpl.Funcs(func_map)

	err := filepath.Walk(template_root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			_, err = tmpl.ParseFiles(path)
		}

		return err
	})

	return tmpl, err
}

func getTemplateFuncs() template.FuncMap {
	func_map := template.FuncMap{
		"_title":        func() string { return "NakshaDB" },
		"_js":           js,
		"_css":          css,
		"_asset":        asset,
		"_image":        image,
		"_db_timestamp": db_timestamp,
	}

	return func_map
}

func asset(path string) string {
	path = strings.Trim(path, " /")
	path = "/assets/" + timestamp + "/" + path
	return path
}

func js(path string) template.HTML {
	path = strings.Trim(path, " /")
	path = "js/" + path
	str := fmt.Sprintf("<script type=\"text/javascript\" src=\"%v\"></script>", asset(path))

	return template.HTML(str)
}

func css(path string) template.HTML {
	path = strings.Trim(path, " /")
	path = "css/" + path
	str := fmt.Sprintf("<link type=\"text/css\" rel=\"stylesheet\" href=\"%v\" />", asset(path))

	return template.HTML(str)
}

func image(path string, alt string, other_attrs string) template.HTML {
	path = strings.Trim(path, " /")
	path = "images/" + path
	str := fmt.Sprintf("<img src=\"%v\" alt=\"%v\" %v/>", asset(path), alt, other_attrs)

	return template.HTML(str)
}

func db_timestamp(ts string) string {
	var res string
	if len(ts) == 0 {
		res = ""
	} else {
		res = strings.Split(ts, ".")[0]
	}

	return res
}

func SetAssetTimestamp(ts string) {
	timestamp = ts
}
