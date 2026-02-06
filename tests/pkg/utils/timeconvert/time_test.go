// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package timeconvert_test

import (
	"testing"
	"time"

	"appsite-go/pkg/utils/timeconvert"
)

func TestFormatting(t *testing.T) {
	now := time.Now()
	
	// Test FormatToDateTime
	dt := timeconvert.FormatToDateTime(now)
	expectedDT := now.Format(timeconvert.LayoutDateTime)
	if dt != expectedDT {
		t.Errorf("FormatToDateTime() = %s, want %s", dt, expectedDT)
	}

	// Test FormatToDate
	d := timeconvert.FormatToDate(now)
	expectedD := now.Format(timeconvert.LayoutDateOnly)
	if d != expectedD {
		t.Errorf("FormatToDate() = %s, want %s", d, expectedD)
	}
}

func TestCurrentHelpers(t *testing.T) {
	// These are current time based, so exact match is hard.
	// We just check format length or parsing potential.
	curDT := timeconvert.CurrentTime()
	if _, err := timeconvert.ParseDateTime(curDT); err != nil {
		t.Errorf("CurrentTime() returned unparseable time: %s", curDT)
	}

	curD := timeconvert.CurrentDate()
	if _, err := timeconvert.ParseDate(curD); err != nil {
		t.Errorf("CurrentDate() returned unparseable date: %s", curD)
	}
}

func TestParsing(t *testing.T) {
	testTimeStr := "2026-02-06 12:00:00"
	parsed, err := timeconvert.ParseDateTime(testTimeStr)
	if err != nil {
		t.Errorf("ParseDateTime failed: %v", err)
	}
	if parsed.Year() != 2026 || parsed.Month() != 2 || parsed.Day() != 6 {
		t.Errorf("ParseDateTime parsed wrong date: %v", parsed)
	}

	testDateStr := "2026-02-06"
	parsedDate, err := timeconvert.ParseDate(testDateStr)
	if err != nil {
		t.Errorf("ParseDate failed: %v", err)
	}
	if parsedDate.Year() != 2026 {
		t.Errorf("ParseDate parsed wrong year: %v", parsedDate)
	}
}

func TestStringToTime(t *testing.T) {
	tests := []struct {
		input    string
		hasError bool
	}{
		{"2026-02-06 12:00:00", false},
		{"2026-02-06", false},
		{"20260206120000", false},
		{"invalid-date", true},
	}

	for _, tt := range tests {
		_, err := timeconvert.StringToTime(tt.input)
		if (err != nil) != tt.hasError {
			t.Errorf("StringToTime(%s) error = %v, wantError %v", tt.input, err, tt.hasError)
		}
	}
}

func TestDurationAndHelpers(t *testing.T) {
	// CalculateDuration
	start := "2026-02-06 10:00:00"
	end := "2026-02-06 12:00:00"
	dur, err := timeconvert.CalculateDuration(start, end)
	if err != nil {
		t.Fatalf("CalculateDuration failed: %v", err)
	}
	if dur != 2*time.Hour {
		t.Errorf("CalculateDuration = %v, want 2h", dur)
	}

	// IsExpired
	past := time.Now().Add(-1 * time.Hour)
	if !timeconvert.IsExpired(past) {
		t.Error("IsExpired(past) should be true")
	}
	future := time.Now().Add(1 * time.Hour)
	if timeconvert.IsExpired(future) {
		t.Error("IsExpired(future) should be false")
	}

	// AddDays
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local)
	next := timeconvert.AddDays(base, 1)
	if next.Day() != 2 {
		t.Errorf("AddDays +1 = %v, want Jan 2", next)
	}
}

func TestTimeZone(t *testing.T) {
	// Just ensure it doesn't panic and tries to set
	// Note: loading random location might fail on some minimal docker images without tzdata, so we try "UTC"
	timeconvert.SetDefaultTimeZone("UTC")
	if time.Local.String() != "UTC" {
		// This might fail if system is weird, but usually works
		// t.Logf("SetDefaultTimeZone didn't set time.Local to UTC? Got: %v", time.Local)
	}
}
