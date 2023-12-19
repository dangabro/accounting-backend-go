package style

import "github.com/xuri/excelize/v2"

const (
	Title                   = iota
	CenterBold              = iota
	Bold                    = iota
	Right                   = iota
	RightBold               = iota
	VerticalCenter          = iota
	VerticalBottom          = iota
	NumberFormat            = iota
	NumberFormatBold        = iota
	Center                  = iota
	VerticalCenterRight     = iota
	VerticalCenterRightBold = iota
	VerticalBottomBold      = iota
)

type style1 struct {
	styles map[int]int
	f      *excelize.File
}

type Style interface {
	GetStyle(int) int
}

func (s *style1) GetStyle(val int) int {
	res, ok := s.styles[val]
	if !ok {
		panic("cannot find style")
	}

	return res
}

func NewStyle(f *excelize.File) Style {
	var res style1
	res.f = f

	fillStyle(&res)

	return &res
}

func fillStyle(s *style1) {
	st := make(map[int]int)
	s.styles = st

	f := s.f

	st[Title] = titleStyle(f)
	st[CenterBold] = centerBold(f)
	st[Bold] = bold(f)
	st[Right] = right(f)
	st[RightBold] = rightBold(f)
	st[VerticalCenter] = verticalCenter(f)
	st[NumberFormat] = numberFormat(f)
	st[NumberFormatBold] = numberFormatBold(f)
	st[VerticalBottom] = verticalBottom(f)
	st[Center] = center(f)
	st[VerticalCenterRight] = verticalCenterRight(f)
	st[VerticalCenterRightBold] = verticalCenterRightBold(f)
	st[VerticalBottomBold] = verticalBottomBold(f)
}

func titleStyle(f *excelize.File) int {
	style := excelize.Style{
		Font:       &excelize.Font{Bold: true, Size: 20, Strike: false, Color: "", ColorIndexed: 0, ColorTheme: nil, ColorTint: 0, VertAlign: ""},
		Alignment:  &excelize.Alignment{Horizontal: "center", Indent: 0, JustifyLastLine: false, ReadingOrder: 0, RelativeIndent: 0, ShrinkToFit: false, TextRotation: 0, Vertical: "center", WrapText: false},
		Protection: nil, NumFmt: 0, DecimalPlaces: 0, CustomNumFmt: nil, Lang: "", NegRed: false,
	}
	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func centerBold(f *excelize.File) int {
	style := excelize.Style{
		Font:       &excelize.Font{Bold: true, Strike: false, Color: "", ColorIndexed: 0, ColorTheme: nil, ColorTint: 0, VertAlign: ""},
		Alignment:  &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Protection: nil, NumFmt: 0, DecimalPlaces: 0, CustomNumFmt: nil, Lang: "", NegRed: false,
	}
	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func bold(f *excelize.File) int {
	style := excelize.Style{
		Font: &excelize.Font{Bold: true},
	}
	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func right(f *excelize.File) int {
	style := excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right"},
	}
	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func rightBold(f *excelize.File) int {
	style := excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Font:      &excelize.Font{Bold: true},
	}
	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func verticalCenter(f *excelize.File) int {
	style := excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center"},
	}
	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func verticalCenterRight(f *excelize.File) int {
	style := excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "right"},
	}
	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func verticalCenterRightBold(f *excelize.File) int {
	style := excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "right"},
		Font:      &excelize.Font{Bold: true},
	}
	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func numberFormat(f *excelize.File) int {
	style := excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right"},
		NumFmt:    4,
	}

	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func numberFormatBold(f *excelize.File) int {
	style := excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Font:      &excelize.Font{Bold: true},
		NumFmt:    4,
	}

	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func verticalBottom(f *excelize.File) int {
	style := excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "bottom"},
	}

	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func verticalBottomBold(f *excelize.File) int {
	style := excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "bottom"},
		Font:      &excelize.Font{Bold: true},
	}

	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}

func center(f *excelize.File) int {
	style := excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
	}

	index, err := f.NewStyle(&style)

	if err != nil {
		panic(err.Error())
	}

	return index
}
