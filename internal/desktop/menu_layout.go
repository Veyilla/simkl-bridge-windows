package desktop

type MenuLayout struct {
	Rows          []Rect
	Separators    []int
	Footer        Rect
	Width, Height int
}

func HomeLayout() MenuLayout {
	return MenuLayout{Width: 368, Height: 552, Rows: []Rect{{16, 142, 336, 44}, {16, 186, 336, 44}, {16, 250, 336, 44}, {16, 294, 336, 44}, {16, 338, 336, 44}, {16, 382, 336, 44}, {16, 446, 336, 44}}, Separators: []int{130, 240, 436}, Footer: Rect{24, 512, 320, 28}}
}
