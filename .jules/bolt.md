## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-05-18 - Streaming io.Reader Reduces GC Pressure
**Learning:** When making `http.NewRequest` calls with large dummy payloads in Go, allocating a large static byte slice (e.g., `make([]byte, 10*1024*1024)`) per request creates significant GC pressure during concurrent tests. Replacing it with a custom streaming `io.Reader` implementation eliminates the allocation completely.
**Action:** Use custom `io.Reader` implementations to stream synthetic dummy payloads instead of pre-allocating large byte arrays. When doing this for `http.NewRequest`, you must explicitly set `req.ContentLength` since it cannot be automatically inferred from a custom reader, otherwise it will default to chunked encoding or fail.
