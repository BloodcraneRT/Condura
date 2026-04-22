## 2025-02-20 - Logarithmic doubling for buffer filling
**Learning:** When efficiently filling a byte slice with a repeated byte pattern in Go (e.g., in custom io.Reader dummy payload generators), a manual for loop is significantly slower than using a logarithmic doubling approach with copy().
**Action:** Use `p[0] = b; for i := 1; i < len(p); i *= 2 { copy(p[i:], p[:i]) }` when generating repeated patterns to reduce buffer filling time by over 20x.
