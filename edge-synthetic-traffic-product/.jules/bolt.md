## 2024-05-03 - [Logarithmic Doubling for Large Payloads]
**Learning:** [Using a naive byte-by-byte loop to fill large buffers (> 1MB) for synthetic payloads causes significant memory/CPU bottlenecks. Go's native `copy()` function with logarithmic doubling is orders of magnitude faster (12.5µs down to 0.5µs in local bench).]
**Action:** [When filling dynamically sized stream chunks or large arrays with a repeated byte pattern, use `p[0] = b; for i := 1; i < len(p); i *= 2 { copy(p[i:], p[:i]) }`. Avoid for small slices to maintain readability.]
