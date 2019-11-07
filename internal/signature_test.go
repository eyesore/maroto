package internal_test

import (
	"fmt"
	"testing"

	"github.com/eyesore/maroto/internal"
	"github.com/eyesore/maroto/internal/mocks"
	"github.com/eyesore/maroto/pkg/props"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewSignature(t *testing.T) {
	signature := internal.NewSignature(&mocks.Pdf{}, &mocks.Math{}, &mocks.Text{})

	assert.NotNil(t, signature)
	assert.Equal(t, fmt.Sprintf("%T", signature), "*internal.signature")
}

func TestSignature_AddSpaceFor(t *testing.T) {
	// Arrange
	pdf := &mocks.Pdf{}
	pdf.On("GetMargins").Return(10.0, 10.0, 10.0, 10.0)
	pdf.On("Line", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	math := &mocks.Math{}
	math.On("GetWidthPerCol", mock.Anything).Return(50.0)

	text := &mocks.Text{}
	text.On("Add", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	signature := internal.NewSignature(pdf, math, text)

	// Act
	signature.AddSpaceFor("label", props.Text{Size: 10.0}, 5, 5, 2)

	// Assert
	pdf.AssertNumberOfCalls(t, "Line", 1)
	pdf.AssertCalled(t, "Line", 114.0, 10.0, 156.0, 10.0)
	text.AssertNumberOfCalls(t, "Add", 1)
	text.AssertCalled(t, "Add", "label", props.Text{Size: 10.0}, 5.0, 2.0, 5.0)
}
