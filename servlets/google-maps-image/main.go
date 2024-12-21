// Note: run `go doc -all` in this package to see all of the types and functions available.
// ./pdk.gen.go contains the domain types from the host where your plugin will run.
package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/extism/go-pdk"
)

var (
	apiKey string
)

func loadKeys() {
	if apiKey != "" { // already loaded
		return
	}
	k, kerr := pdk.GetConfig("api-key")
	if !kerr {
		panic("missing required configuration")
	}
	apiKey = k
}

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
// It takes CallToolRequest as input (The incoming tool request from the LLM)
// And returns CallToolResult (The servlet's response to the given tool call)
func Call(input CallToolRequest) (CallToolResult, error) {
	loadKeys()

	args := input.Params.Arguments.(map[string]any)
	switch input.Params.Name {
	case "google-maps-static-api-center":
		return fetchStaticImage(args)
	case "google-maps-static-api-markers":
		return fetchStaticImage(args)
	default:
		return callToolError("unknown tool " + input.Params.Name), nil
	}
}

func fetchStaticImage(args map[string]any) (CallToolResult, error) {
	center, ok := args["center"].(string)
	if !ok {
		return callToolError("center is required"), nil
	}
	zoom, ok := args["zoom"].(float64)
	if !ok {
		return callToolError("zoom is required"), nil
	}
	size, ok := args["size"].(string)
	if !ok {
		return callToolError("size is required"), nil
	}
	maptype, ok := args["maptype"].(string)
	if !ok {
		maptype = "roadmap"
	}
	q := url.Values{}

	q.Set("format", "png")
	q.Set("key", apiKey)

	q.Set("center", center)
	q.Set("zoom", fmt.Sprint(zoom))
	q.Set("size", size)
	q.Set("scale", "2")
	q.Set("maptype", maptype)

	if markers, ok := args["markers"].([]any); ok {
		for _, marker := range markers {
			if m, ok := marker.(map[string]any); ok {
				params := []string{}

				if label, lab := m["label"]; lab {
					params = append(params, fmt.Sprintf("label:%s", label))
				}
				if color, col := m["color"]; col {
					color = strings.Replace(fmt.Sprint(color), "#", "0x", -1)
					params = append(params, fmt.Sprintf("color:%s", color))
				}
				if size, siz := m["size"]; siz {
					params = append(params, fmt.Sprintf("size:%s", size))
				}
				if icon, ico := m["icon"]; ico {
					params = append(params, fmt.Sprintf("icon:%s", icon))
				}
				if location, lok := m["location"]; lok {
					params = append(params, fmt.Sprint(location))
				}
				q.Add("markers", strings.Join(params, "|"))
			}
		}
	}
	if path, ok := args["path"].([]any); ok {
		for _, p := range path {
			if pa, ok := p.(map[string]any); ok {
				params := []string{}
				if color, col := pa["color"]; col {
					color = strings.Replace(fmt.Sprint(color), "#", "0x", -1)
					params = append(params, fmt.Sprintf("color:%s", color))
				}
				if weight, wei := pa["weight"]; wei {
					params = append(params, fmt.Sprintf("weight:%f", weight))
				}
				if fillcolor, fil := pa["fillcolor"]; fil {
					fillcolor = strings.Replace(fmt.Sprint(fillcolor), "#", "0x", -1)
					params = append(params, fmt.Sprintf("fillcolor:%s", fillcolor))
				}
				if geodesic, geo := pa["geodesic"]; geo {
					params = append(params, fmt.Sprintf("geodesic:%t", geodesic))
				}
				if points, poi := pa["points"]; poi {
					params = append(params, fmt.Sprint(points))
				}
				q.Add("path", strings.Join(params, "|"))
			}
		}
	}
	if style, ok := args["style"].([]any); ok {
		for _, s := range style {
			if st, ok := s.(map[string]any); ok {
				params := []string{}
				if feature, fea := st["feature"]; fea {
					params = append(params, fmt.Sprintf("feature:%s", feature))
				}
				if element, ele := st["element"]; ele {
					params = append(params, fmt.Sprintf("element:%s", element))
				}
				if color, col := st["color"]; col {
					color = strings.Replace(fmt.Sprint(color), "#", "0x", -1)
					params = append(params, fmt.Sprintf("color:%s", color))
				}
				if weight, wei := st["weight"]; wei {
					params = append(params, fmt.Sprintf("weight:%f", weight))
				}
				if visibility, vis := st["visibility"]; vis {
					params = append(params, fmt.Sprintf("visibility:%s", visibility))
				} else {
					params = append(params, "visibility:on")
				}
				q.Add("style", strings.Join(params, "|"))
			}
		}
	}
	if visible, ok := args["visible"].([]any); ok {
		for _, v := range visible {
			if vi, ok := v.(map[string]any); ok {
				params := []string{}
				if location, loc := vi["location"]; loc {
					params = append(params, fmt.Sprint(location))
				}
				q.Add("visible", strings.Join(params, "|"))
			}
		}
	}

	url := fmt.Sprint("https://maps.googleapis.com/maps/api/staticmap?", q.Encode())

	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	res := req.Send()

	return callToolSuccess(res.Body()), nil
}

func Describe() (ListToolsResult, error) {
	return ListToolsResult{Tools: []ToolDescription{
		{
			Name: "google-maps-static-api-center",
			Description: "Returns an image, centered on the given `center`." +
				"For each request, you can specify the location of the map, the size of the image, the zoom level, " +
				"the type of map, and the placement of optional markers at locations on the map. " +
				"You can additionally label your markers using alphanumeric characters.",
			InputSchema: schema{
				"type": "object",
				"properties": props{
					"center":   prop("string", "(required)  defines the center of the map, equidistant from all edges of the map. This parameter takes a location as either a comma-separated {latitude,longitude} pair (e.g. '40.714728,-73.998672') or a string address (e.g. 'city hall, new york, ny') identifying a unique location on the face of the earth."),
					"zoom":     prop("number", "(required) defines the zoom level of the map, which determines the magnification level of the map. This parameter takes a numerical value corresponding to the zoom level of the region desired."),
					"size":     prop("string", "(required) defines the rectangular dimensions of the map image. This parameter takes a string of the form {horizontal_value}x{vertical_value}. For example, 500x400 defines a map 500 pixels wide by 400 pixels high."),
					"maptype":  prop("string", "(optional) defines the type of map to construct. There are several possible maptype values, including roadmap, satellite, hybrid, and terrain."),
					"language": prop("string", "(optional) defines the language to use for display of labels on map tiles."),
					"region":   prop("string", "(optional) defines the appropriate borders to display, based on geo-political sensitivities. Accepts a region code specified as a two-character ccTLD ('top-level domain') value. "),
					"map_id":   prop("string", "(optional) specifies the identifier for a specific map. The Map ID associates a map with a particular style or feature, and must belong to the same project as the API key used to initialize the map."),
					"markers": arrprop("array", "(optional) define one or more markers to attach to the image at specified locations.", schema{
						"type": "object",
						"properties": props{
							"label":    prop("string", "(optional) defines a single uppercase alphanumeric character to be displayed within the marker. Note that default and mid sized markers are the only markers capable of displaying an alphanumeric-character parameter. tiny and small markers are not capable of displaying an alphanumeric-character."),
							"color":    prop("string", "(optional) defines a color for the marker. Accepts a label from the set {black, brown, green, purple, yellow, blue, gray, orange, red, white}"),
							"size":     prop("string", "(optional) defines the size of the marker. Accepts a label of {tiny, mid, small, normal, large}"),
							"icon":     prop("string", "(optional) defines a custom icon to display at the marker's location. The provided icon must be a publicly accessible URL that does not require HTTP authentication."),
							"location": prop("string", "(required) defines the location of the marker on the map. This parameter takes a location as either a comma-separated {latitude,longitude} pair (e.g. '40.714728,-73.998672') or a string address (e.g. 'city hall, new york, ny') identifying a unique location on the face of the earth."),
						},
					}),
					"path": arrprop("array", "(optional) defines a single path of two or more connected points to overlay on the image at specified locations.", schema{
						"type": "object",
						"properties": props{
							"color":     prop("string", "(optional) defines a color for the path. Accepts a label from the set {black, brown, green, purple, yellow, blue, gray, orange, red, white}"),
							"weight":    prop("number", "(optional) defines the weight of the path in pixels."),
							"fillcolor": prop("string", "(optional) defines the fill color of the path. This parameter takes a color in hexadecimal format with a leading # sign (e.g. #000000)."),
							"geodesic":  prop("boolean", "(optional) specifies whether to draw each segment of the path as a geodesic (true) or as a straight line on the Mercator projection (false)."),
							"points":    prop("string", "(required) defines a single path of two or more connected points to overlay on the image at specified locations. This parameter takes a string of point definitions separated by the pipe character (|), or an encoded polyline using the enc: prefix within the location declaration of the path."),
						},
					}),
					"visible": arrprop("array", "(optional) specifies one or more locations that should remain visible on the map, though no markers or other indicators will be displayed.", schema{
						"type": "object",
						"properties": props{
							"location": prop("string", "(required) defines the location of the marker on the map. This parameter takes a location as either a comma-separated {latitude,longitude} pair (e.g. '40.714728,-73.998672') or a string address (e.g. 'city hall, new york, ny') identifying a unique location on the face of the earth."),
						},
					}),
					"style": arrprop("array", "(optional) defines a custom style to alter the presentation of a specific feature (roads, parks, and other features) of the map.", schema{
						"type": "object",
						"properties": props{
							"feature":    prop("string", "(required) defines the feature to style, applying the style rules to all points of the feature. Accepted values are 'all' and 'road'."),
							"element":    prop("string", "(required) defines the element to which the style is applied. An element is a feature of a map, and may be one of the following values: {all, geometry, labels}."),
							"color":      prop("string", "(optional) defines the color of the feature. This can be specified in hexadecimal (as #RRGGBB, #RRGGBBAA, #RGB, or #RGBA), or in HSLA notation. For example, '#ffcc00' will set a golden yellow color."),
							"weight":     prop("number", "(optional) defines the weight of the feature, in pixels. This setting affects polylines, where the value width is measured in pixels."),
							"visibility": prop("string", "(optional) defines the visibility of the feature. The setting can be 'on', 'off', or 'simplified'."),
						},
					}),
				},
				"required": []string{"center", "zoom", "size"},
			},
		},
		{
			Name: "google-maps-static-api-markers",
			Description: "Returns an image, centered around the given `markers`. This should be preferred when there is a path or multiple markers. " +
				"For each request, you can specify the location of the map, the size of the image, the zoom level, " +
				"the type of map, and the placement of optional markers at locations on the map. " +
				"You can additionally label your markers using alphanumeric characters.",
			InputSchema: schema{
				"type": "object",
				"properties": props{
					"center":   prop("string", "(optional)  defines the center of the map, equidistant from all edges of the map. This parameter takes a location as either a comma-separated {latitude,longitude} pair (e.g. '40.714728,-73.998672') or a string address (e.g. 'city hall, new york, ny') identifying a unique location on the face of the earth."),
					"zoom":     prop("number", "(required) defines the zoom level of the map, which determines the magnification level of the map. This parameter takes a numerical value corresponding to the zoom level of the region desired."),
					"size":     prop("string", "(required) defines the rectangular dimensions of the map image. This parameter takes a string of the form {horizontal_value}x{vertical_value}. For example, 500x400 defines a map 500 pixels wide by 400 pixels high."),
					"maptype":  prop("string", "(optional) defines the type of map to construct. There are several possible maptype values, including roadmap, satellite, hybrid, and terrain."),
					"language": prop("string", "(optional) defines the language to use for display of labels on map tiles."),
					"region":   prop("string", "(optional) defines the appropriate borders to display, based on geo-political sensitivities. Accepts a region code specified as a two-character ccTLD ('top-level domain') value. "),
					"map_id":   prop("string", "(optional) specifies the identifier for a specific map. The Map ID associates a map with a particular style or feature, and must belong to the same project as the API key used to initialize the map."),
					"markers": arrprop("array", "(required) define one or more markers to attach to the image at specified locations.", schema{
						"type": "object",
						"properties": props{
							"label":    prop("string", "(optional) defines a single uppercase alphanumeric character to be displayed within the marker. Note that default and mid sized markers are the only markers capable of displaying an alphanumeric-character parameter. tiny and small markers are not capable of displaying an alphanumeric-character."),
							"color":    prop("string", "(optional) defines a color for the marker. Accepts a label from the set {black, brown, green, purple, yellow, blue, gray, orange, red, white}"),
							"size":     prop("string", "(optional) defines the size of the marker. Accepts a label of {tiny, mid, small, normal, large}"),
							"icon":     prop("string", "(optional) defines a custom icon to display at the marker's location. The provided icon must be a publicly accessible URL that does not require HTTP authentication."),
							"location": prop("string", "(required) defines the location of the marker on the map. This parameter takes a location as either a comma-separated {latitude,longitude} pair (e.g. '40.714728,-73.998672') or a string address (e.g. 'city hall, new york, ny') identifying a unique location on the face of the earth."),
						},
					}),
					"path": arrprop("array", "(optional) defines a single path of two or more connected points to overlay on the image at specified locations.", schema{
						"type": "object",
						"properties": props{
							"color":     prop("string", "(optional) defines a color for the path. Accepts a label from the set {black, brown, green, purple, yellow, blue, gray, orange, red, white}"),
							"weight":    prop("number", "(optional) defines the weight of the path in pixels."),
							"fillcolor": prop("string", "(optional) defines the fill color of the path. This parameter takes a color in hexadecimal format with a leading # sign (e.g. #000000)."),
							"geodesic":  prop("boolean", "(optional) specifies whether to draw each segment of the path as a geodesic (true) or as a straight line on the Mercator projection (false)."),
							"points":    prop("string", "(required) defines a single path of two or more connected points to overlay on the image at specified locations. This parameter takes a string of point definitions separated by the pipe character (|), or an encoded polyline using the enc: prefix within the location declaration of the path."),
						},
					}),
					"visible": arrprop("array", "(optional) specifies one or more locations that should remain visible on the map, though no markers or other indicators will be displayed.", schema{
						"type": "object",
						"properties": props{
							"location": prop("string", "(required) defines the location of the marker on the map. This parameter takes a location as either a comma-separated {latitude,longitude} pair (e.g. '40.714728,-73.998672') or a string address (e.g. 'city hall, new york, ny') identifying a unique location on the face of the earth."),
						},
					}),
					"style": arrprop("array", "(optional) defines a custom style to alter the presentation of a specific feature (roads, parks, and other features) of the map.", schema{
						"type": "object",
						"properties": props{
							"feature":    prop("string", "(required) defines the feature to style, applying the style rules to all points of the feature. Accepted values are 'all' and 'road'."),
							"element":    prop("string", "(required) defines the element to which the style is applied. An element is a feature of a map, and may be one of the following values: {all, geometry, labels}."),
							"color":      prop("string", "(optional) defines the color of the feature. This can be specified in hexadecimal (as #RRGGBB, #RRGGBBAA, #RGB, or #RGBA), or in HSLA notation. For example, '#ffcc00' will set a golden yellow color."),
							"weight":     prop("number", "(optional) defines the weight of the feature, in pixels. This setting affects polylines, where the value width is measured in pixels."),
							"visibility": prop("string", "(optional) defines the visibility of the feature. The setting can be 'on', 'off', or 'simplified'."),
						},
					}),
				},
				"required": []string{"markers", "zoom", "size"},
			},
		},
	}}, nil
}

func some[T any](t T) *T {
	return &t
}

type SchemaProperty struct {
	Type        string  `json:"type"`
	Description string  `json:"description,omitempty"`
	Items       *schema `json:"items,omitempty"`
}

func prop(tpe, description string) SchemaProperty {
	return SchemaProperty{Type: tpe, Description: description}
}

func arrprop(tpe, description string, items schema) SchemaProperty {
	return SchemaProperty{Type: tpe, Description: description, Items: &items}
}

type schema = map[string]any
type props = map[string]SchemaProperty

func callToolSuccess(bytes []byte) (res CallToolResult) {
	b64s := base64.StdEncoding.EncodeToString(bytes)
	res.Content = []Content{{Type: ContentTypeImage, Data: some(b64s), MimeType: some("image/png")}}
	return
}

func callToolError(msg string) (res CallToolResult) {
	res.IsError = some(true)
	res.Content = []Content{{Type: ContentTypeText, Text: some(msg)}}
	return
}
