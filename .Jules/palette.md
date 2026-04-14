## 2024-05-18 - Hover-Only Element Keyboard Accessibility
**Learning:** Hover-only elements (like the task 'Delete' button) with `opacity-0 group-hover:opacity-100` are completely inaccessible to keyboard-only users as they cannot visually see where their focus is.
**Action:** Always pair `opacity-0` hover effects on interactive elements with `focus-visible:opacity-100` and clear focus rings (e.g., `focus-visible:ring-2`) to ensure keyboard users can discover and use the element.
## 2026-04-07 - Accessible Forms & Keyboard Navigation

**Learning:** Forms using labels without explicit `htmlFor` attributes to associate with their target `id` create a barrier for screen readers. Further, actions hidden behind `opacity-0` hover states completely fail for keyboard navigation users who rely on tabbing, making core functionality unreachable.

**Action:** Always pair `<label htmlFor="id">` with `<input id="id">`. Ensure that any interaction hidden behind a mouse hover also includes a `focus-visible` override (e.g., `focus-visible:opacity-100`) combined with visible focus rings (`focus-visible:ring-2`) and keyboard outlines disabled (`focus-visible:outline-none`) to safely expose functionality to keyboard-only users.
## 2024-05-19 - Actionable Empty States
**Learning:** Plain text empty states ("No tasks currently scheduled") are unhelpful to users and miss an opportunity to guide them. Adding visual weight (dashed borders, icons) and a call-to-action button that programmatically focuses the relevant input form (`document.getElementById("id").focus()`) significantly reduces friction for first-time users and improves accessibility.
**Action:** When designing empty states for lists or tables, always include an actionable button that directs focus to the element needed to populate that list, and ensure the icon used is decorative with `aria-hidden="true"`.
## 2024-04-14 - Playwright Native Select Option Attach State
**Learning:** In Playwright tests, when waiting for dynamically loaded options in a native `<select>` element to become available before interacting with them, default visibility checks (like `state: 'visible'`) often fail and cause timeouts because `<option>` tags are technically non-visible elements within the select.
**Action:** Always use `{ state: 'attached' }` when using `waitForSelector` for native `<select>` options (e.g., `page.waitForSelector('select option[value="mock"]', { state: 'attached' })`) to reliably proceed with test automation after the DOM updates.

## 2024-04-14 - Test Artifact Pollution Prevention
**Learning:** When using temporary automated tools like Playwright to test frontend changes, it's very easy to accidentally commit large binary artifacts (like `.webm` videos) or test dependencies (`playwright`) to the repository, polluting `package.json` and lockfiles.
**Action:** Always place test artifacts in temporary directories OUTSIDE the repository scope, or explicitly `rm` them. Always revert temporary dependency additions (e.g., `git restore package.json pnpm-lock.yaml`) before committing to strictly adhere to the persona constraint of not adding new dependencies.
