# Desktop mockups — Learn the Grammar

Static, standalone HTML/CSS mockups of the desktop user app. Open any file directly
in a browser — they link to each other by filename, so the flow is click-through.

## Flow

1. `Login.html` — sign in / sign up (tab switch is live JS)
2. `StartLesson.html` — home screen, has an "unfinished lesson" state
3. `Generating.html` — loading state after "Start a new lesson" (skeleton + spinner overlay)
4. `Lesson.html` — the 5 sentences to translate
5. `Grading.html` — loading state after "Submit for grading" (dimmed answers + spinner overlay)
6. `LessonResult.html` — scores, advice, grammar concepts, sentence-by-sentence review
7. `Settings.html` — level picker + collapsible topics grid (live JS)
8. `Dictionary.html` — word list with topic filter and sort (live JS)
9. `Profile.html` — average score, vocabulary progress, grammar concept ratings

`Settings.html` and `Dictionary.html` have real vanilla-JS behavior (no framework) —
read the `<script>` block at the bottom of each for the intended interaction logic
(collapsible topics, single-topic search filter, sort by stars/alphabet). Treat that
JS as a spec to reimplement in React, not as production code.

## Design tokens (see `:root` CSS vars in each file)

- Background `#F6F5F1`, surface `#FFFFFF`, border `#E3E1DA`
- Text `#201F1C` / secondary `#6B6862`
- Accent (buttons, links, active states) `#2C5F8A`, dark `#1F4763`, soft bg `#E8F0F6`
- Gold accent (stars, logo swash) `#C99A3B`
- Success `#1A7F37` / soft `#DAFBE1`, Danger `#CF222E` / soft `#FFEBE9`, Warning `#C1780A` / soft `#FBF1DE`
- Fonts: Sora (headings/display), Plus Jakarta Sans (body), Caveat (handwritten logo wordmark) — all via Google Fonts
- Radius `12px` on cards, `10px` on inputs/buttons

## Not built yet

- Admin console (topics/word generation/meanings) — pending an updated PRD, since
  topics/words are now sourced from a Google Sheet rather than managed in-app.
- Mobile — desktop only so far.

Source of truth for behavior: `doc/PRD.md` in the repo.
