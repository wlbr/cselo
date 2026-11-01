package postgresql

import "testing"

func TestValidatePage(t *testing.T) {
	validCases := []struct {
		limit  int
		offset int
	}{
		{limit: 1, offset: 0},
		{limit: 100, offset: 25},
	}
	for _, testCase := range validCases {
		if err := validatePage(testCase.limit, testCase.offset); err != nil {
			t.Errorf("validatePage(%d, %d) returned %v", testCase.limit, testCase.offset, err)
		}
	}

	invalidCases := []struct {
		limit  int
		offset int
	}{
		{limit: 0, offset: 0},
		{limit: 101, offset: 0},
		{limit: 1, offset: -1},
	}
	for _, testCase := range invalidCases {
		if err := validatePage(testCase.limit, testCase.offset); err == nil {
			t.Errorf("validatePage(%d, %d) returned nil error", testCase.limit, testCase.offset)
		}
	}
}
