## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-05-19 - Streaming Synthetic Payloads Avoids GC Pressure
**Learning:** In synthetic test runners, pre-allocating large byte arrays (e.g., a 10MB slice for upload tests) inside functions creates significant garbage collection pressure and increases memory usage. Using a custom `io.Reader` to stream these payloads on-the-fly dynamically generates the bytes while minimizing allocations.
**Action:** Always prefer streaming data with custom `io.Reader` implementations rather than allocating large, static byte arrays. Be aware that you must explicitly set `req.ContentLength` when creating an `http.NewRequest` with a custom reader since it cannot automatically infer the payload size.
