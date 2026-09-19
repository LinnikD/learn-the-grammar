# User scenarios: Language Grammar & Writing Trainer

## Roles

There are three roles. A person moves from one to the next by signing up and by being given admin rights.

| Role | Who | What they can do |
| --- | --- | --- |
| Guest | Someone who is not logged in. | Open the login and sign up screen. Cannot take a lesson. |
| User | A logged-in account. | Use the user app: settings, start lesson, lesson, lesson result, dictionary and profile. |
| Admin | A logged-in account with admin rights. | Everything a User can do, plus the admin console. |

The first admin is created when the system is set up, not through the app.

## How to read this document

These scenarios describe what a person sees and does, step by step, using the roles above. They cover only what is visible to the person; the internal logic (how words are chosen, what the LLM receives, how ratings change) is in the PRD and is not repeated here. Each scenario lists the requirements (FR) it covers.

Each scenario has the same parts:

- **Role:** who does it.
- **Before:** the conditions that must be true.
- **Steps:** the main path.
- **Alternatives and errors:** what happens when the main path is not followed.
- **Result:** what is true afterwards.
- **Requirements:** the FR IDs from the PRD.

Where the PRD does not say what should happen, the scenario says so and points to the list of gaps at the end. Gaps are questions for the author, not decisions.

## Scenario 1: First launch

**Role:** a Guest who becomes a User.

**Before:** the person has no account. The admin has already added topics and generated words (see Scenario 4).

**Steps**

1. The Guest opens the site and sees the login and sign up screen.
2. The Guest chooses to sign up and enters an email and a password. No email confirmation is needed.
3. The account is created and the person is signed in automatically.
4. After the first sign in, which happens right after sign up, the person is taken to the Settings screen to choose the language level (A1, A2 or B1). A1 is selected by default. All topics are enabled by default, and the User can turn off the ones they do not want.
5. The User confirms the settings and goes to the Start lesson screen.
6. The User starts the first lesson (Scenario 2).

**Alternatives and errors**

- **No words available.** If the main dictionary has no words (for example all topics are turned off, or the admin has not generated words for the chosen level), the lesson does not start. The User is offered to set up topics and returns to Settings (FR-30).
- **Email already registered, invalid email, weak password.** Sign up fails with an error message. The only password rule is a length of more than 8 characters.
- **The User skips Settings.** The default level is A1, so a level is always set; skipping Settings does not block a lesson.

**Result:** a signed-in User with a level and a set of enabled topics, ready to start a lesson.

**Requirements:** FR-29, FR-26, FR-1, FR-2, FR-30.

## Scenario 2: Regular lesson

**Role:** User.

**Before:** the User is signed in, the main dictionary has words, and there is no unfinished lesson (or the User chose to start a new one, see Scenario 3).

**Steps**

1. On the Start lesson screen the User starts a new lesson. A loader blocks the screen until the sentences are ready.
2. The Lesson screen shows 5 sentences in the User's native language, with a field for each translation.
3. The User translates each sentence into the target language.
4. Optional: the User asks for a hint on one word of their choice in a sentence and sees its translation. Hints never change the score.
5. The User submits all translations. A loader blocks the screen until the result is ready.
6. The Lesson result screen shows the score for each criterion (grammar, spelling, vocabulary and meaning), the overall score, up to 3 short pieces of advice, and a correct or alternative version of every sentence, even those that were right.
7. The dictionary and the profile reflect the result of the lesson.

**Alternatives and errors**

- **No words available.** The lesson does not start; see Scenario 1.
- **LLM fails while writing the sentences.** The User sees the error message "Oops, something went wrong" (FR-33). The User can press the start button again.
- **LLM fails while checking the translations.** The User sees the same error message. The form is not reset: the typed translations stay, and the User can submit again.
- **The User leaves before submitting.** See Scenario 3.

**Result:** the lesson is finished, and the dictionary and the profile reflect its result.

**Requirements:** FR-11, FR-10, FR-12, FR-13, FR-33.

## Scenario 3: Unfinished lesson

**Role:** User.

**Before:** the User started a lesson and left before submitting the translations, for example by closing the browser or opening another screen.

**Steps**

1. The sentences and the translations typed so far are saved. How saving works is decided at implementation.
2. Later the User opens the Start lesson screen.
3. The screen shows that an unfinished lesson exists and offers two actions: continue it, or end it without a result and start a new one.
4. **Continue.** The Lesson screen opens with the same sentences and the translations typed so far. The User goes on from step 3 of Scenario 2.
5. **End and start a new one.** The old lesson is closed without a result. It changes no ratings and no scores. A new lesson starts from step 1 of Scenario 2.

**Alternatives and errors**

- **Settings changed in between.** If the User changes the level or turns a topic off while a lesson is unfinished, that lesson does not change.

**Result:** either the old lesson continues, or it is closed with no effect and a new one begins.

**Requirements:** FR-32.

## Scenario 4: Admin fills the content

**Role:** Admin.

**Before:** the system is set up and the first admin exists. The Admin is signed in and opens the admin console. Only accounts with admin rights can open it.

**Steps**

1. The Admin creates a Google Sheet that the system can read and lists the topics on the topics sheet.
2. The Admin runs the word generation script, outside the app. It fills another sheet with meanings and translations in English, Russian and Greek for each topic and level.
3. The Admin reviews the result in the sheet. To ban a wrong meaning, the Admin edits the sheet.
4. In the admin console the Admin presses "Update words from connected sheet". All meanings from the sheet are loaded into the system, and the ids given to new meanings are written back into the sheet.
5. Users can now start lessons with the new words.

**Alternatives and errors**

- **A non-admin opens the admin console.** It is available only to accounts with admin rights; everyone else gets a 404 page.
- **Loading the sheet again.** A row with an id updates the existing meaning and keeps the users' ratings; a row without an id is added as a new meaning.
- **The script is run again for a topic and level that already have words.** Duplicates are skipped; the script must be idempotent.
- **The sheet cannot be read, or loading fails.** If the system has no access to the sheet, the Admin sees an error. If loading stops partway, the Admin sees a message that it was interrupted. Loading can simply be started again: it is idempotent, so nothing is duplicated.
- **Banning a meaning or deleting a topic.** Both are done through the sheet. How they are marked there is an open question in the PRD.

**Result:** topics and meanings exist, so Users can start lessons.

**Requirements:** FR-34, FR-20, FR-21, FR-38, FR-23, FR-28.

## Scenarios not yet described

These scenarios are TBD. Do not implement them without an explicit discussion with the author.

| Scenario | What it would cover | Requirements |
| --- | --- | --- |
| Sign in | A returning User signs in; a wrong password shows an error; password recovery through a link sent to the email. | FR-26, FR-36, FR-37 |
| Change settings | The User changes the level or turns topics on and off; what this does to the dictionary and the ratings. | FR-1, FR-2, FR-31 |
| Browse the dictionary | The User opens the dictionary and its tabs (main dictionary, disabled topics, higher levels); words without a rating. | FR-3, FR-4, FR-5, FR-22 |
| View the profile | The User sees the average score of the last N lessons with a breakdown, vocabulary progress and grammar concept ratings. | FR-16, FR-17, FR-35 |
| Ban a meaning or delete a topic | The Admin removes a meaning or a topic through the sheet; how it is marked there is an open question in the PRD. | FR-23, FR-28 |
| Add a language | Postponed. A new translation is produced for the existing meanings. | FR-19 |

## Gaps found while writing

There are no open gaps at the moment. Questions that come up while writing more scenarios are listed here.
