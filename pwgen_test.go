package main

import (
	"strings"
	"testing"
)

func TestGenerateDefault(t *testing.T) {
	p, err := Generate(Options{Length: 16, Count: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p) != 16 {
		t.Fatalf("want length 16, got %d", len(p))
	}
	if !strings.ContainsAny(p, Lower) {
		t.Fatalf("password %q has no lowercase", p)
	}
}

func TestGuaranteesEnabledSets(t *testing.T) {
	for i := 0; i < 50; i++ {
		p, err := Generate(Options{Length: 24, UseUpper: true, UseDigits: true, UseSymbols: true})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.ContainsAny(p, Upper) {
			t.Fatalf("missing upper in %q", p)
		}
		if !strings.ContainsAny(p, Digits) {
			t.Fatalf("missing digit in %q", p)
		}
		if !strings.ContainsAny(p, Symbols) {
			t.Fatalf("missing symbol in %q", p)
		}
	}
}

func TestExcludeAmbiguous(t *testing.T) {
	for i := 0; i < 50; i++ {
		p, err := Generate(Options{Length: 32, UseUpper: true, UseDigits: true, ExcludeAmbig: true})
		if err != nil {
			t.Fatal(err)
		}
		if strings.ContainsAny(p, Ambiguous) {
			t.Fatalf("ambiguous char in %q", p)
		}
	}
}

func TestNoRepeatChars(t *testing.T) {
	p, err := Generate(Options{Length: 20, NoRepeatChars: true})
	if err != nil {
		t.Fatal(err)
	}
	seen := make(map[rune]bool)
	for _, r := range p {
		if seen[r] {
			t.Fatalf("repeated char %q in %q", r, p)
		}
		seen[r] = true
	}
}

func TestValidation(t *testing.T) {
	if _, err := Generate(Options{Length: 3}); err != ErrBadLength {
		t.Fatalf("want ErrBadLength, got %v", err)
	}
	if _, err := Generate(Options{Length: 257}); err != ErrBadLength {
		t.Fatalf("want ErrBadLength, got %v", err)
	}
	if _, err := GenerateN(Options{Length: 16, Count: 0}); err != ErrBadCount {
		t.Fatalf("want ErrBadCount, got %v", err)
	}
	if _, err := GenerateN(Options{Length: 16, Count: 101}); err != ErrBadCount {
		t.Fatalf("want ErrBadCount, got %v", err)
	}
}

func TestGenerateN(t *testing.T) {
	ps, err := GenerateN(Options{Length: 16, Count: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 5 {
		t.Fatalf("want 5, got %d", len(ps))
	}
	unique := make(map[string]bool)
	for _, p := range ps {
		unique[p] = true
	}
	if len(unique) < 2 {
		t.Fatal("expected distinct passwords in a 5-batch")
	}
}

func TestUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		p, err := Generate(Options{Length: 16})
		if err != nil {
			t.Fatal(err)
		}
		if seen[p] {
			t.Fatalf("duplicate password generated: %q", p)
		}
		seen[p] = true
	}
}

func TestGeneratePassphrase(t *testing.T) {
	// deterministic fake rng
	i := 0
	randInt := func(max int) int {
		i++
		return (i * 7) % max
	}
	p, err := GeneratePassphrase(randInt, 4, "-")
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(p, "-")
	if len(parts) != 4 {
		t.Fatalf("want 4 words, got %d: %q", len(parts), p)
	}
	found := false
	for _, w := range parts {
		w = strings.TrimRight(w, "0123456789")
		for _, valid := range Wordlist {
			if w == valid {
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatalf("no wordlist words found in %q", p)
	}
}

func TestGeneratePassphraseValidation(t *testing.T) {
	randInt := func(max int) int { return 0 }
	if _, err := GeneratePassphrase(randInt, 2, "-"); err == nil {
		t.Fatal("too few words should error")
	}
	if _, err := GeneratePassphrase(randInt, 13, "-"); err == nil {
		t.Fatal("too many words should error")
	}
}

func TestGeneratePassphraseHasDigit(t *testing.T) {
	randInt := func(max int) int { return 3 }
	p, _ := GeneratePassphrase(randInt, 5, "-")
	hasDigit := false
	for _, r := range p {
		if r >= '0' && r <= '9' {
			hasDigit = true
		}
	}
	if !hasDigit {
		t.Fatalf("passphrase %q should contain a digit", p)
	}
}

func TestGeneratePattern(t *testing.T) {
	i := 0
	randInt := func(max int) int {
		i++
		return (i * 11) % max
	}
	p, err := GeneratePattern(randInt, "Aaaa-9999")
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 9 {
		t.Fatalf("len = %d, want 9: %q", len(p), p)
	}
	if p[4] != '-' {
		t.Fatalf("literal dash missing: %q", p)
	}
	if !strings.ContainsAny(p[:4], "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		t.Fatalf("no upper: %q", p)
	}
}

func TestGeneratePatternAllLiteral(t *testing.T) {
	randInt := func(max int) int { return 0 }
	p, err := GeneratePattern(randInt, "hello-world")
	if err != nil {
		t.Fatal(err)
	}
	if p != "hello-world" {
		t.Fatalf("literal passthrough failed: %q", p)
	}
}
