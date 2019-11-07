package internal

import (
	"strings"

	"github.com/eyesore/maroto/pkg/consts"
	"github.com/eyesore/maroto/pkg/props"
	"github.com/jung-kurt/gofpdf"
)

// Text is the abstraction which deals of how to add text inside PDF
type Text interface {
	Add(text string, fontFamily props.Text, marginTop float64, actualCol float64, qtdCols float64)
	VariableAdd(text string, fontFamily props.Text, marginTop float64, actualCol float64, qtdCols float64, widthMultiplier int)
}

type text struct {
	pdf  gofpdf.Pdf
	math Math
	font Font
}

// NewText create a Text
func NewText(pdf gofpdf.Pdf, math Math, font Font) *text {
	return &text{
		pdf,
		math,
		font,
	}
}

// Add a text inside a cell.
func (s *text) Add(text string, textProp props.Text, marginTop float64, actualCol float64, qtdCols float64) {
	actualWidthPerCol := s.math.GetWidthPerCol(qtdCols)

	translator := s.pdf.UnicodeTranslatorFromDescriptor("")
	s.font.SetFont(textProp.Family, textProp.Style, textProp.Size)

	textTranslated := translator(text)
	stringWidth := s.pdf.GetStringWidth(textTranslated)
	words := strings.Split(textTranslated, " ")

	if stringWidth < actualWidthPerCol || textProp.Extrapolate || len(words) == 1 {
		s.addLine(textProp, actualCol, actualWidthPerCol, marginTop, stringWidth, textTranslated)
	} else {
		currentlySize := 0.0
		actualLine := 0
		lines := []string{}
		lines = append(lines, "")

		for _, word := range words {
			if s.pdf.GetStringWidth(word+" ")+currentlySize < actualWidthPerCol {
				lines[actualLine] = lines[actualLine] + word + " "
				currentlySize += s.pdf.GetStringWidth(word + " ")
			} else {
				lines = append(lines, "")
				actualLine++
				lines[actualLine] = lines[actualLine] + word + " "
				currentlySize = s.pdf.GetStringWidth(word + " ")
			}
		}

		for index, line := range lines {
			lineWidth := s.pdf.GetStringWidth(line)
			s.addLine(textProp, actualCol, actualWidthPerCol, marginTop+float64(index)*textProp.Size/2.0, lineWidth, line)
		}
	}
}

// Add a text inside a cell for variable columns.
func (s *text) VariableAdd(text string, textProp props.Text, marginTop float64, actualCol float64, qtdCols float64, widthMultiplier int) {
	actualWidthPerCol := s.math.GetWidthPerCol(qtdCols)
	actualWidth := actualWidthPerCol * float64(widthMultiplier)

	translator := s.pdf.UnicodeTranslatorFromDescriptor("")
	s.font.SetFont(textProp.Family, textProp.Style, textProp.Size)

	textTranslated := translator(text)
	stringWidth := s.pdf.GetStringWidth(textTranslated)
	words := strings.Split(textTranslated, " ")

	if stringWidth < actualWidth || textProp.Extrapolate || len(words) == 1 {
		s.addVariableLine(textProp, actualCol, actualWidthPerCol, actualWidth, marginTop, stringWidth, textTranslated)
	} else {
		currentlySize := 0.0
		actualLine := 0
		lines := []string{}
		lines = append(lines, "")

		for _, word := range words {
			if s.pdf.GetStringWidth(word+" ")+currentlySize < actualWidth {
				lines[actualLine] = lines[actualLine] + word + " "
				currentlySize += s.pdf.GetStringWidth(word + " ")
			} else {
				lines = append(lines, "")
				actualLine++
				lines[actualLine] = lines[actualLine] + word + " "
				currentlySize = s.pdf.GetStringWidth(word + " ")
			}
		}

		for index, line := range lines {
			lineWidth := s.pdf.GetStringWidth(line)
			s.addVariableLine(textProp, actualCol, actualWidthPerCol, actualWidth, marginTop+float64(index)*textProp.Size/2.0, lineWidth, line)
		}
	}
}

func (s *text) addLine(textProp props.Text, actualCol, actualWidthPerCol, marginTop, stringWidth float64, textTranslated string) {
	left, top, _, _ := s.pdf.GetMargins()

	if textProp.Align == consts.Left {
		s.pdf.Text(actualCol*actualWidthPerCol+left, marginTop+top, textTranslated)
		return
	}

	var modifier float64 = 2

	if textProp.Align == consts.Right {
		modifier = 1
	}

	dx := (actualWidthPerCol - stringWidth) / modifier

	s.pdf.Text(dx+actualCol*actualWidthPerCol+left, marginTop+top, textTranslated)
}

func (s *text) addVariableLine(textProp props.Text, actualCol, actualWidthPerCol, actualWidth, marginTop, stringWidth float64, textTranslated string) {
	left, top, _, _ := s.pdf.GetMargins()

	if textProp.Align == consts.Left {
		s.pdf.Text(actualCol*actualWidthPerCol+left, marginTop+top, textTranslated)
		return
	}

	var modifier float64 = 2

	if textProp.Align == consts.Right {
		modifier = 1
	}

	dx := (actualWidth - stringWidth) / modifier

	s.pdf.Text(dx+actualCol*actualWidthPerCol+left, marginTop+top, textTranslated)
}
