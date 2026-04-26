## 2024-04-26 - Large Array Fills
**Learning:** Using `copy()` in a logarithmic doubling approach (`for j := 1; j < len(chunk); j *= 2 { copy(chunk[j:], chunk[:j]) }`) is significantly faster than a `for` loop iteration when filling large byte arrays with a single repeated value (e.g., ~12x faster for 1MB chunk). This optimization should not be used for small fixed-size arrays per instructions, but is highly effective for large slices.
**Action:** Replace `for` loop array fills with `copy()` doubling for large payloads/chunks.
