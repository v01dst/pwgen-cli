package main

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strconv"
	"strings"
)

const (
	Lower     = "abcdefghijklmnopqrstuvwxyz"
	Upper     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits    = "0123456789"
	Symbols   = "!@#$%^&*()-_=+[]{}<>?~"
	Ambiguous = "il1Lo0O`'|"
)

// Options controls password generation.
type Options struct {
	Length        int
	UseUpper      bool
	UseDigits     bool
	UseSymbols    bool
	ExcludeAmbig  bool
	NoRepeatChars bool
	Count         int
}

// ErrNoCharsets is returned when every character set is disabled.
var ErrNoCharsets = errors.New("pwgen: at least one character set must be enabled")

// ErrBadLength is returned for lengths outside [4, 256].
var ErrBadLength = errors.New("pwgen: length must be between 4 and 256")

// ErrBadCount is returned for counts outside [1, 100].
var ErrBadCount = errors.New("pwgen: count must be between 1 and 100")

// Generate produces one password satisfying the options, guaranteeing at
// least one character from every enabled set.
func Generate(opts Options) (string, error) {
	if opts.Length < 4 || opts.Length > 256 {
		return "", ErrBadLength
	}

	sets := []string{Lower}
	if opts.UseUpper {
		sets = append(sets, Upper)
	}
	if opts.UseDigits {
		sets = append(sets, Digits)
	}
	if opts.UseSymbols {
		sets = append(sets, Symbols)
	}

	pool := strings.Join(sets, "")
	if opts.ExcludeAmbig {
		pool = stripChars(pool, Ambiguous)
		for i, s := range sets {
			sets[i] = stripChars(s, Ambiguous)
		}
	}

	// A set can vanish entirely after ambiguity filtering; drop empty sets.
	filtered := sets[:0]
	for _, s := range sets {
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	sets = filtered

	if pool == "" || len(sets) == 0 {
		return "", ErrNoCharsets
	}

	used := make(map[rune]bool)
	var out strings.Builder
	out.Grow(opts.Length)

	// Guarantee one char from each enabled set first.
	for _, set := range sets {
		if len(out.String()) >= opts.Length {
			break
		}
		c, err := randChar(set, used, opts.NoRepeatChars)
		if err != nil {
			return "", err
		}
		out.WriteRune(c)
	}

	for out.Len() < opts.Length {
		c, err := randChar(pool, used, opts.NoRepeatChars)
		if err != nil {
			return "", err
		}
		out.WriteRune(c)
	}

	// Shuffle so guaranteed chars are not always at the front.
	runes := []rune(out.String())
	shuffle(runes)
	return string(runes), nil
}

// GenerateN produces n passwords.
func GenerateN(opts Options) ([]string, error) {
	if opts.Count < 1 || opts.Count > 100 {
		return nil, ErrBadCount
	}
	passwords := make([]string, 0, opts.Count)
	for i := 0; i < opts.Count; i++ {
		p, err := Generate(opts)
		if err != nil {
			return nil, err
		}
		passwords = append(passwords, p)
	}
	return passwords, nil
}

func stripChars(s, remove string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(remove, r) {
			return -1
		}
		return r
	}, s)
}

func randChar(pool string, used map[rune]bool, noRepeat bool) (rune, error) {
	runes := []rune(pool)
	// With noRepeat, retry until an unused char is drawn (bounded to avoid
	// pathological loops when the pool is nearly exhausted).
	for attempt := 0; attempt < 128; attempt++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(runes))))
		if err != nil {
			return 0, err
		}
		c := runes[n.Int64()]
		if !noRepeat || !used[c] {
			used[c] = true
			return c, nil
		}
	}
	return 0, errors.New("pwgen: pool too small for no-repeat generation")
}

// Wordlist for passphrase mode — short, unambiguous, memorable words.
var Wordlist = []string{
	"apple", "arrow", "basil", "beach", "birch", "bison", "blade", "brave",
	"brisk", "camel", "candle", "cedar", "chalk", "cherry", "cinder", "citrus",
	"clover", "comet", "copper", "coral", "cosmos", "cotton", "crane", "crisp",
	"crystal", "daisy", "dawn", "delta", "denim", "dune", "ember", "fable",
	"fern", "fiber", "flint", "flora", "frost", "garnet", "ginger", "granite",
	"harbor", "hazel", "henna", "ivory", "jasper", "jungle", "kelp", "kite",
	"lagoon", "lantern", "lava", "lemon", "lilac", "linen", "lotus", "lunar",
	"maple", "marble", "meadow", "mint", "morning", "moss", "nectar", "nimbus",
	"oasis", "ocean", "onyx", "opal", "orchid", "otter", "palm", "pebble",
	"pepper", "petal", "pilot", "pine", "plume", "prairie", "quartz", "quill",
	"radish", "rapid", "raven", "reef", "ripple", "river", "rocket", "saffron",
	"sage", "salmon", "sandal", "sapphire", "satin", "savanna", "shadow", "shore",
	"signal", "silver", "solar", "spruce", "star", "storm", "summer", "sunset",
	"tangerine", "teak", "thistle", "thunder", "tide", "timber", "topaz", "tulip",
	"umber", "valley", "velvet", "vertex", "violet", "walnut", "willow", "zenith",
}

// GeneratePassphrase builds a passphrase of n words joined by a separator,
// with a random digit appended to one word for strength.
func GeneratePassphrase(randInt func(max int) int, words int, separator string) (string, error) {
	if words < 3 || words > 12 {
		return "", errors.New("pwgen: word count must be between 3 and 12")
	}
	parts := make([]string, words)
	for i := 0; i < words; i++ {
		parts[i] = Wordlist[randInt(len(Wordlist))]
	}
	digit := randInt(10)
	parts[randInt(words)] += strconv.Itoa(digit)
	return strings.Join(parts, separator), nil
}

func shuffle(runes []rune) {
	for i := len(runes) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue // keep current order on failure; still random-ish
		}
		j := n.Int64()
		runes[i], runes[j] = runes[j], runes[i]
	}
}
