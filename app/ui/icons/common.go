package icons

import "github.com/maddalax/htmgo/framework/h"

func TrashIcon() *h.Element {
	return h.Svg(
		h.Class("h-full w-full"),
		h.Attribute("viewBox", "0 0 24 24"),
		h.Attribute("fill", "none"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Path(
			h.Attribute("d", "M10 11V17"),
			h.Attribute("stroke", "#000000"),
			h.Attribute("stroke-width", "2"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Path(
			h.Attribute("d", "M14 11V17"),
			h.Attribute("stroke", "#000000"),
			h.Attribute("stroke-width", "2"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Path(
			h.Attribute("d", "M4 7H20"),
			h.Attribute("stroke", "#000000"),
			h.Attribute("stroke-width", "2"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Path(
			h.Attribute("d", "M6 7H12H18V18C18 19.6569 16.6569 21 15 21H9C7.34315 21 6 19.6569 6 18V7Z"),
			h.Attribute("stroke", "#000000"),
			h.Attribute("stroke-width", "2"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Path(
			h.Attribute("d", "M9 5C9 3.89543 9.89543 3 11 3H13C14.1046 3 15 3.89543 15 5V7H9V5Z"),
			h.Attribute("stroke", "#000000"),
			h.Attribute("stroke-width", "2"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
	)
}

func Question() *h.Element {
	return h.Svg(
		h.Class("h-full w-full"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Attribute("viewBox", "0 0 24 24"),
		h.Attribute("fill", "none"),
		h.Attribute("stroke", "#000000"),
		h.Attribute("stroke-width", "2"),
		h.Attribute("stroke-linecap", "round"),
		h.Attribute("stroke-linejoin", "round"),
		h.Tag(
			"circle",
			h.Attribute("cx", "12"),
			h.Attribute("cy", "12"),
			h.Attribute("r", "10"),
		),
		h.Path(
			h.Attribute("d", "M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"),
		),
		h.Tag(
			"line",
			h.Attribute("x1", "12"),
			h.Attribute("y1", "17"),
			h.Attribute("x2", "12.01"),
			h.Attribute("y2", "17"),
		),
	)
}
