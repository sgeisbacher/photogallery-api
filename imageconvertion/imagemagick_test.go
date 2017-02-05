package imageconvertion

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestBuildDimension(t *testing.T) {
	RegisterTestingT(t)

	testDataTable := []struct {
		mediaFormat MediaFormat
		expectedStr string
		expectedErr bool
	}{
		{MediaFormat{}, "", true},
		{MediaFormat{120, 0, "abc"}, "120x", false},
		{MediaFormat{0, 120, "abc"}, "x120", false},
		{MediaFormat{240, 120, "abc"}, "240x120", false},
	}

	for _, testData := range testDataTable {
		dimensionStr, err := BuildDimension(testData.mediaFormat)
		Expect(dimensionStr).To(Equal(testData.expectedStr))
		Expect(err != nil).To(Equal(testData.expectedErr))
	}
}
