package handler

import (
	"os"
	"regexp"
	"testing"
)

// These tests exercise DefaultTagPattern and resolveTagPreservation.
//
// init() requires the REGISTRY_* / access-key environment variables and
// panics without them, so run the tests with dummy values:
//
//	REGISTRY_ACCESS_KEY=test REGISTRY_PROJECT_ID=test REGISTRY_REGION=fr-par \
//	REGISTRY_NAMESPACE=test REGISTRY_SECRET_KEY=test go test ./...

func TestDefaultTagPatternMatchesSemverAndLatest(t *testing.T) {
	re := regexp.MustCompile(DefaultTagPattern)

	preserved := []string{
		"1.2.3",
		"1.2.3-rc.1",
		"1.0.0-alpha+001",
		"0.1.0",
		"10.20.30",
		"1.0.0+20130313144700",
		"latest",
		"latest-arm64",
	}
	for _, tag := range preserved {
		if !re.MatchString(tag) {
			t.Errorf("DefaultTagPattern should preserve %q, it does not match", tag)
		}
	}

	deleted := []string{
		"feature-x",
		"v1.2.3",
		"1.2",
		"1.2.3.4",
		"sha256-abc123",
	}
	for _, tag := range deleted {
		if re.MatchString(tag) {
			t.Errorf("DefaultTagPattern should not preserve %q, it matches", tag)
		}
	}
}

func TestResolveTagPreservation(t *testing.T) {
	// An invalid pattern falls back to the default and forces dry-run mode.
	re, forceDryRun, err := resolveTagPreservation("^(broken(")
	if err == nil {
		t.Fatal("expected a compile error for an invalid pattern, got nil")
	}
	if !forceDryRun {
		t.Fatal("expected dry-run to be forced when the configured pattern is invalid")
	}
	for _, tag := range []string{"1.2.3", "1.2.3-rc.1", "latest"} {
		if !re.MatchString(tag) {
			t.Errorf("fallback pattern should preserve %q, it does not match", tag)
		}
	}
	if re.MatchString("feature-x") {
		t.Error("fallback pattern should not preserve feature-x")
	}

	// An empty value (no patterns configured) keeps the original semantics:
	// the empty pattern matches every tag, so everything is preserved, no
	// error, no forced dry-run.
	re, forceDryRun, err = resolveTagPreservation("")
	if err != nil {
		t.Fatalf("expected no error for an empty pattern, got %v", err)
	}
	if forceDryRun {
		t.Fatal("expected dry-run not to be forced for an empty pattern")
	}
	if !re.MatchString("1.2.3") || !re.MatchString("anything") {
		t.Error("empty pattern should preserve every tag")
	}

	// A valid user pattern is used as-is.
	re, forceDryRun, err = resolveTagPreservation("^v\\d+$|^release-\\d+$")
	if err != nil {
		t.Fatalf("expected no error for a valid pattern, got %v", err)
	}
	if forceDryRun {
		t.Fatal("expected dry-run not to be forced for a valid pattern")
	}
	if !re.MatchString("v12") || !re.MatchString("release-7") {
		t.Error("user pattern should match its own tags")
	}
	if re.MatchString("1.2.3") {
		t.Error("user pattern should not match tags it does not name")
	}
}

// TestResolveTagPreservationIsolated guards the init() environment: the test
// binary only reaches these tests because the REGISTRY_* variables were set
// (see the package comment).
func TestResolveTagPreservationIsolated(t *testing.T) {
	if os.Getenv("REGISTRY_ACCESS_KEY") == "" {
		t.Skip("run with REGISTRY_ACCESS_KEY set (see package comment)")
	}
}
