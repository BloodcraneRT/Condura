## 2024-05-24 - Unlinked Labels & Inaccessible Hover Buttons
**Learning:** Found that multiple form inputs lacked explicit `id` and `htmlFor` attributes linking them to their labels, and that hover-only buttons (like "Delete") became invisible during keyboard navigation due to missing focus states.
**Action:** Ensure all form labels have explicit `id` and `htmlFor` pairings, and ensure all hover-revealed interactive elements include `focus-visible:opacity-100` alongside their hover state for keyboard users.
