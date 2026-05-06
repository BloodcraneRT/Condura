## 2024-05-06 - Initial

## 2024-05-06 - Precalculating Deterministic Sequences
**Learning:** For large byte slices (e.g., 1MB+ streams), dynamically evaluating a math expression (like `(i*13)%256`) byte-by-byte in a `for` loop is computationally expensive and slow compared to bulk memory copying.
**Action:** Precalculate the repeating pattern (e.g., 256 bytes) into a small base slice, then use a logarithmic doubling approach with `copy()` (`p[0] = b; for i := 1; i < len(p); i *= 2 { copy(p[i:], p[:i]) }`) to fill the rest. Do not use this micro-optimization for small slices (<1500 bytes) to maintain readability.
