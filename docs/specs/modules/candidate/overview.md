# Candidate Module

## Purpose

The Candidate module manages public-user interview-preparation profile data. It stores stable, user-editable business information that is specific to the interview-preparation product and should not live in the Auth module.

## Responsibilities

- Manage candidate profile data for public users
- Store selected interview target-role and experience-level keys for the candidate
- Store user-editable presentation/profile fields such as full name, bio, and location
- Store preparation preferences and goals that are part of the product domain
- Store selected preferred topic keys as normalized preference records, not JSON blobs
- Track onboarding completion after the candidate selects role, experience level, and topics
- Own candidate-facing profile concepts while delegating binary files to Filevault
- Expose candidate data to other modules through a portal when needed
- Support gradual onboarding, where a candidate profile may start minimal at registration and be completed later

## Domain Main Entities

| Entity | Description |
| ------ | ----------- |
| `CandidateProfile` | One-to-one interview-preparation profile for a public user |
| `CandidateTopicPreference` | Preferred topic list for a candidate, ordered by priority |

## Boundary

- `auth` owns identity, credentials, sessions, and permissions
- `candidate` owns stable interview-preparation profile data, onboarding state, and selected catalog keys
- `interview` owns the target-role, experience-level, and topic catalogs used by onboarding and interview generation
- `analytics` owns derived metrics such as streaks, total interviews, average scores, achievements, and topic performance
- `filevault` owns profile photo/avatar storage; `candidate` only defines the business association for those files

See ERD.md for entity relationships.
