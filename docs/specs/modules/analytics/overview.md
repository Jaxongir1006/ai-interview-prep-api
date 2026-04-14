# Analytics Module

## Purpose

The Analytics module tracks derived candidate progress over time. It stores aggregated performance metrics, topic-level proficiency, and earned achievements based on completed interview activity.

## Responsibilities

- Store candidate progress summaries such as streaks, total interviews, average scores, and total time spent
- Store topic-level proficiency derived from interview results
- Store achievement definitions and awarded achievements
- Expose analytics data for dashboards, recommendations, and progress reports
- Keep derived performance data separate from editable profile data
- Provide authenticated candidate dashboard read models by combining analytics aggregates with candidate profile data and interview activity

## Domain Main Entities

| Entity | Description |
| ------ | ----------- |
| `CandidateProgressSummary` | One-to-one derived summary of a candidate's overall progress |
| `CandidateTopicStat` | Per-topic aggregated performance for a candidate |
| `AchievementDefinition` | Admin-defined achievement catalog entry |
| `CandidateAchievement` | Achievement earned by a candidate |

## Boundary

- `candidate` owns stable profile data and topic preferences
- `interview` owns interview sessions, submitted answers, and raw outcomes
- `analytics` owns derived aggregates and badges built from interview outcomes
- Strengths and weaknesses should be inferred from topic stats, not stored as free-form profile fields
- Dashboard read models may aggregate data through portals, but analytics must not take ownership of editable profile fields or raw interview answers

See ERD.md for entity relationships.
