
## 2026-04-07 - Focus visibility for hover actions & explicit label links
**Learning:** Hover-only interactive elements completely disappear for keyboard users unless explicitly managed. Additionally, implicit form labels are sometimes insufficient for reliable screen-reader and click-to-focus behavior.
**Action:** Always add `focus-visible:opacity-100 focus-visible:ring-2` (or similar) to elements relying on `group-hover:opacity-100`. Also ensure all form inputs use explicit `id` and `htmlFor` mappings for optimal accessibility.
