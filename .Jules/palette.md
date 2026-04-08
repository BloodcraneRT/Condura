## 2024-05-18 - Hover-Only Element Keyboard Accessibility
**Learning:** Hover-only elements (like the task 'Delete' button) with `opacity-0 group-hover:opacity-100` are completely inaccessible to keyboard-only users as they cannot visually see where their focus is.
**Action:** Always pair `opacity-0` hover effects on interactive elements with `focus-visible:opacity-100` and clear focus rings (e.g., `focus-visible:ring-2`) to ensure keyboard users can discover and use the element.

## 2026-04-07 - Accessible Forms & Keyboard Navigation
**Learning:** Forms using labels without explicit `htmlFor` attributes to associate with their target `id` create a barrier for screen readers. Further, actions hidden behind `opacity-0` hover states completely fail for keyboard navigation users who rely on tabbing, making core functionality unreachable.
**Action:** Always pair `<label htmlFor="id">` with `<input id="id">`. Ensure that any interaction hidden behind a mouse hover also includes a `focus-visible` override (e.g., `focus-visible:opacity-100`) combined with visible focus rings (`focus-visible:ring-2`) and keyboard outlines disabled (`focus-visible:outline-none`) to safely expose functionality to keyboard-only users.

## 2024-05-20 - Form Validation & Action Cues
**Learning:** Required form fields lacking visual cues (like an asterisk) and action buttons that do not clearly indicate when a form is incomplete lead to a confusing user experience, where users might submit forms prematurely and face validation errors instead of receiving upfront guidance.
**Action:** Always include a visual indicator (like a red asterisk) for required fields and explicitly disable submission buttons (e.g., with `disabled` state and `disabled:opacity-50 disabled:cursor-not-allowed`) when required form fields are not yet filled, reducing user friction.
