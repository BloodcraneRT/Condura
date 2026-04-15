## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-05-18 - Use Custom io.Reader for Generating Synthetic Payloads
**Learning:** Pre-allocating large byte arrays (e.g., `make([]byte, 10*1024*1024)`) repeatedly during synthetic upload testing causes high GC pressure and unnecessary memory usage.
**Action:** Use a custom `io.Reader` implementation to stream dummy data dynamically for requests. When doing this for `http.NewRequest`, explicitly set `req.ContentLength` since the client cannot automatically infer the size from a custom reader.
