package main

import (
	"testing"

	"github.com/muesli/termenv"
)

func TestApplyColorFallback(t *testing.T) {
	tests := []struct {
		name               string
		colorMode          bool
		colorAsciiMode     bool
		profile            termenv.Profile
		wantColorMode      bool
		wantColorAsciiMode bool
	}{
		{"default mode without color", true, false, termenv.Ascii, false, false},
		{"color ASCII without color", false, true, termenv.Ascii, false, false},
		{"color mode with ANSI", true, false, termenv.ANSI, true, false},
		{"color ASCII with truecolor", false, true, termenv.TrueColor, false, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			colorMode, colorAsciiMode := applyColorFallback(test.colorMode, test.colorAsciiMode, test.profile)
			if colorMode != test.wantColorMode || colorAsciiMode != test.wantColorAsciiMode {
				t.Fatalf("applyColorFallback() = (%t, %t), want (%t, %t)", colorMode, colorAsciiMode, test.wantColorMode, test.wantColorAsciiMode)
			}
		})
	}
}

func TestFitDimensions(t *testing.T) {
	tests := []struct {
		name                  string
		origW, origH          int
		maxWidth, maxHeight   int
		wantWidth, wantHeight int
	}{
		{"width limit", 1600, 900, 80, 0, 80, 45},
		{"height limit", 1600, 900, 80, 20, 36, 20},
		{"portrait", 600, 1200, 80, 24, 12, 24},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			width, height := fitDimensions(test.origW, test.origH, test.maxWidth, test.maxHeight)
			if width != test.wantWidth || height != test.wantHeight {
				t.Fatalf("fitDimensions() = (%d, %d), want (%d, %d)", width, height, test.wantWidth, test.wantHeight)
			}
		})
	}
}
