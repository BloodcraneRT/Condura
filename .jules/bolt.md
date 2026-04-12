## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.
## 2024-05-19 - Streaming Dummy Payloads for Large Requests
**Learning:** Eagerly allocating large byte slices (e.g., 10MB) for synthetic testing payloads inside a frequently executed function (like `runUploadTest`) significantly increases GC pressure and memory spikes. Replacing it with a custom `io.Reader` implementation (`dummyPayloadReader`) effectively eliminates this allocation.
**Action:** When creating large, repeated synthetic payloads for `http.NewRequest`, use a custom `io.Reader` to stream the payload dynamically, but remember to explicitly set `req.ContentLength` since Go's HTTP client cannot infer length from custom readers.
