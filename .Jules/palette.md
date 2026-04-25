## 2024-05-18 - Hover-Only Element Keyboard Accessibility
**Learning:** Hover-only elements (like the task 'Delete' button) with `opacity-0 group-hover:opacity-100` are completely inaccessible to keyboard-only users as they cannot visually see where their focus is.
**Action:** Always pair `opacity-0` hover effects on interactive elements with `focus-visible:opacity-100` and clear focus rings (e.g., `focus-visible:ring-2`) to ensure keyboard users can discover and use the element.
## 2026-04-07 - Accessible Forms & Keyboard Navigation

**Learning:** Forms using labels without explicit `htmlFor` attributes to associate with their target `id` create a barrier for screen readers. Further, actions hidden behind `opacity-0` hover states completely fail for keyboard navigation users who rely on tabbing, making core functionality unreachable.

**Action:** Always pair `<label htmlFor="id">` with `<input id="id">`. Ensure that any interaction hidden behind a mouse hover also includes a `focus-visible` override (e.g., `focus-visible:opacity-100`) combined with visible focus rings (`focus-visible:ring-2`) and keyboard outlines disabled (`focus-visible:outline-none`) to safely expose functionality to keyboard-only users.
## 2024-05-19 - Actionable Empty States
**Learning:** Plain text empty states ("No tasks currently scheduled") are unhelpful to users and miss an opportunity to guide them. Adding visual weight (dashed borders, icons) and a call-to-action button that programmatically focuses the relevant input form (`document.getElementById("id").focus()`) significantly reduces friction for first-time users and improves accessibility.
**Action:** When designing empty states for lists or tables, always include an actionable button that directs focus to the element needed to populate that list, and ensure the icon used is decorative with `aria-hidden="true"`.
## 2024-05-20 - Form Required Indicators and Loading States
**Learning:** Required forms without explicit visual indicators cause friction, and asynchronous submission buttons without loading states and disabled classes (like `disabled:opacity-70 disabled:cursor-not-allowed`) can lead to double submissions and poor user feedback.
**Action:** Always include clear visual indicators (e.g., `<span aria-hidden="true" className="text-rose-500 ml-1">*</span>`) for required fields to aid users. Also, ensure async form submit handlers utilize state (`isSubmitting`) wrapped in a `.finally()` block to consistently toggle disabled states and "Saving..." text on submit buttons.
## 2025-04-25 - Destructive Action Confirmations
**Learning:** Implementing destructive UI actions (such as deleting tasks or configurations) without safeguards leads to accidental data loss and poor user trust.
**Action:** Always wrap destructive actions in a confirmation dialog (e.g., `window.confirm`) to provide a safety net for users before executing API delete requests.
