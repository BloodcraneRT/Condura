## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-05-19 - Using io.Reader streaming instead of pre-allocated slices for memory optimization
**Learning:** Pre-allocating large static byte slices (like a 10MB dummy payload array) for synthetic traffic tests creates significant unnecessary memory pressure and garbage collection overhead, especially under concurrent load. Streaming dummy payloads directly to an HTTP request using a custom `io.Reader` implementation eliminates this large upfront allocation entirely. However, when doing this for `http.NewRequest`, you must explicitly set `req.ContentLength` since it cannot be automatically inferred from a custom interface.
**Action:** For performance optimization in synthetic test runners, avoid allocating large static byte slices repeatedly inside functions. Use a custom `io.Reader` implementation to stream synthetic dummy payloads instead to reduce GC pressure and memory usage, and remember to explicitly set the content length.

## 2024-05-18 - Logarithmic Doubling for Fast Array Population
**Learning:** For populating large arrays (like a 1MB payload in `runBERTTest`) with a repeating pattern, using Go's `copy()` with a logarithmic doubling approach is significantly faster than using a loop that calculates modulo arithmetic byte-by-byte. The benchmarks showed the `copy()` method took ~0.3ms vs ~1ms for the loop approach.
**Action:** When working with large byte arrays that need to be populated with repetitive data, pre-calculate the base pattern and then use `for i := len(base); i < len(array); i *= 2 { copy(array[i:], array[:i]) }` instead of looping over the whole array.
