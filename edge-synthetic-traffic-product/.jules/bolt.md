## 2024-04-23 - Fast Slice Filling with copy()
**Learning:** In Go, when filling a byte slice with a repeated pattern (like dummy payload generation), a manual `for` loop is extremely slow.
**Action:** Always prefer a logarithmic doubling approach using the built-in `copy()` function (e.g., `p[0]=b; for i:=1; i<len(p); i*=2 { copy(p[i:], p[:i]) }`), which leverages optimized memory move instructions for an `O(log n)` speedup.
