## 2024-05-24 - Optimization of dummyPayloadReader
**Learning:** Found that a manual loop to fill a byte slice (`for i := int64(0); i < toRead; i++ { p[i] = r.b }`) in `dummyPayloadReader.Read` is very slow compared to the logarithmic `copy` method.
**Action:** Replace `for loop` with `p[0] = r.b; for i := 1; i < len; i *= 2 { copy(p[i:], p[:i]) }` for performance improvements.
