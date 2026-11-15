package analytics

import (
	"testing"
	"time"
)

func TestParseDateRange(t *testing.T) {
	tests := []struct {
		name, start, end, wantStart, wantEnd string
		wantError                            bool
	}{
		{name: "all time"},
		{name: "start only", start: "2026-10-01", wantStart: "2026-10-01"},
		{name: "end only", end: "2026-10-03", wantEnd: "2026-10-04"},
		{name: "inclusive same day", start: "2026-10-03", end: "2026-10-03", wantStart: "2026-10-03", wantEnd: "2026-10-04"},
		{name: "first representable date", start: "0001-01-01", wantStart: "0001-01-01"},
		{name: "bad start", start: "not-a-date", wantError: true},
		{name: "bad end", end: "2026-02-30", wantError: true},
		{name: "reversed", start: "2026-10-03", end: "2026-10-01", wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := ParseDateRange(tt.start, tt.end)
			if (err != nil) != tt.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tt.wantError)
			}
			if tt.wantError {
				return
			}
			for _, bound := range []struct {
				name string
				got  *time.Time
				want string
			}{{"start", start, tt.wantStart}, {"end", end, tt.wantEnd}} {
				if bound.want == "" {
					if bound.got != nil {
						t.Errorf("%s = %v, want nil", bound.name, bound.got)
					}
					continue
				}
				if bound.got == nil || bound.got.Format(time.DateOnly) != bound.want || bound.got.Location() != time.UTC {
					t.Errorf("%s = %v, want %s UTC", bound.name, bound.got, bound.want)
				}
			}
		})
	}
}
