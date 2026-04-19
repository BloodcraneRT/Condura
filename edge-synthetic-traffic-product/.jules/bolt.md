## 2025-04-19 - Fast slice filling for custom readers
**Learning:** Using a simple `for` loop to manually copy bytes one-by-one into a slice is a significant performance bottleneck when creating large dummy payloads in custom `io.Reader` implementations.
**Action:** Use a logarithmic doubling approach with `copy()` for slice filling (`p[0] = b; for i := 1; i < len(p); i *= 2 { copy(p[i:], p[:i]) }`), which is dramatically faster.
