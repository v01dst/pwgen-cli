package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"

)

const version = "1.0.0"

func main() {
	var (
		length   = flag.Int("l", 16, "password length (4-256)")
		count    = flag.Int("n", 1, "how many passwords to generate (1-100)")
		noUpper  = flag.Bool("no-upper", false, "exclude uppercase letters")
		noDigits = flag.Bool("no-digits", false, "exclude digits")
		symbols  = flag.Bool("s", false, "include symbols")
		noAmbig  = flag.Bool("A", false, "exclude ambiguous chars (il1Lo0O etc)")
		unique   = flag.Bool("u", false, "no repeated characters within a password")
		strength = flag.Bool("strength", false, "print entropy estimate")
		ver      = flag.Bool("version", false, "print version")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "pwgen-cli %s — secure password generator\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: pwgen-cli [flags]\n\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  pwgen-cli -l 24\n")
		fmt.Fprintf(os.Stderr, "  pwgen-cli -l 32 -s -n 5\n")
		fmt.Fprintf(os.Stderr, "  pwgen-cli -l 12 -A -u -strength\n")
	}

	flag.Parse()

	if *ver {
		fmt.Println("pwgen-cli", version)
		return
	}

	opts := Options{
		Length:        *length,
		UseUpper:      !*noUpper,
		UseDigits:     !*noDigits,
		UseSymbols:    *symbols,
		ExcludeAmbig:  *noAmbig,
		NoRepeatChars: *unique,
		Count:         *count,
	}

	passwords, err := GenerateN(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	for _, p := range passwords {
		if *strength {
			fmt.Printf("%s\t%s\n", p, entropyLabel(p))
		} else {
			fmt.Println(p)
		}
	}
}

func entropyBits(p string) float64 {
	var pool int
	sets := []string{Lower, Upper, Digits, Symbols}
	for _, set := range sets {
		if strings.ContainsAny(p, set) {
			pool += len(set)
		}
	}
	if pool == 0 {
		return 0
	}
	return float64(len(p)) * math.Log2(float64(pool))
}

func entropyLabel(p string) string {
	bits := int(entropyBits(p))
	switch {
	case bits >= 100:
		return fmt.Sprintf("%dbits-critical", bits)
	case bits >= 75:
		return fmt.Sprintf("%dbits-strong", bits)
	case bits >= 50:
		return fmt.Sprintf("%dbits-fair", bits)
	default:
		return fmt.Sprintf("%dbits-weak", bits)
	}
}
