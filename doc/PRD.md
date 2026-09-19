# PRD: Learn the Grammar — Language Grammar & Writing Trainer

**Status:** Draft, implementation-ready except for the vocabulary selection algorithm in FR-8  
**Version:** 2.0  
**Language:** English  

## 1. Overview

Learn the Grammar is a web application that trains active grammar, spelling, and word choice through translation exercises. The LLM writes sentences in the User's native language, the User translates them into the language being learned, and the LLM grades the translations. The application keeps a personal dictionary and uses each User's progress to prepare future Lessons.

The first release supports Russian as the native/source language and Greek as the target language. The data model must allow additional languages to be added later without recreating meanings or losing existing per-language progress.

This document is the product source of truth for v1. Every acceptance criterion is identified by its functional requirement and criterion number, for example `FR-14.AC-3`.

## 2. Goals, scope, and non-goals

### Goals

- **G-1:** Improve active grammar, spelling, and word choice through regular translation practice.
- **G-2:** Adapt Lessons to the User's current vocabulary progress so weak and new words can receive practice.
- **G-3:** Show recent Lesson performance and progress through the active dictionary.

### V1 scope

- Language pair: Russian as native/source language and Greek as target language.
- Catalog languages: English, Russian, and Greek are stored for every Meaning.
- Levels: `A1`, `A2`, and `B1`.
- Platform: a responsive web User app and a separate Admin console.
- Access: only authenticated Users can take or continue Lessons.
- Configuration: Lesson size, average window, and words-per-sentence targets are deployment settings rather than User settings.

### Non-goals

- **NG-1:** Listening, speaking, or pronunciation practice.
- **NG-2:** Levels above `B1`.
- **NG-3:** Free-form conversation with the LLM.
- **NG-4:** User-selectable language pairs in v1.

## 3. Lesson flow

A Lesson has three stages:

1. The LLM generates several Russian source sentences and hidden suggested Greek translations.
2. The User enters Greek translations and submits any current set of answers, including incomplete answers.
3. The LLM grades the answers, returns feedback and per-entity rating deltas, and the application finalizes the result exactly once.

## 4. Roles

- **Guest:** an unauthenticated visitor.
- **User:** an authenticated learner.
- **Admin:** an authenticated user with administrative access.

## 5. Glossary and core domain model

- **Level:** one of `A1`, `A2`, or `B1`.
- **Topic:** a vocabulary category. A Topic can be active or inactive.
- **Meaning:** a language-independent concept. It belongs to exactly one Topic and one Level, and can be active or inactive.
- **Translation:** one language-specific expression of a Meaning. English, Russian, and Greek are stored for each Meaning in v1, with at most one Translation per language.
- **Translation rating:** a User's progress for one specific Translation. It is either `unseen` or a numeric value in `[-1, +3]`.
- **Grammar Concept:** a language-specific grammar item with an unlock Level.
- **Grammar Concept rating:** a User's progress for one Grammar Concept, in `[-1, 10]`.
- **Lesson:** a generated, unfinished, abandoned, or completed set of source sentences and the snapshots needed to grade them consistently.
- **Completed lesson:** a Lesson for which a valid grading result has been successfully finalized.

## 6. Screens

This inventory defines which v1 screens exist and their purpose. Detailed layout is not prescribed by this PRD.

### User app

| Screen | Purpose |
| --- | --- |
| Login and sign-up | Sign in, create an account, and request password recovery. |
| Settings | Select the Level, enable or disable active Topics, and explicitly save changes. |
| Start Lesson | Start a Lesson; when an unfinished Lesson exists, continue it or end it and start another. |
| Lesson | Show source sentences, one answer field per sentence, contextual word hints, and Submit. |
| Lesson result | Show criterion scores, overall score, advice, and a correct/reference/alternative version for every sentence. |
| Dictionary | Browse main-dictionary, disabled-Topic, and higher-Level vocabulary with personal ratings. |
| Profile | Show recent score averages, vocabulary progress, active Grammar Concept ratings, and current settings. |

### Admin console

| Screen | Purpose |
| --- | --- |
| Update words | Run `Update words from connected sheet` and show validation, access, interruption, and completion results. |

## 7. Global invariants

1. A User can have at most one unfinished Lesson at a time.
2. A Lesson can be finalized at most once.
3. Rating and statistics side effects from a Lesson can be applied at most once.
4. Topic, Meaning, Translation, and Grammar Concept identity is stable by system ID.
5. Deactivation is non-destructive: historical ratings and completed Lesson data are retained.
6. Existing unfinished Lessons use immutable generation-time snapshots and are not changed by later Settings or catalog changes.
7. LLM output is authoritative for semantic judgments, word usage, and rating deltas. The application validates the response shape and value ranges but does not replace the LLM's semantic judgment with string matching or local lemmatization.

## 8. Functional requirements and acceptance criteria

### Epic A — Setup and access

### FR-1 — Initial settings

A new User must review and explicitly save the initial learning settings.

- **FR-1.AC-1:** A new User is created with Level `A1`.
- **FR-1.AC-2:** All active Topics are enabled for a new User.
- **FR-1.AC-3:** After the first successful sign-in, the User is directed to Settings.
- **FR-1.AC-4:** Settings contains an explicit `Save` action.
- **FR-1.AC-5:** Onboarding is not complete until Settings have been successfully saved.
- **FR-1.AC-6:** If the User leaves before saving, the next sign-in directs the User to Settings again.
- **FR-1.AC-7:** After the initial successful Save, later sign-ins do not automatically direct the User to Settings.
- **FR-1.AC-8:** A Topic activated or added after the User registered becomes enabled for that User automatically.
- **FR-1.AC-9:** The application does not show a special notification or marker solely because a Topic was enabled automatically.

### FR-2 — Topic settings

A User can control which active Topics participate in learning.

- **FR-2.AC-1:** A User can disable an enabled Topic.
- **FR-2.AC-2:** A User can re-enable a disabled Topic.
- **FR-2.AC-3:** Translations from disabled Topics are excluded from the main dictionary.
- **FR-2.AC-4:** Translations from disabled Topics are not selected for new Lessons.
- **FR-2.AC-5:** Disabling a Topic does not delete or change existing Translation ratings.
- **FR-2.AC-6:** Re-enabling a Topic makes its existing ratings applicable again.
- **FR-2.AC-7:** Topic-setting changes do not change an existing unfinished Lesson.
- **FR-2.AC-8:** Topic-setting changes do not change completed Lesson results.
- **FR-2.AC-9:** Topic-setting changes take effect only after an explicit successful Save.

### FR-18 — Supported language pair

The v1 learner experience uses Russian as the native language and Greek as the target language.

- **FR-18.AC-1:** The v1 UI supports Russian source sentences and Greek answers.
- **FR-18.AC-2:** A User cannot change the language pair in v1.
- **FR-18.AC-3:** Translation ratings belong to a specific Translation and therefore to a specific language, not directly to a Meaning.

### FR-26 — Lesson access

Only authenticated Users can enter the Lesson flow.

- **FR-26.AC-1:** A Guest cannot start a Lesson.
- **FR-26.AC-2:** A Guest cannot continue an unfinished Lesson.
- **FR-26.AC-3:** An authenticated User or Admin can open the Lesson flow.
- **FR-26.AC-4:** An unauthenticated request does not create a Lesson or invoke LLM generation.

### FR-29 — Sign-up

A Guest can create an account with email and password.

- **FR-29.AC-1:** Sign-up accepts an email address and password.
- **FR-29.AC-2:** A password must contain at least 8 characters.
- **FR-29.AC-3:** The email address is validated before account creation.
- **FR-29.AC-4:** Email addresses are normalized and compared case-insensitively.
- **FR-29.AC-5:** Sign-up with an already registered normalized email fails.
- **FR-29.AC-6:** Sign-up with an invalid email fails.
- **FR-29.AC-7:** Sign-up with a password shorter than 8 characters fails.
- **FR-29.AC-8:** Email confirmation is not required in v1.
- **FR-29.AC-9:** After successful registration, the User is signed in automatically.
- **FR-29.AC-10:** The initial-settings flow in FR-1 starts after automatic sign-in.

### FR-31 — Change level

A User can change the current learning Level.

- **FR-31.AC-1:** The User can select `A1`, `A2`, or `B1`.
- **FR-31.AC-2:** A saved Level change applies to new Lessons.
- **FR-31.AC-3:** The main dictionary is recalculated for the new Level.
- **FR-31.AC-4:** Existing Translation ratings are preserved.
- **FR-31.AC-5:** Existing Grammar Concept ratings are preserved.
- **FR-31.AC-6:** Completed Lesson history and averages are preserved.
- **FR-31.AC-7:** An existing unfinished Lesson continues to use the settings snapshot with which it was created.

### FR-36 — Incorrect password

The system must reject an incorrect password without authenticating the account.

- **FR-36.AC-1:** Sign-in with a registered normalized email and an incorrect password fails.
- **FR-36.AC-2:** The account is not authenticated.
- **FR-36.AC-3:** No new authenticated session is created.

### FR-37 — Password recovery

A User can reset a forgotten password through a time-limited, single-use link.

- **FR-37.AC-1:** A User can request a password-reset link by email.
- **FR-37.AC-2:** A reset link expires 4 hours after it is issued.
- **FR-37.AC-3:** A reset link can be used only once.
- **FR-37.AC-4:** At most one reset link is valid for a User at any time.
- **FR-37.AC-5:** Issuing a new reset link invalidates the previous link.
- **FR-37.AC-6:** A successfully used link becomes invalid immediately.
- **FR-37.AC-7:** After a successful reset, the old password no longer permits sign-in.
- **FR-37.AC-8:** The new password follows the same requirements as a sign-up password.

### Epic B — Dictionary

### FR-3 — Topic and Level organization

Meanings are organized by Topic and Level.

- **FR-3.AC-1:** Each Meaning belongs to exactly one Topic.
- **FR-3.AC-2:** Each Meaning belongs to exactly one Level: `A1`, `A2`, or `B1`.
- **FR-3.AC-3:** A Topic can contain no Meanings for a particular Level.
- **FR-3.AC-4:** An empty `(Topic, Level)` combination is valid.

### FR-4 — Main dictionary by Level

The main dictionary contains active, enabled content at or below the User's current Level.

- **FR-4.AC-1:** At `A1`, the main dictionary includes eligible `A1` Translations.
- **FR-4.AC-2:** At `A2`, the main dictionary includes eligible `A1` and `A2` Translations.
- **FR-4.AC-3:** At `B1`, the main dictionary includes eligible `A1`, `A2`, and `B1` Translations.
- **FR-4.AC-4:** Translations above the current Level are excluded.
- **FR-4.AC-5:** Translations in disabled or inactive Topics are excluded.
- **FR-4.AC-6:** Translations of inactive Meanings are excluded.
- **FR-4.AC-7:** Moving a Translation into or out of the main dictionary does not delete its rating.

### FR-5 — Browse dictionary

A User can inspect vocabulary grouped by its current relationship to the main dictionary.

- **FR-5.AC-1:** The User can open a Dictionary view.
- **FR-5.AC-2:** The view has separate groups for the main dictionary, disabled Topics, and higher Levels.
- **FR-5.AC-3:** Each visible Translation is placed in the group determined by the User's current saved Settings.
- **FR-5.AC-4:** A Translation with no completed graded usage is shown as `unseen`, without a numeric rating.
- **FR-5.AC-5:** `unseen` is distinct from numeric rating `0`.
- **FR-5.AC-6:** Once a numeric rating exists, it is shown to the User.
- **FR-5.AC-7:** Meanings and Topics that are inactive are not shown anywhere in the User-facing Dictionary.

### FR-7 — Translation rating

The system tracks each User's mastery of each target-language Translation.

- **FR-7.AC-1:** A rating belongs to `(User, Translation)`, not to a Meaning as a whole.
- **FR-7.AC-2:** Before the first rating event, a Translation is `unseen`.
- **FR-7.AC-3:** The first rating event creates a numeric rating.
- **FR-7.AC-4:** A numeric rating is limited to `[-1, +3]`.
- **FR-7.AC-5:** The grading LLM can return `+1`, `-1`, or `0` for a Translation occurrence.
- **FR-7.AC-6:** `+1` increases the numeric rating.
- **FR-7.AC-7:** `-1` decreases the numeric rating.
- **FR-7.AC-8:** `0` does not change the numeric rating.
- **FR-7.AC-9:** A rating never exceeds `+3`.
- **FR-7.AC-10:** A rating never falls below `-1`.
- **FR-7.AC-11:** Rating `+3` means mastered.
- **FR-7.AC-12:** A net positive update at `+3` leaves the rating at `+3`.
- **FR-7.AC-13:** A net negative update at `+3` lowers the rating.
- **FR-7.AC-14:** A net negative update at `-1` leaves the rating at `-1`.
- **FR-7.AC-15:** A net `+1` update at `-1` changes the rating to `0`.
- **FR-7.AC-16:** A mastered Translation remains eligible for Lessons.

### FR-19 — Future language support

The data architecture must support additional languages even though v1 exposes only Russian-to-Greek learning.

- **FR-19.AC-1:** Meaning identity is independent of any native or target language.
- **FR-19.AC-2:** A Translation for a new language can be added to an existing Meaning without creating a new Meaning.
- **FR-19.AC-3:** Adding a new-language Translation does not change ratings for existing Translations.
- **FR-19.AC-4:** A User-facing language-management UI is not required in v1.

### FR-21 — Meaning translations

A Meaning has a single canonical Translation per language in v1.

- **FR-21.AC-1:** A Meaning can have no more than one Translation for the same language.
- **FR-21.AC-2:** Synonyms are not stored as multiple same-language Translations of one Meaning in v1.
- **FR-21.AC-3:** Import rejects a row that would violate the one-Translation-per-language rule.
- **FR-21.AC-4:** Required Topic, Level, and Translation fields must be present and valid before an import can begin.
- **FR-21.AC-5:** Every importable v1 Meaning row contains exactly one English, one Russian, and one Greek Translation.
- **FR-21.AC-6:** A new Meaning receives a stable system ID when first loaded, and that ID is written back to the connected Sheet.

### FR-22 — Personal dictionary

Dictionary progress is personal and language-specific.

- **FR-22.AC-1:** Ratings for the same Translation are independent between Users.
- **FR-22.AC-2:** A User's rating for one Translation does not affect another Translation of the same Meaning.
- **FR-22.AC-3:** A future target language can have separate User ratings for its Translations.

### FR-23 — Deactivate and reactivate a Meaning

An Admin can remove a Meaning from active use without deleting historical data.

- **FR-23.AC-1:** Deactivating a Meaning sets it to inactive rather than hard-deleting it.
- **FR-23.AC-2:** An inactive Meaning is excluded from every User's dictionary.
- **FR-23.AC-3:** An inactive Meaning is excluded from new Lesson generation.
- **FR-23.AC-4:** Existing Translation ratings and historical Lesson data are retained.
- **FR-23.AC-5:** Reactivating the Meaning removes the inactive state.
- **FR-23.AC-6:** After reactivation, existing User ratings become applicable again.
- **FR-23.AC-7:** Deactivation or reactivation does not alter an existing unfinished Lesson snapshot.
- **FR-23.AC-8:** The Admin controls the active/inactive state through the connected Sheet; a subsequent successful import applies it.

### FR-28 — Deactivate a Topic

Removing a Topic from active use is non-destructive.

- **FR-28.AC-1:** Deleting a Topic through the connected-Sheet workflow sets it to inactive rather than hard-deleting it.
- **FR-28.AC-2:** An inactive Topic and all of its Meanings are excluded from User-facing dictionaries.
- **FR-28.AC-3:** An inactive Topic and all of its Meanings are excluded from new Lesson generation.
- **FR-28.AC-4:** User Topic settings, Translation ratings, and historical Lesson data are retained.
- **FR-28.AC-5:** Topic deactivation does not alter existing unfinished Lesson snapshots.
- **FR-28.AC-6:** Profile and progress calculations include only active Topics.
- **FR-28.AC-7:** The Admin can add a Topic at any time by listing it in the connected Sheet and running a successful import.
- **FR-28.AC-8:** A newly imported active Topic follows FR-1's automatic-enable behavior for existing Users.

### Epic C — Grammar

### FR-6 — Grammar Concept unlock

Grammar Concepts become available cumulatively by Level.

- **FR-6.AC-1:** Each Grammar Concept belongs to one language.
- **FR-6.AC-2:** Each Grammar Concept has one unlock Level.
- **FR-6.AC-3:** `A1` activates `A1` concepts.
- **FR-6.AC-4:** `A2` activates `A1` and `A2` concepts.
- **FR-6.AC-5:** `B1` activates `A1`, `A2`, and `B1` concepts.
- **FR-6.AC-6:** A concept above the current Level is not supplied for new Lesson generation.
- **FR-6.AC-7:** Lowering the Level does not delete historical Grammar Concept ratings.
- **FR-6.AC-8:** Raising the Level again restores access to the previously stored rating.

### FR-24 — Grammar Concept ratings

The system tracks each User's progress for each Grammar Concept.

- **FR-24.AC-1:** A rating belongs to `(User, Grammar Concept)`.
- **FR-24.AC-2:** The rating range is `[-1, 10]`.
- **FR-24.AC-3:** The grading LLM returns `+1`, `-1`, or `0` for each evaluated Grammar Concept.
- **FR-24.AC-4:** `+1`, `-1`, and `0` respectively increase, decrease, or preserve the rating.
- **FR-24.AC-5:** Rating updates saturate at `-1` and `10`.
- **FR-24.AC-6:** Grammar Concept ratings do not affect Lesson generation in v1.

### FR-35 — Grammar profile

The Profile shows progress for currently active Grammar Concepts.

- **FR-35.AC-1:** The Profile shows Grammar Concepts unlocked at the current Level.
- **FR-35.AC-2:** Higher-Level concepts are not shown.
- **FR-35.AC-3:** A hidden higher-Level concept retains its historical rating.
- **FR-35.AC-4:** Returning to a Level that unlocks the concept shows its retained rating.
- **FR-35.AC-5:** Concepts for inactive content or unsupported languages are not shown in the v1 Profile.

### Epic D — Lesson generation

### FR-8 — Generation input and vocabulary selection

Lesson generation uses the User's saved configuration and selected vocabulary.

- **FR-8.AC-1:** Generation receives the User's current saved Level.
- **FR-8.AC-2:** Generation receives active Grammar Concepts unlocked at the current or lower Levels.
- **FR-8.AC-3:** Generation receives selected Translations and their Meanings from the main dictionary.
- **FR-8.AC-4:** Disabled or inactive content is not supplied as selected Lesson vocabulary.
- **FR-8.AC-5:** Content above the User's current Level is not supplied.
- **FR-8.AC-6:** The vocabulary target equals `configured sentence count × configured words per sentence for the current Level`.
- **FR-8.AC-7:** Default words-per-sentence values are `A1=4`, `A2=7`, and `B1=10`.
- **FR-8.AC-8:** Words-per-sentence values are deployment configuration, not User settings.
- **FR-8.AC-9:** If the main dictionary contains fewer eligible Translations than the target, all eligible Translations are supplied.
- **FR-8.AC-10:** When the supplied list is smaller than the target, the LLM may use additional vocabulary outside the supplied list to form natural sentences.
- **FR-8.AC-11:** If the main dictionary is empty, generation is not invoked.

> **BLOCKED — vocabulary selection algorithm:** The rule for selecting the requested mix of new, mistaken, and mastered Translations has not been specified. An implementation agent must not invent this algorithm. FR-8 is not implementation-ready until this rule is provided.

### FR-9 — Generation output

The LLM returns the complete material and metadata needed to take and later grade a Lesson.

- **FR-9.AC-1:** Successful generation returns the configured number of Russian source sentences.
- **FR-9.AC-2:** Each source sentence has a hidden suggested Greek translation.
- **FR-9.AC-3:** Suggested translations are not displayed in the main Lesson-taking UI.
- **FR-9.AC-4:** The response identifies which supplied Translations were used and which were not used.
- **FR-9.AC-5:** LLM usage metadata is the source of truth; the application does not independently lemmatize Greek text to determine usage.
- **FR-9.AC-6:** A missing required sentence, suggested translation, or usage field makes the response unusable.
- **FR-9.AC-7:** An unusable response is handled as an LLM failure under FR-33.
- **FR-9.AC-8:** A failed generation does not create a partial Lesson.

### FR-27 — Deployment configuration

Operational lesson-size and profile-window values are configured at deployment time.

- **FR-27.AC-1:** The default number of sentences per Lesson is `5`.
- **FR-27.AC-2:** The default rolling-average window `N` is `5` completed Lessons.
- **FR-27.AC-3:** The default words per sentence are `A1=4`, `A2=7`, and `B1=10`.
- **FR-27.AC-4:** Sentence count, `N`, and words-per-sentence values are deployment configuration.
- **FR-27.AC-5:** Users cannot change these values through the v1 UI.
- **FR-27.AC-6:** A new Lesson uses the configuration current when that Lesson is generated.
- **FR-27.AC-7:** The Profile uses the currently configured `N`.

### FR-30 — Empty main dictionary

A Lesson cannot start without at least one eligible Translation.

- **FR-30.AC-1:** Before generation, the application checks for at least one eligible active Translation in the main dictionary.
- **FR-30.AC-2:** If no eligible Translation exists, no Lesson is created.
- **FR-30.AC-3:** The LLM generation service is not invoked.
- **FR-30.AC-4:** The User sees an error state explaining that a Lesson cannot be generated.
- **FR-30.AC-5:** The User can navigate to Settings to change the Level or enabled Topics.

### FR-33 — LLM failure

Generation and grading failures must be retryable without duplicate side effects.

- **FR-33.AC-1:** For a failed generation or grading call, the backend performs exactly one automatic retry.
- **FR-33.AC-2:** If the retry succeeds, the flow continues normally.
- **FR-33.AC-3:** If the retry also fails, the User sees `Oops, something went wrong`.
- **FR-33.AC-4:** After failed grading, entered translations remain saved.
- **FR-33.AC-5:** The User can press Submit again after failed grading.
- **FR-33.AC-6:** A repeated or concurrent Submit cannot finalize the same Lesson twice.
- **FR-33.AC-7:** Translation ratings, Grammar Concept ratings, and Lesson statistics are applied at most once.
- **FR-33.AC-8:** After failed generation, the User can start generation again.
- **FR-33.AC-9:** While an LLM request is pending, the relevant screen is blocked by a loading state to prevent conflicting actions.
- **FR-33.AC-10:** A structurally valid but product-invalid response, including missing required data or out-of-range values, is treated as an unusable response and follows the same retry behavior.

### Epic E — Taking a Lesson

### FR-10 — Contextual word hint

A User can request a base-form hint without revealing the completed answer.

- **FR-10.AC-1:** The User can request a hint for a particular word in a source sentence.
- **FR-10.AC-2:** The LLM generates the hint using the sentence context.
- **FR-10.AC-3:** The hint returns the dictionary or base form of the target-language word.
- **FR-10.AC-4:** The hint does not reveal the exact inflected form required by the sentence.
- **FR-10.AC-5:** The hint does not reveal the full suggested sentence translation.
- **FR-10.AC-6:** Requesting a hint does not directly change Lesson scores.
- **FR-10.AC-7:** Requesting a hint does not directly change a Translation rating.
- **FR-10.AC-8:** Requesting a hint does not directly change a Grammar Concept rating.

### FR-11 — Submit translations

A User can submit any current set of answers, including incomplete answers, for grading.

- **FR-11.AC-1:** The User can enter a separate Greek translation for each source sentence.
- **FR-11.AC-2:** The User can submit with one or more answer fields empty.
- **FR-11.AC-3:** An empty answer is passed to the grading LLM as no answer rather than blocked by the UI.
- **FR-11.AC-4:** The grading LLM can return negative word-rating deltas for vocabulary not demonstrated because an answer was empty.
- **FR-11.AC-5:** Submit grades a snapshot of the answers current at submission time.
- **FR-11.AC-6:** While grading is in progress, another concurrent Submit cannot create a second finalization.

### FR-32 — Unfinished Lesson

The User can safely leave and continue the latest unfinished Lesson.

- **FR-32.AC-1:** A User can have at most one unfinished Lesson.
- **FR-32.AC-2:** An unfinished Lesson stores an immutable snapshot of generated sentences and hidden suggested translations.
- **FR-32.AC-3:** It stores the relevant Level, Topic, selected-vocabulary, Grammar Concept, and generation-configuration snapshots.
- **FR-32.AC-4:** Entered answers are saved automatically using a debounced or equivalent mechanism.
- **FR-32.AC-5:** Losing the last few seconds of typing is acceptable.
- **FR-32.AC-6:** After refresh, browser close, or re-login, the User can continue from the latest saved state.
- **FR-32.AC-7:** Starting a Lesson while an unfinished Lesson exists offers `Continue` or `End and start new`.
- **FR-32.AC-8:** `Continue` opens the same sentences and the latest saved answers.
- **FR-32.AC-9:** `End and start new` closes the old Lesson without a grading result.
- **FR-32.AC-10:** Ending an unfinished Lesson does not change ratings or Lesson statistics.
- **FR-32.AC-11:** Later Settings or catalog changes do not alter the unfinished Lesson snapshot.

### Epic F — Grading and feedback

### FR-12 — Lesson score

The LLM grades three integer criteria, and the application calculates the overall score.

- **FR-12.AC-1:** The LLM returns an integer score from `1` to `10` for **Grammar**.
- **FR-12.AC-2:** The LLM returns an integer score from `1` to `10` for **Spelling**.
- **FR-12.AC-3:** The LLM returns an integer score from `1` to `10` for **Vocabulary and meaning**.
- **FR-12.AC-4:** The application, not the LLM, calculates the overall score.
- **FR-12.AC-5:** Overall score is the arithmetic mean of the three criterion scores.
- **FR-12.AC-6:** Overall score is displayed with exactly two decimal places.
- **FR-12.AC-7:** Translation ratings are not inputs to the overall Lesson score.
- **FR-12.AC-8:** Grammar Concept ratings are not inputs to the overall Lesson score.
- **FR-12.AC-9:** For example, scores `8`, `7`, and `8` produce overall score `7.67`.

### FR-13 — Feedback

The completed-Lesson result explains the score and provides reference translations.

- **FR-13.AC-1:** The result shows all three criterion scores.
- **FR-13.AC-2:** The result shows the overall score.
- **FR-13.AC-3:** The LLM can return from `0` to `3` textual advice items.
- **FR-13.AC-4:** Advice is free text and can call attention to a grammar or language issue.
- **FR-13.AC-5:** For every source sentence, the result shows a correct, reference, or alternative Greek version.
- **FR-13.AC-6:** A reference version is shown even when the User's answer was correct.
- **FR-13.AC-7:** The reference version can be semantically equivalent to the User's answer.

### FR-25 — Grammar grading

The grading LLM decides one delta for every relevant Grammar Concept.

- **FR-25.AC-1:** The grading response returns `+1`, `-1`, or `0` for every relevant Grammar Concept.
- **FR-25.AC-2:** The LLM decides the final delta for a concept.
- **FR-25.AC-3:** The application does not independently aggregate separate grammar occurrences into another delta.
- **FR-25.AC-4:** Deltas are applied according to FR-24.
- **FR-25.AC-5:** Updates are applied only during successful Lesson finalization.
- **FR-25.AC-6:** A retry or repeated finalization request cannot apply the same update twice.

### Epic G — Ratings and statistics

### FR-14 — Translation rating updates

The grading LLM evaluates every occurrence of supplied vocabulary that was used in the Lesson.

- **FR-14.AC-1:** For every used Translation occurrence, the grading response returns `+1`, `-1`, or `0`.
- **FR-14.AC-2:** Every occurrence is evaluated separately.
- **FR-14.AC-3:** The LLM response is the source of truth for correctness and delta.
- **FR-14.AC-4:** The application does not infer correctness through string matching.
- **FR-14.AC-5:** A supplied Translation not used in the generated sentences receives no rating update.
- **FR-14.AC-6:** For each Translation, all occurrence deltas in the Lesson are summed first.
- **FR-14.AC-7:** The summed delta is applied once to the pre-Lesson rating and then clamped to `[-1, +3]`.
- **FR-14.AC-8:** This calculation is independent of sentence or occurrence order.
- **FR-14.AC-9:** Rating updates are applied only after successful grading.
- **FR-14.AC-10:** An abandoned Lesson does not change ratings.
- **FR-14.AC-11:** One grading result can be applied at most once.

### FR-15 — Store Lesson scores

Only successfully graded Lessons enter score history.

- **FR-15.AC-1:** A Lesson enters score history only after a valid grading result is successfully finalized.
- **FR-15.AC-2:** A generated-only Lesson is not completed.
- **FR-15.AC-3:** A grading-failed Lesson is not completed and remains retryable.
- **FR-15.AC-4:** An abandoned unfinished Lesson is not completed.
- **FR-15.AC-5:** A completed Lesson stores all criterion scores and the overall score.
- **FR-15.AC-6:** A Lesson is stored in statistics at most once.

### FR-16 — Average Lesson scores

The Profile summarizes recent completed-Lesson performance.

- **FR-16.AC-1:** The Profile calculates averages over the latest `N` completed Lessons.
- **FR-16.AC-2:** If fewer than `N` completed Lessons exist, all available completed Lessons are used.
- **FR-16.AC-3:** If no completed Lessons exist, the Profile shows a no-data state rather than score `0`.
- **FR-16.AC-4:** Overall average uses the overall scores from the same Lessons.
- **FR-16.AC-5:** Each of the three criteria is averaged independently over the same Lessons.
- **FR-16.AC-6:** Level or Topic changes do not modify historical Lesson averages.

### FR-17 — Vocabulary progress

The Profile reports mastery of the current main dictionary.

- **FR-17.AC-1:** Progress includes only Translations in the current main dictionary.
- **FR-17.AC-2:** Each eligible Translation contributes at most `3` stars.
- **FR-17.AC-3:** Rating `-1` contributes `0` stars.
- **FR-17.AC-4:** `unseen` contributes `0` stars.
- **FR-17.AC-5:** Disabled-Topic, higher-Level, inactive-Topic, and inactive-Meaning content is excluded from both numerator and denominator.
- **FR-17.AC-6:** Progress is `sum(max(numeric rating, 0)) / (eligible Translation count × 3) × 100`.
- **FR-17.AC-7:** If the eligible Translation count is `0`, progress is displayed as `0%`.

### Epic H — Administrative catalog tools

### FR-20 — Vocabulary generation script

An administrative script can generate vocabulary rows for catalog experiments without regenerating an already populated Topic/Level set.

- **FR-20.AC-1:** The Admin lists Topics in the connected Sheet before generating vocabulary.
- **FR-20.AC-2:** The script accepts a Topic and Level as generation scope.
- **FR-20.AC-3:** Generated rows contain the Topic, Level, and exactly one English, one Russian, and one Greek Translation.
- **FR-20.AC-4:** Generated rows are written to the connected Sheet for Admin review; the application does not automatically approve or semantically review them.
- **FR-20.AC-5:** Before generation, the script checks whether the Sheet already contains one or more rows for the same `(Topic, Level)`.
- **FR-20.AC-6:** If any such row exists, the script skips generation for that entire `(Topic, Level)` pair.
- **FR-20.AC-7:** The script does not attempt semantic word-by-word duplicate detection.
- **FR-20.AC-8:** Re-running the script against an unchanged Sheet does not append another generated set for an already populated `(Topic, Level)` pair.

### FR-34 — Admin route access

Administrative routes do not reveal themselves to non-Admins.

- **FR-34.AC-1:** An Admin can open authorized Admin routes.
- **FR-34.AC-2:** An authenticated non-Admin receives HTTP `404` for an Admin route.
- **FR-34.AC-3:** A Guest also receives HTTP `404` for an Admin route rather than being redirected to sign-in.
- **FR-34.AC-4:** A rejected request does not disclose whether the Admin route exists.
- **FR-34.AC-5:** The first Admin account is created during system setup and cannot be created through public sign-up.

### FR-38 — Import and update catalog data from a Sheet

An Admin can validate and apply Sheet data as idempotent catalog upserts.

- **FR-38.AC-1:** Before applying any row, the import validates the complete input for missing required fields, invalid Topics, invalid Levels, missing required Translations, duplicate same-language Translations, and malformed IDs.
- **FR-38.AC-2:** If any row fails validation, the complete import fails before catalog changes are applied.
- **FR-38.AC-3:** The failure reports enough row-level information for the Admin to correct the Sheet.
- **FR-38.AC-4:** A row without a `meaning_id` represents a new Meaning.
- **FR-38.AC-5:** Creating a new Meaning assigns a stable system `meaning_id` and associates it with the source row so a rerun does not create a duplicate.
- **FR-38.AC-6:** A row containing a `meaning_id` that does not exist in the system is an error; the import does not create a Meaning with the supplied unknown ID.
- **FR-38.AC-7:** A row containing an existing `meaning_id` updates that Meaning's Topic, Level, active state, and Translations as supplied.
- **FR-38.AC-8:** Translation and other historical User ratings remain attached to their stable entity IDs when catalog data changes.
- **FR-38.AC-9:** Meaning identity follows `meaning_id`; v1 does not attempt semantic identity or version checks when an existing ID's content is changed.
- **FR-38.AC-10:** After successful full-input validation, rows may be applied incrementally rather than in one all-or-nothing transaction.
- **FR-38.AC-11:** If application stops partway through, already applied rows remain applied.
- **FR-38.AC-12:** Re-running the same import resumes safely and does not duplicate rows already applied.
- **FR-38.AC-13:** The final state after a successful rerun is the same as if the validated import had completed without interruption.
- **FR-38.AC-14:** If the system cannot access the connected Sheet, no import is applied and the Admin sees an access error.
- **FR-38.AC-15:** If application is interrupted after processing has begun, the Admin sees an interruption result and can safely run the import again.
- **FR-38.AC-16:** A successful import reports completion to the Admin.

## 9. LLM response contract requirements

Concrete schemas can be defined during implementation, but they must enforce the following product-level rules:

### Generation response

- Exactly the configured number of source sentences.
- One hidden suggested target-language translation per source sentence.
- Explicit used/not-used metadata for every supplied Translation.
- Stable references to supplied Translation and Meaning IDs; free-text word matching is insufficient.

### Grading response

- Integer criterion scores from `1` to `10` for Grammar, Spelling, and Vocabulary and meaning.
- Between `0` and `3` textual advice items.
- One correct/reference/alternative target-language sentence for every source sentence.
- One `+1`, `-1`, or `0` delta for every used Translation occurrence.
- One `+1`, `-1`, or `0` delta for every relevant Grammar Concept.
- Stable entity/occurrence references sufficient to apply updates deterministically and idempotently.

Any missing required field, invalid enum value, out-of-range score, unresolved entity reference, or inconsistent occurrence reference makes the response unusable and invokes FR-33.

## 10. Risks

| Risk | Required mitigation direction |
| --- | --- |
| LLM grading is inconsistent or misses errors | Use a fixed rubric per criterion and maintain representative grading examples for regression checks. |
| A valid alternative translation is marked wrong | Grade meaning and grammar rather than matching only one reference; the reference translation is not the sole acceptable answer. |
| LLM word-usage metadata is wrong | Use stable IDs in the response and spot-check usage results against real Greek inflections; do not substitute fragile string matching. |
| LLM cost grows with Lesson size or User count | Track cost per Lesson and keep sentence and vocabulary targets deployment-configurable. |
| An interrupted import leaves only part of a valid Sheet applied | Make per-row application idempotent and surface the interruption so a rerun converges to the intended final state. |
| Catalog edits change the semantics associated with an existing ID | Accept this limitation for experimental v1 data; add Meaning versioning before production-grade catalog governance is required. |

## 11. Out of scope for v1

- User-selectable native or target languages.
- User-facing management of additional languages.
- Grammar Concept ratings influencing generation.
- Semantic duplicate detection in the vocabulary generation script.
- Meaning versioning or automatic semantic identity protection when Admin data changes under an existing ID.
- Hard deletion of Meanings, Topics, ratings, or completed-Lesson history through the specified Admin workflows.

## 12. Open blocker

Only one product decision remains open:

- **FR-8 vocabulary selection algorithm:** define how new, mistaken, and mastered Translations are mixed when selecting vocabulary for a new Lesson, including behavior when a category contains fewer eligible Translations than its target share.

No implementation should infer or invent this rule from ratings, progress percentages, or the words-per-sentence configuration.

## 13. Out of scope of this document

This PRD defines product behavior and observable acceptance criteria. It does not define:

- technology stack, deployment, or local-development instructions;
- concrete storage technology or physical table layout;
- exact HTTP API schemas;
- exact LLM prompts or serialization formats beyond the product-level response requirements in Section 9;
- detailed visual design or per-screen layout;
- the exact Grammar Concept catalog for each Level, which is content supplied separately.

Repository implementation documents can define these details, but they must preserve the invariants and acceptance criteria in this PRD. Any scenario blocked by FR-8 must remain unimplemented until the vocabulary selection algorithm is explicitly approved.
