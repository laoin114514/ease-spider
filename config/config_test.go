package config

import (
	"testing"
	"time"
)

func TestParseInterval(t *testing.T) {
	cases := []struct {
		raw  string
		want time.Duration
		ok   bool
	}{
		{"30s", 30 * time.Second, true},
		{"5m", 5 * time.Minute, true},
		{"2h", 2 * time.Hour, true},
		{"1d", 24 * time.Hour, true},
		{" 3m ", 3 * time.Minute, true},
		{"5", 5 * time.Minute, true},  // 不带单位按分钟
		{"5x", 5 * time.Minute, true}, // 单位无法识别也按分钟（历史行为）
		{"", 0, false},
		{"abc", 0, false},
		{"0m", 0, false},
		{"-5m", 0, false},
	}
	for _, tc := range cases {
		got, ok := ParseInterval(tc.raw)
		if ok != tc.ok || (ok && got != tc.want) {
			t.Errorf("ParseInterval(%q) = (%v, %v), want (%v, %v)", tc.raw, got, ok, tc.want, tc.ok)
		}
	}
}

func TestIntervalOr(t *testing.T) {
	const def = 7 * time.Minute
	if got := IntervalOr("", def); got != def {
		t.Errorf("空值应返回默认值 %v, got %v", def, got)
	}
	if got := IntervalOr("abc", def); got != def {
		t.Errorf("非法值应返回默认值 %v, got %v", def, got)
	}
	if got := IntervalOr("90s", def); got != 90*time.Second {
		t.Errorf("IntervalOr(90s) = %v, want 90s", got)
	}
}

func TestTimerFrequencyEntriesCoversAllKeys(t *testing.T) {
	want := []string{
		"cf_official_contests", "cf_official_problems", "cf_records",
		"cf_team_contests", "cf_team_contests_problems",
		"dingding", "luogu_records", "niuke_records",
	}
	got := TimerFrequency().Entries()
	if len(got) != len(want) {
		t.Fatalf("Entries() 返回 %d 项, want %d", len(got), len(want))
	}
	for i, e := range got {
		if e.Key != want[i] {
			t.Errorf("Entries()[%d].Key = %q, want %q", i, e.Key, want[i])
		}
	}
}
