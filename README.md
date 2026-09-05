<div align="center">

# 🔑 pwgen-cli

**Cryptographically secure password generator — tiny, fast, no dependencies.**

[![CI](https://github.com/v01dst/pwgen-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/v01dst/pwgen-cli/actions/workflows/ci.yml)
![License](https://img.shields.io/badge/license-MIT-8A2BE2)
![Go](https://img.shields.io/badge/go-1.27-00ADD8?logo=go&logoColor=white)
![Tests](https://img.shields.io/badge/tests-passing-brightgreen)
![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-blue)

`crypto/rand` · `zero deps` · `entropy meter` · `single binary`
</div>

---

## ✨ Features

- **🔐 Crypto-grade randomness** — backed by Go's `crypto/rand`, never `math/rand`
- **🎛️ Full control** — length, uppercase, digits, symbols, ambiguous chars, repeat policy
- **✅ Set guarantees** — every enabled character set appears at least once, then shuffled
- **📏 Entropy meter** — instant strength readout in bits (`-strength`)
- **🪶 Single binary** — pure stdlib, compiles everywhere, runs anywhere
- **🚫 Anti-footgun** — ambiguity filter (`-A`) for humans, no-repeat mode (`-u`) for_PINs

## 🚀 Quick Start

```bash
git clone https://github.com/v01dst/pwgen-cli
cd pwgen-cli
go build -o pwgen .
./pwgen -l 24
```

Or with Docker:

```bash
docker build -t pwgen .
docker run --rm pwgen -l 32 -s
```

## 📖 Usage

```
pwgen-cli [flags]

  -l int          password length (default 16, range 4-256)
  -n int          how many passwords (default 1, range 1-100)
  -s              include symbols
  -no-upper       exclude uppercase letters
  -no-digits      exclude digits
  -A              exclude ambiguous chars (il1Lo0O and friends)
  -u              no repeated characters within a password
  -strength       print entropy estimate next to each password
  -version        print version
```

### Examples

```bash
$ pwgen-cli -l 24
9tXcQvR2mKe5wHy7pB3nJd8f

$ pwgen-cli -l 32 -s -n 3 -strength
mV7w(EdXO*FrH^rl-11^=Q$7	153bits-critical
q8?D-Lj*Pe{d+n]DmCCE6<XO	153bits-critical
Q(mS0jav-^Pm%akxh#7B7i>X	153bits-critical

$ pwgen-cli -l 6 -no-upper -no-digits -u   # human-friendly
kvwtnr
```

### Exit codes

| Code | Meaning                     |
|------|-----------------------------|
| `0`  | Success                     |
| `1`  | Invalid flags / generation failure |

## 🧱 Tech Stack

| Layer    | Tech                          |
|----------|-------------------------------|
| Language | Go 1.27 (stdlib only)         |
| RNG      | crypto/rand                   |
| Testing  | go test (table + statistical) |
| CI       | GitHub Actions                |
| Packaging| Docker (multi-stage, alpine)  |

---

<div align="center">

Built with ⚡ by **v01dst**

[![GitHub](https://img.shields.io/badge/github-v01dst-181717?logo=github)](https://github.com/v01dst)
[![Discord](https://img.shields.io/badge/discord-9p.1-5865F2?logo=discord&logoColor=white)](https://discord.com/users/9p.1)

</div>
