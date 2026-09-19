# User scenarios: Learn the Grammar

**Purpose:** This document complements the PRD. It describes the visible journeys, screens, decisions, and recovery paths that designers and UI implementation agents need. The PRD remains the source of truth for domain rules, calculations, LLM contracts, and catalog/import behavior.

## How to use this document

Use each scenario to design the stated screens and their visible states. Do not introduce behavior that conflicts with the linked functional requirements.

This document intentionally does **not** repeat:

- rating formulas, clamps, or how multiple word occurrences are aggregated;
- LLM prompts, response schemas, retries, or usage-detection rules;
- deployment configuration values and import implementation details;
- the vocabulary-selection algorithm in `FR-8`, which remains blocked in the PRD.

## Roles

| Role | Visible access |
| --- | --- |
| Guest | Can use the interface-language selector and sign-in, sign-up, and password-recovery screens. Cannot enter lessons. |
| User | Can use Settings, Start Lesson, Lesson, Lesson Result, Dictionary, and Profile. |
| Admin | Has all User access and can open the Admin Console. |

The first Admin account is created during system setup. It cannot be created through public sign-up.

## Scenario 0 — Select interface language

**Role:** Guest, User, or Admin

**Design surfaces:** Public access screens, application navigation, Admin Console

### Main path

1. On a first visit, the app selects the first supported language from the browser preference order, or English when there is no supported preference.
2. A Guest can open the locale selector on sign-in, sign-up, or password-recovery screens and select Russian, Greek, or English.
3. The interface changes immediately.
4. An authenticated User or Admin can select a locale from application navigation; the selection also applies immediately in the Admin Console and is retained for later sign-ins.

### Result

The visitor can use the interface in a supported language without changing the fixed Russian-to-Greek learning pair.

**Requirements:** FR-39

## Scenario 1 — Sign up and first-time setup

**Role:** Guest who becomes a User  
**Design surfaces:** Sign up, Settings, Start Lesson

### Main path

1. The Guest opens the app and chooses **Sign up**.
2. The Guest enters an email address and password and submits the form.
3. The account is created and the new User is signed in automatically.
4. The app opens **Settings**. `A1` is selected and all currently active Topics are enabled.
5. The User may change the Level and enabled Topics.
6. The User selects **Save**.
7. The app opens **Start Lesson**.

### Visible alternatives and errors

- Invalid registration details show a validation error and keep the User on the sign-up screen.
- If the User leaves Settings before the first successful Save, the next sign-in returns to Settings.
- If the User tries to start a Lesson but has no eligible vocabulary, the app explains that a Lesson cannot be started and offers navigation to Settings.

### Result

The User has completed onboarding and has saved learning settings.

**Requirements:** FR-1, FR-2, FR-26, FR-29, FR-30, FR-31

## Scenario 2 — Sign in and recover access

**Role:** Guest or returning User  
**Design surfaces:** Sign in, Password recovery, Set new password, Settings or Start Lesson

### Sign in

1. The Guest opens **Sign in** and enters email and password.
2. After successful sign-in, the app opens Settings if onboarding is incomplete; otherwise it opens Start Lesson.

### Password recovery

1. The Guest selects **Forgot password**.
2. The Guest enters the account email and requests a recovery link.
3. After opening a valid recovery link, the User enters and confirms a new password.
4. The app confirms that access has been restored and allows the User to sign in with the new password.

### Visible alternatives and errors

- Incorrect sign-in credentials show an error and do not sign the person in.
- An invalid, expired, or already-used recovery link shows an error and offers the option to request a new link.
- A new password that does not meet validation rules is rejected on the reset screen.

### Result

The returning User either has an authenticated session or has a clear next action to recover access.

**Requirements:** FR-1, FR-26, FR-36, FR-37

## Scenario 3 — Start and complete a Lesson

**Role:** User  
**Design surfaces:** Start Lesson, loading state, Lesson, Hint, Lesson Result

### Main path

1. The User opens **Start Lesson** and selects **Start lesson**.
2. The screen shows a loading state while the Lesson is prepared.
3. The **Lesson** screen shows the generated Russian source sentences and one Greek answer field for each sentence.
4. The User enters any number of answers. Empty fields remain valid answers to submit.
5. Optional: the User selects a source word and requests a hint. The app shows a target-language dictionary/base form, not a completed sentence.
6. The User selects **Submit**.
7. The screen shows a loading state while answers are graded.
8. The **Lesson Result** screen shows Grammar, Spelling, and Vocabulary and meaning scores; the overall score; zero to three advice items; and a reference or alternative Greek sentence for every source sentence.

### Visible alternatives and errors

- If preparation fails, the app shows `Oops, something went wrong` and lets the User try starting again.
- If grading fails, the app shows the same error and returns the User to the Lesson with entered answers intact, ready to submit again.
- While preparation or grading is running, the relevant action is unavailable to prevent a second request.

### Result

The User can review the completed Lesson. Dictionary and Profile views reflect the completed result.

**Requirements:** FR-9, FR-10, FR-11, FR-12, FR-13, FR-14, FR-15, FR-25, FR-33

## Scenario 4 — Resume or discard an unfinished Lesson

**Role:** User  
**Design surfaces:** Start Lesson, Resume choice, Lesson

### Main path

1. The User leaves a Lesson before it has a result, then later opens **Start Lesson**.
2. The screen explains that an unfinished Lesson exists and offers **Continue** and **End and start new**.
3. If the User chooses **Continue**, the same source sentences and the latest saved answers open in the Lesson screen.
4. If the User chooses **End and start new**, the unfinished Lesson is discarded without a result, and the app begins a new Lesson.

### Visible alternatives and errors

- Changing the current Level or Topic settings in the meantime does not change the unfinished Lesson the User resumes.
- If grading previously failed, the User can return to the same saved answers and submit again.

### Result

The User either resumes the same unfinished work or intentionally starts fresh. Discarding does not produce a Lesson Result.

**Requirements:** FR-32, FR-33

## Scenario 5 — Change learning settings

**Role:** User  
**Design surfaces:** Settings, Dictionary, Start Lesson, Profile

### Main path

1. The User opens **Settings**.
2. The User selects `A1`, `A2`, or `B1` and enables or disables active Topics.
3. The User selects **Save**.
4. The app confirms the saved settings. The Dictionary and future Lessons use the updated selection.

### Visible alternatives and errors

- Leaving without Save does not apply the changes.
- Existing ratings and completed-Lesson history remain available after changing Level or Topics.
- An unfinished Lesson remains unchanged and can still be resumed from Start Lesson.

### Result

Future lesson content and the current main dictionary match the saved settings.

**Requirements:** FR-1, FR-2, FR-4, FR-6, FR-31, FR-32

## Scenario 6 — Browse the Dictionary

**Role:** User  
**Design surfaces:** Dictionary tabs or grouped sections

### Main path

1. The User opens **Dictionary**.
2. The User can browse the current main dictionary, vocabulary from disabled Topics, and vocabulary above the current Level in separate sections.
3. Each visible Translation shows either `unseen` or its personal numeric rating.
4. The User changes Settings and returns to Dictionary; the groups update to match the saved Level and Topic choices.

### Visible alternatives and errors

- `unseen` is visually distinct from numeric `0`.
- Inactive Meanings and inactive Topics are not shown.

### Result

The User can understand which vocabulary is currently available for learning and which vocabulary is outside the current main dictionary.

**Requirements:** FR-3, FR-4, FR-5, FR-7, FR-18, FR-22, FR-23, FR-28

## Scenario 7 — View progress in the Profile

**Role:** User  
**Design surfaces:** Profile

### Main path

1. The User opens **Profile**.
2. The Profile shows the recent overall average and separate averages for Grammar, Spelling, and Vocabulary and meaning.
3. The Profile shows vocabulary progress for the current main dictionary.
4. The Profile shows the User's ratings for Grammar Concepts unlocked at the current Level.
5. The Profile also shows the current learning settings.

### Visible alternatives and empty states

- Before the User completes a Lesson, score areas show a no-data state rather than a score of zero.
- If the current main dictionary is empty, vocabulary progress is shown as `0%`.
- Grammar Concepts above the current Level are not shown.

### Result

The User can see current progress without mistaking unavailable or not-yet-collected data for poor performance.

**Requirements:** FR-16, FR-17, FR-27, FR-31, FR-35

## Scenario 8 — Update learning content

**Role:** Admin  
**Design surfaces:** Connected `Topics` and `Meanings` tables, vocabulary-generation binary, Admin Console, import status

### Main path

1. The Admin adds a Topic to the `Topics` table with `active=true` and `filled=false`. The Admin does not add word rows manually.
2. The Admin runs the vocabulary-generation binary outside the application. In one run, it uses the LLM to generate words for every active, unfilled Topic at every v1 Level and writes them to the `Meanings` table.
3. After generating all Levels for a Topic successfully, the binary marks that Topic `filled=true`. Topics already filled or inactive are skipped on later runs.
4. The Admin opens the **Admin Console** and selects **Update words from connected sheet**.
5. The console shows that the import completed successfully.
6. New or updated active content becomes available to Users according to their settings.

### Manage content through the Sheet

- The Admin can add an active, unfilled Topic; it becomes available after vocabulary generation and import.
- The Admin can set a Meaning row inactive; after the next successful import it is no longer shown to Users or used in new Lessons.
- The Admin can reactivate an inactive Meaning by setting its `active` flag to `true` and importing again.
- The Admin can set a Topic inactive; after import, that Topic and its words are no longer available to Users or new Lessons.

### Visible alternatives and errors

- A Guest or non-Admin who attempts to open an Admin route sees a not-found page.
- If LLM generation fails for a Topic or Level, the binary leaves the Topic unfilled, reports the failure, and the Admin can run it again.
- If the Sheet cannot be accessed, the console shows an access error and the Admin can correct access and try again.
- If the Sheet contains invalid data, the console reports the rows that must be corrected before the import can proceed.
- If an import is interrupted, the console reports the interruption and the Admin can run it again.

### Result

The Admin has a clear, recoverable path to make reviewed catalog content available without exposing administrative controls to other users.

**Requirements:** FR-20, FR-21, FR-23, FR-28, FR-34, FR-38

## Design coverage checklist

Before a design or UI implementation is considered complete, it must include:

- primary, loading, validation-error, empty, and recovery states named in the relevant scenario;
- explicit first-time Settings Save and unfinished-Lesson choice states;
- a clear distinction between `unseen`, numeric ratings, no-data scores, and `0%` vocabulary progress;
- all application-owned visible text available in each supported interface locale, while learning-pair and LLM content remains unchanged by the interface locale;
- user-safe error states that retain entered Lesson answers after grading failure;
- no visible Admin navigation or route disclosure for Guests and non-Admins.

For rules that are not visible in a particular journey, consult the PRD rather than extending this document.
