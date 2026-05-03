## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-05-19 - Using io.Reader streaming instead of pre-allocated slices for memory optimization
**Learning:** Pre-allocating large static byte slices (like a 10MB dummy payload array) for synthetic traffic tests creates significant unnecessary memory pressure and garbage collection overhead, especially under concurrent load. Streaming dummy payloads directly to an HTTP request using a custom `io.Reader` implementation eliminates this large upfront allocation entirely. However, when doing this for `http.NewRequest`, you must explicitly set `req.ContentLength` since it cannot be automatically inferred from a custom interface.
**Action:** For performance optimization in synthetic test runners, avoid allocating large static byte slices repeatedly inside functions. Use a custom `io.Reader` implementation to stream synthetic dummy payloads instead to reduce GC pressure and memory usage, and remember to explicitly set the content length.

## 2024-05-20 - Logarithmic Doubling for Large Repeating Payloads
**Learning:** Generating large deterministic repeating payloads (like a 1MB pseudo-random binary sequence) byte-by-byte using mathematical modulo operations inside a loop has significant CPU overhead. By precalculating the base repeating pattern (e.g., 256 bytes) and using logarithmic doubling with `copy()`, performance improved by ~6x (from ~1.03ms to ~0.17ms per iteration) due to optimized memory block copying vs individual math evaluations.
**Action:** When filling large buffers with a repeating pattern, use the logarithmic doubling approach with `copy()` instead of evaluating the pattern logically for each byte.
