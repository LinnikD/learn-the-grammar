# PRD: Language Grammar & Writing Trainer

## Overview

A web application that trains grammar, spelling and word choice through translation exercises. The LLM writes sentences in the user's native language, the user translates them into the language being learned, and the LLM grades the translations. The app keeps a personal dictionary with a rating for every word and uses it to choose the words for the next lessons.

The first version is a personal project built for its author, who is learning Greek from Russian. It targets learners at levels A1 to B1 who want to practice writing, grammar and spelling rather than passive recognition.

A lesson has three steps:

1. The LLM writes several sentences in the user's native language.
2. The user translates them into the target language and submits.
3. The LLM checks the translations, scores them and gives feedback.

## Goals, scope and non-goals

**Goals**

- G-1: Improve the user's active grammar, spelling and word choice through regular translation practice.
- G-2: Adapt lessons to what the user already knows, so practice focuses on weak words.
- G-3: Show progress: recent scores and how much of the dictionary is learned.

**First version**

- Language pair: Russian (native) to Greek (target). English is stored for every word, so it can be added quickly as a native or target language.
- Levels: A1, A2, B1.
- Platform: one site with two parts: the user app, which works on desktop and in a mobile browser, and the admin console, which only the admin can open.
- Access: only logged-in users can take a lesson.
- Postponed: letting the user choose N for the average score; access without login for easier onboarding; adding a language through the admin console.

**Non-goals**

- NG-1: Listening, speaking and pronunciation practice.
- NG-2: Levels above B1.
- NG-3: Free-form conversation with the LLM.

## Glossary

| Term | Meaning |
| --- | --- |
| Lesson | One round: several sentences to translate, then grading. |
| Level | A1, A2 or B1. The user picks one at setup and can change it later. |
| Topic | A theme such as family or travel, defined by the admin. |
| Meaning | The unit shared by all users: one meaning of a word, with a topic, a level and its translations. The same English word can have several meanings. |
| Translation | One word in one language for a meaning. English, Russian and Greek are stored for every meaning. |
| Dictionary | A user's own list of words in the language being learned, each with the user's own rating. |
| Rating | A user's number for a word, from -1 to +3, shown as up to 3 stars. |
| Grammar concept | One grammar construction, tied to a language and a level. Each user has a personal rating per concept. |
| Admin | A user account with admin rights that manages topics and words. |
| User app | The part of the site where users take lessons, browse their dictionary and see their profile. Works on desktop and in a mobile browser. |
| Admin console | The part of the site available only to the admin, where topics and words are managed. |
| Main dictionary | The part of the user's dictionary that lessons use: words of the user's level and lower levels from enabled topics. |

## Screens

A skeleton only: which screens exist and what each is for. Layout and details are agreed with the agent before each feature. User scenarios are described in a separate document (see Out of scope of this document).

**User app**

| Screen | Purpose |
| --- | --- |
| Login and sign up | Sign in; sign up with an email and a password, without email confirmation. Signing in is required before taking a lesson. |
| Settings | Choose the language level and enable or disable topics. |
| Start lesson | Starts a new lesson. If an unfinished lesson exists, offers two actions: continue it, or end it without a result and start a new one. |
| Lesson | Shows the sentences, a field for each translation, a hint request for a word, and submit. |
| Lesson result | Scores per criterion and overall, advice, and a correct or alternative version of every sentence. |
| Dictionary | The user's words with ratings, in tabs: main dictionary, disabled topics, higher levels. |
| Profile | Average score over the last N lessons with a breakdown, vocabulary progress, the rating of each grammar concept, and current settings. |

**Admin console**

| Screen | Purpose |
| --- | --- |
| Topics | Add a topic; delete a topic together with its meanings. |
| Word generation | Generate meanings with translations for a topic and a level. |
| Meanings | Browse meanings with their translations and topic; ban a meaning. |
| Languages (postponed) | Add a language: the LLM produces its translations for the existing meanings. |

## Levels, topics and grammar

The only required setting is the language level.

- **Levels:** A1, A2, B1. Each level adds words and grammar on top of the lower ones. The user can change the level later: the main dictionary and the unlocked grammar change with it, while word ratings and the average score of the last lessons are kept.
- **Topics:** all topics are enabled by default, and the user can turn off the ones they are not interested in. Disabled topics are not used in lessons. The ratings of their words are kept, so they are still there when the topic is enabled again, and disabling a topic does not change the average score of the last lessons.
- **Words per topic and level:** each topic has its own word set for every level. A level can have no words for a topic: politics is too complex for A1, so it stays empty there.
- **Grammar per level:** each level unlocks new constructions, for example tenses, comparative forms of adjectives and compound sentences. Unlocking is cumulative: grammar from lower levels stays available. The exact list per level is TBD. Each construction is a grammar concept, and the same concepts feed the tips shown to the user and the instructions given to the LLM.
- **Lesson content:** a lesson uses words of the user's level and lower levels from enabled topics, and all grammar unlocked so far.

## Dictionary and ratings

- **Main dictionary:** the words of the user's level and of all lower levels from enabled topics. Lessons use only these words.
- **Separate tabs:** words of disabled topics and words of higher levels are shown in their own tabs, so the user can still browse them.
- **Rating:** each word has a rating from -1 to +3, shown as up to 3 stars. Every word starts at 0. A correct answer adds 1 and a mistake subtracts 1. Negative values mark words the user does not know, so they can be given more often.
- **Not seen yet:** words that have not appeared in any lesson are shown in the dictionary without a rating.
- **Which words change:** only words that appear in the lesson's sentences change rating after the lesson; skipped words keep their rating.
- **Mastery:** a word is mastered when its rating reaches +3. There is no decay over time, and mastered words can still appear in lessons. This rule is a draft.
- **Choosing words for a lesson:** each lesson mixes three kinds of words: new ones, ones the user made mistakes in, and mastered ones. The exact selection algorithm is TBD.
- **Grammar concepts:** each concept has its own rating per user, separate from word ratings. A concept's rating goes up by 1 when it is used correctly and down by 1 when it is used with a mistake, and it ranges from -1 to 10 (up to 10 stars). In the first version it is used for statistics only and does not affect lesson generation.

## Lesson flow

- **Number of sentences:** 5 at the start; a configurable value, not chosen by the user in the first version.
- **What the LLM receives:** the user's level, the words selected by rating, and the grammar unlocked for the level. The number of words depends on the level: as a starting point 4 words per sentence at A1, 7 at A2 and 10 at B1, which for 5 sentences is 20, 35 and 50 words. If the main dictionary has fewer words than that, all of them are sent, and the prompt says that words outside the list may be used. If there are no words at all, the lesson does not start and the user is offered to set up topics.
- **What the LLM does:** tries to write the sentences using these words. It may use other words and may skip some of the given ones, but then it lists the skipped words. It also returns a suggested translation of each sentence, hidden and used only for hints.
- **User:** enters translations of all sentences and submits them for checking.
- **Unfinished lesson:** an unfinished lesson is saved. On the start lesson screen the user can continue it, or end it without a result and start a new one. A lesson ended without a result changes no ratings and no scores.
- **LLM errors:** if the LLM does not respond or returns an unusable answer, the user sees an error message in the interface ("Oops, something went wrong").
- **Hints:** shown only when the user asks. A hint gives the translation of one word chosen by the user, not of the whole sentence. Using a hint does not affect the score.
- **Request format:** the exact format of requests and responses is not defined here (see Out of scope of this document).

## Grading and feedback

The LLM grades the lesson on three criteria:

- **Grammar:** whether forms, tenses, agreement and word order are correct.
- **Spelling:** whether words are written correctly, including accents.
- **Vocabulary and meaning:** whether the words fit the context and the translation keeps the meaning of the source sentence.

Each criterion gets an integer score from 1 to 10. The overall lesson score is the average of the criterion scores, from 1 to 10 with two decimals (for example 7.55). The exact grading prompt is TBD.

The user sees the result of the lesson:

- the score for each criterion and the overall score;
- up to 3 short pieces of advice from the LLM (character limit TBD);
- for every sentence, a correct or alternative version, even when the user's sentence is right.

The LLM also marks each grammar concept of the user's level as correct, mistaken or not used. A concept marked not used keeps its rating.

## Profile and statistics

- **Average score:** mean lesson score over the last N lessons, with a breakdown by criterion. N is 5 at the start and configurable; later the user can choose the range.
- **Vocabulary progress:** the share of stars earned in the main dictionary. Each word can give up to 3 stars, words with rating -1 count as 0 stars, and words of disabled topics and higher levels are not counted.
- **Grammar progress:** the rating of each grammar concept.
- **Settings:** language pair, level and enabled topics.

```latex
mastered\_pct = \frac{star\_count}{total\_dictionary\_words \times 3} \times 100
```

## Content administration

All content is managed in the admin console; users cannot add words. Admin rights belong to specific accounts, and the first admin is created when the system is set up, not through the app.

- **Topics:** the admin adds topics, which must exist before word generation runs, and can add new ones at any time.
- **Word generation:** for each topic and level, the LLM generates a set of meanings with translations in English, Russian and Greek. Generated words are not reviewed before use.
- **Ban:** the admin can ban a meaning. A banned meaning is removed from every user's dictionary and is never used in lessons or statistics.
- **Delete a topic:** deletes all its meanings and removes them from every user's dictionary.
- **New language (postponed):** the LLM produces the new translation for the existing meanings. Existing ratings are not affected.

## Requirements

Requirements are grouped into epics. IDs are stable and do not follow the order within an epic. Each epic names the sections with its rules. Epic H (content) must provide topics and words before lessons in Epics D to G can work, and Epic A sets the level and topics that Epics B to G use.

### Epic A: Setup and access

Rules: Levels, topics and grammar.

| ID | Requirement |
| --- | --- |
| FR-26 | Only logged-in users can take a lesson. |
| FR-1 | The user selects a level during setup; all topics are enabled by default. |
| FR-2 | The user can disable and re-enable topics. |
| FR-18 | The first version supports the pair Russian to Greek. |
| FR-29 | A new user can sign up with an email and a password; email confirmation is not required. |
| FR-31 | The user can change the level later; word ratings and the average score of the last lessons are kept. |

### Epic B: Dictionary

Rules: Dictionary and ratings; Levels, topics and grammar.

| ID | Requirement |
| --- | --- |
| FR-3 | Words are organized by topic, with a word set per level; a topic may have no words at a given level. |
| FR-4 | The user's dictionary includes words of their level and of all lower levels. |
| FR-22 | Each user has a personal dictionary with their own rating per word. |
| FR-7 | A word's rating is from -1 to +3, starts at 0, and +3 means mastered. |
| FR-5 | The user can browse the dictionary; unseen words have no rating; words of disabled topics and higher levels are in separate tabs. |
| FR-19 | English or another language can be added as a native or target language without regenerating existing words. |

### Epic C: Grammar concepts

Rules: Levels, topics and grammar (grammar per level); Dictionary and ratings (grammar concepts).

| ID | Requirement |
| --- | --- |
| FR-6 | Each level unlocks grammar constructions, cumulatively. |
| FR-24 | Grammar concepts belong to a language and a level; each user has a personal rating per concept. |

### Epic D: Lesson generation

Rules: Lesson flow.

| ID | Requirement |
| --- | --- |
| FR-8 | For each lesson the LLM gets the user's level, the words selected by rating and the allowed grammar. |
| FR-9 | The LLM returns sentences in the native language, each with a hidden suggested translation, and lists the words it skipped. |
| FR-27 | The number of sentences per lesson (default 5) and N for the average score (default 5) are configurable values. |
| FR-30 | If the main dictionary has no words, the lesson does not start and the user is offered to set up topics. |
| FR-33 | If the LLM does not respond or returns an unusable answer, the user sees an error message in the interface. |

### Epic E: Taking a lesson

Rules: Lesson flow.

| ID | Requirement |
| --- | --- |
| FR-11 | The user can enter translations of all sentences and submit them. |
| FR-10 | The user can request a hint for a word of their choice; hints do not affect the score. |
| FR-32 | An unfinished lesson is saved; the user can continue it, or end it without a result and start a new one. |

### Epic F: Grading and feedback

Rules: Grading and feedback.

| ID | Requirement |
| --- | --- |
| FR-12 | The LLM gives an integer score from 1 to 10 for each criterion; the user sees the result of the lesson. |
| FR-13 | The user sees the score, up to 3 pieces of advice and a correct or alternative version of every sentence. |
| FR-25 | The LLM marks each grammar concept as correct, mistaken or not used, and the user's concept ratings and statistics are updated from it. |

### Epic G: Ratings and progress

Rules: Dictionary and ratings; Profile and statistics.

| ID | Requirement |
| --- | --- |
| FR-14 | Word ratings are updated from the graded answers, only for words that appear in the sentences. |
| FR-15 | The score of every completed lesson is stored. |
| FR-16 | The profile shows the average score over the last N lessons, with a breakdown by criterion. |
| FR-17 | The profile shows the share of stars earned in the main dictionary. |
| FR-35 | The profile shows the user's rating for each grammar concept. |

### Epic H: Content administration

Rules: Content administration.

| ID | Requirement |
| --- | --- |
| FR-20 | The admin defines topics and has the LLM generate words with translations for each topic and level. |
| FR-21 | Each generated meaning is saved with its topic, level and translations in English, Russian and Greek. |
| FR-23 | The admin can ban a meaning; it is removed from every user's dictionary. |
| FR-28 | The admin can add topics at any time and delete a whole topic together with its meanings. |
| FR-34 | The admin console is available only to accounts with admin rights; the first admin is created when the system is set up. |

## Risks

| Risk | Mitigation idea |
| --- | --- |
| The LLM grades inconsistently or misses errors, so scores feel unfair | Fixed rubric per criterion, spot checks on sample answers |
| A correct alternative translation is marked wrong | Grade meaning and grammar, not match to one reference; accept valid variants |
| Mastery rule marks words as known too early or too late | Make the rule configurable and review it against real usage |
| LLM cost per lesson grows with users | Track cost per lesson early; cap sentences per lesson |
| The LLM's list of skipped words is wrong, so ratings of the wrong words change (Greek word forms make text matching hard) | Spot-check the skipped-word lists on real lessons early |

## Open questions

- [ ] What grammar does each level unlock?
- [ ] How does the app choose the words for a lesson: how many new, mistaken and mastered words, and in what order? (algorithm TBD)

## Out of scope of this document

This document describes what the product does, not how it is built.

- **Tech stack, deployment and local development:** the repository README and AGENTS.md.
- **API contract:** `api/openapi.yaml`.
- **LLM exchange format** (requests, responses and prompts for lesson generation, grading and hints): TODO, to be designed in a separate file (path TBD). Do not implement this part without an explicit protocol.
- **Storage:** technology and data layout are not part of this document; the entities in the glossary are logical only.
- **User scenarios:** TODO, described in a separate document (path TBD).
