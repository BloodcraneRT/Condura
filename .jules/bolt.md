## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-05-18 - Avoid Large Allocations in HTTP Requests with Custom Readers
**Learning:** Pre-allocating large static byte slices (e.g., `make([]byte, 10*1024*1024)`) to pass into `http.NewRequest` inside recurring tasks significantly increases GC pressure and memory allocations. We can stream these payloads dynamically by implementing a simple custom `io.Reader`. However, since `http.NewRequest` cannot determine the content length of a custom reader automatically, `req.ContentLength` must be set explicitly, or the request might be sent as chunked, causing potential server-side issues or discrepancies.
**Action:** Use custom `io.Reader` implementations to stream large dummy payloads instead of pre-allocating byte slices, but always remember to explicitly set `req.ContentLength`.
