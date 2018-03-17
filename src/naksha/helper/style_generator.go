package helper

import (
	"fmt"
	"net/http"
)

type StyleGenerator struct {
	req *http.Request
}

func (sg *StyleGenerator) PolygonStyle() string {
	poly := fmt.Sprintf(
		"<PolygonSymbolizer fill=\"%v\" fill-opacity=\"%v\" />",
		sg.req.PostFormValue("fill"),
		sg.req.PostFormValue("fill_opacity"),
	)
	line := fmt.Sprintf(
		"<LineSymbolizer stroke=\"%v\" stroke-width=\"%v\" stroke-opacity=\"%v\" />",
		sg.req.PostFormValue("stroke"),
		sg.req.PostFormValue("stroke_width"),
		sg.req.PostFormValue("stroke_opacity"),
	)

	return fmt.Sprintf("<Rule>%v%v</Rule>", poly, line)
}

func (sg *StyleGenerator) LineStringStyle() string {
	line := fmt.Sprintf(
		"<LineSymbolizer stroke=\"%v\" stroke-width=\"%v\" stroke-opacity=\"%v\" />",
		sg.req.PostFormValue("stroke"),
		sg.req.PostFormValue("stroke_width"),
		sg.req.PostFormValue("stroke_opacity"),
	)

	return fmt.Sprintf("<Rule>%v</Rule>", line)
}

func (sg *StyleGenerator) PointStyle() string {
	symbol := "<MarkersSymbolizer fill=\"%v\" opacity=\"%v\" " +
		"stroke=\"%v\" stroke-width=\"%v\" stroke-opacity=\"%v\" " +
		"width=\"%v\" height=\"%v\" placement=\"point\" marker-type=\"ellipse\" />"
	marker := fmt.Sprintf(
		symbol,
		sg.req.PostFormValue("fill"),
		sg.req.PostFormValue("fill_opacity"),
		sg.req.PostFormValue("stroke"),
		sg.req.PostFormValue("stroke_width"),
		sg.req.PostFormValue("stroke_opacity"),
		sg.req.PostFormValue("width"),
		sg.req.PostFormValue("height"),
	)

	return fmt.Sprintf("<Rule>%v</Rule>", marker)
}

func MakeStyleGenerator(req *http.Request) StyleGenerator {
	return StyleGenerator{req}
}
