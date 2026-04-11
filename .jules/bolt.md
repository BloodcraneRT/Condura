## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-04-11 - Stream large mock payloads directly
**Learning:** For synthetic tests uploading large mock payloads (e.g., 10MB test data), pre-allocating large byte arrays consumes excessive memory and causes GC pressure.
**Action:** Implement a custom `io.Reader` that streams dummy bytes on the fly. When creating the request with `http.NewRequest("POST", target, reader)`, explicitly set `req.ContentLength` as it cannot be automatically inferred from a custom `io.Reader`.
