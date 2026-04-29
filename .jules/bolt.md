## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-05-19 - Using io.Reader streaming instead of pre-allocated slices for memory optimization
**Learning:** Pre-allocating large static byte slices (like a 10MB dummy payload array) for synthetic traffic tests creates significant unnecessary memory pressure and garbage collection overhead, especially under concurrent load. Streaming dummy payloads directly to an HTTP request using a custom `io.Reader` implementation eliminates this large upfront allocation entirely. However, when doing this for `http.NewRequest`, you must explicitly set `req.ContentLength` since it cannot be automatically inferred from a custom interface.
**Action:** For performance optimization in synthetic test runners, avoid allocating large static byte slices repeatedly inside functions. Use a custom `io.Reader` implementation to stream synthetic dummy payloads instead to reduce GC pressure and memory usage, and remember to explicitly set the content length.

## 2024-05-20 - Logarithmic Doubling for Fast Payload Generation
**Learning:** Using a single byte-assignment loop to fill a large pre-allocated slice (e.g., 1MB+) or dynamic stream chunks for synthetic tests is inefficient and CPU-bound.
**Action:** When filling large or dynamically sized byte slices with repeating patterns, pre-calculate the base repeating sequence (if any) up to its period, then use logarithmic doubling with `copy()` (`for i := baseLen; i < len(p); i *= 2 { copy(p[i:], p[:i]) }`) to fill the remainder of the slice significantly faster.
