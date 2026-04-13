## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.
## 2024-05-19 - LimitReader content length with http.NewRequest
**Learning:** `http.NewRequest` cannot automatically infer the `ContentLength` of a request when given an `io.LimitReader` (unlike when given a `*bytes.Reader`). If you stream a dummy payload using an `io.LimitReader` for upload tests to avoid memory allocations, you MUST explicitly set `req.ContentLength` or the request will be sent without a Content-Length header or with chunked encoding in a way that may affect test accuracy.
**Action:** Always explicitly set `req.ContentLength = size` when replacing static byte buffers with streaming readers for HTTP payload generation in synthetic tests.
