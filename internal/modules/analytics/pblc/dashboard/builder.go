package dashboard

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/progress"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/candidate"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/filevault"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/interview"

	"github.com/code19m/errx"
)

const (
	CodeDashboardUnavailable = "DASHBOARD_UNAVAILABLE"

	Range7D  = "7d"
	Range30D = "30d"
	Range90D = "90d"
	RangeAll = "all"
)

type Builder struct {
	domainContainer *domain.Container
	portalContainer *portal.Container
}

func NewBuilder(domainContainer *domain.Container, portalContainer *portal.Container) *Builder {
	return &Builder{
		domainContainer: domainContainer,
		portalContainer: portalContainer,
	}
}

type Option struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID              string  `json:"id"`
	FullName        *string `json:"full_name"`
	Email           *string `json:"email"`
	AvatarURL       *string `json:"avatar_url"`
	TargetRole      *Option `json:"target_role"`
	ExperienceLevel *Option `json:"experience_level"`
}

type StatValue struct {
	Value          int    `json:"value"`
	DeltaPercent   int    `json:"delta_percent"`
	DeltaDirection string `json:"delta_direction"`
}

type NullableStatValue struct {
	Value          *int   `json:"value"`
	DeltaPercent   int    `json:"delta_percent"`
	DeltaDirection string `json:"delta_direction"`
}

type StreakValue struct {
	Value    int  `json:"value"`
	IsRecord bool `json:"is_record"`
}

type Stats struct {
	TotalInterviews      StatValue         `json:"total_interviews"`
	AverageScore         NullableStatValue `json:"average_score"`
	CurrentStreakDays    StreakValue       `json:"current_streak_days"`
	TotalPracticeSeconds StatValue         `json:"total_practice_seconds"`
}

type PerformanceSummary struct {
	AverageScore        *int  `json:"average_score"`
	ScoreDeltaPercent   int   `json:"score_delta_percent"`
	InterviewsCompleted int   `json:"interviews_completed"`
	PracticeSeconds     int64 `json:"practice_seconds"`
}

type PerformancePoint struct {
	Date                string `json:"date"`
	Label               string `json:"label"`
	AverageScore        *int   `json:"average_score"`
	InterviewsCompleted int    `json:"interviews_completed"`
	PracticeSeconds     int64  `json:"practice_seconds"`
}

type Performance struct {
	Range   string             `json:"range"`
	Summary PerformanceSummary `json:"summary"`
	Points  []PerformancePoint `json:"points"`
}

type TopicPerformance struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Score              *int     `json:"score"`
	QuestionsAnswered  int      `json:"questions_answered"`
	CorrectnessRate    *float64 `json:"correctness_rate"`
	AverageTimeSeconds *int64   `json:"average_time_seconds"`
	Trend              string   `json:"trend"`
	Level              string   `json:"level"`
}

type HighlightTopic struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Score             *int    `json:"score"`
	QuestionsAnswered int     `json:"questions_answered"`
	Reason            string  `json:"reason"`
	RecommendedAction *string `json:"recommended_action,omitempty"`
}

type Topics struct {
	Items  []TopicPerformance `json:"items"`
	Weak   []HighlightTopic   `json:"weak"`
	Strong []HighlightTopic   `json:"strong"`
}

type ActivityItem struct {
	SessionID       string   `json:"session_id"`
	Title           string   `json:"title"`
	Status          string   `json:"status"`
	Score           *int     `json:"score"`
	StartedAt       string   `json:"started_at"`
	CompletedAt     *string  `json:"completed_at"`
	DurationSeconds int64    `json:"duration_seconds"`
	QuestionCount   int      `json:"question_count"`
	AnsweredCount   int      `json:"answered_count"`
	Topics          []Option `json:"topics"`
}

type RecentActivity struct {
	Items      []ActivityItem `json:"items"`
	NextCursor *string        `json:"next_cursor"`
}

type RecommendedTopic struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Priority string `json:"priority"`
	Reason   string `json:"reason"`
}

type NextInterview struct {
	TargetRole               *Option  `json:"target_role"`
	ExperienceLevel          *Option  `json:"experience_level"`
	Topics                   []Option `json:"topics"`
	Difficulty               string   `json:"difficulty"`
	QuestionCount            int      `json:"question_count"`
	EstimatedDurationSeconds int64    `json:"estimated_duration_seconds"`
}

type Recommendations struct {
	RecommendedTopics []RecommendedTopic `json:"recommended_topics"`
	NextInterview     NextInterview      `json:"next_interview"`
}

type Overview struct {
	User            User            `json:"user"`
	Stats           Stats           `json:"stats"`
	Performance     Performance     `json:"performance"`
	Topics          Topics          `json:"topics"`
	RecentActivity  RecentActivity  `json:"recent_activity"`
	Recommendations Recommendations `json:"recommendations"`
}

func NormalizeRange(v string) string {
	if v == "" {
		return Range7D
	}
	return v
}

func IsValidRange(v string) bool {
	switch v {
	case Range7D, Range30D, Range90D, RangeAll:
		return true
	default:
		return false
	}
}

func (b *Builder) User(ctx context.Context, userCtx *auth.UserContext) (User, error) {
	u := User{
		ID:    userCtx.UserID,
		Email: userCtx.Email,
	}

	profile, hasProfile, err := b.profile(ctx, userCtx.UserID)
	if err != nil {
		return u, errx.Wrap(err)
	}
	if !hasProfile {
		return u, nil
	}

	u.FullName = profile.FullName
	u.TargetRole = optionPtr(profile.TargetRole)
	u.ExperienceLevel = optionPtr(profile.ExperienceLevel)
	u.AvatarURL = b.avatarURL(ctx, profile.ID)
	return u, nil
}

func (b *Builder) Stats(ctx context.Context, userID string, rangeValue string) (Stats, error) {
	current, previous, err := b.currentAndPreviousSessions(ctx, userID, rangeValue, nil)
	if err != nil {
		return Stats{}, errx.Wrap(err)
	}

	summary, hasSummary, err := b.progressSummary(ctx, userID)
	if err != nil {
		return Stats{}, errx.Wrap(err)
	}

	currentMetrics := metricsFromSessions(current)
	previousMetrics := metricsFromSessions(previous)
	streak := StreakValue{}
	if hasSummary {
		streak.Value = summary.CurrentStreak
		streak.IsRecord = summary.CurrentStreak > 0 && summary.CurrentStreak >= summary.LongestStreak
	}

	return Stats{
		TotalInterviews: StatValue{
			Value:          currentMetrics.Interviews,
			DeltaPercent:   deltaPercent(currentMetrics.Interviews, previousMetrics.Interviews),
			DeltaDirection: deltaDirection(currentMetrics.Interviews, previousMetrics.Interviews),
		},
		AverageScore: NullableStatValue{
			Value:          currentMetrics.AverageScore,
			DeltaPercent:   deltaPercentPtr(currentMetrics.AverageScore, previousMetrics.AverageScore),
			DeltaDirection: deltaDirectionPtr(currentMetrics.AverageScore, previousMetrics.AverageScore),
		},
		CurrentStreakDays: streak,
		TotalPracticeSeconds: StatValue{
			Value:          int(currentMetrics.PracticeSeconds),
			DeltaPercent:   deltaPercent64(currentMetrics.PracticeSeconds, previousMetrics.PracticeSeconds),
			DeltaDirection: deltaDirection64(currentMetrics.PracticeSeconds, previousMetrics.PracticeSeconds),
		},
	}, nil
}

func (b *Builder) Performance(
	ctx context.Context,
	userID string,
	rangeValue string,
	topicID *string,
) (Performance, error) {
	current, previous, err := b.currentAndPreviousSessions(ctx, userID, rangeValue, topicID)
	if err != nil {
		return Performance{}, errx.Wrap(err)
	}

	currentMetrics := metricsFromSessions(current)
	previousMetrics := metricsFromSessions(previous)
	return Performance{
		Range: rangeValue,
		Summary: PerformanceSummary{
			AverageScore:        currentMetrics.AverageScore,
			ScoreDeltaPercent:   deltaPercentPtr(currentMetrics.AverageScore, previousMetrics.AverageScore),
			InterviewsCompleted: currentMetrics.Interviews,
			PracticeSeconds:     currentMetrics.PracticeSeconds,
		},
		Points: performancePoints(current),
	}, nil
}

func (b *Builder) TopicOption(topicID *string) *Option {
	if topicID == nil {
		return nil
	}
	return &Option{ID: *topicID, Name: displayName(*topicID)}
}

func (b *Builder) Topics(ctx context.Context, userID string, _ string) (Topics, error) {
	stats, err := b.domainContainer.TopicStatRepo().List(ctx, progress.TopicStatFilter{
		UserID: &userID,
	})
	if err != nil {
		return Topics{}, errx.Wrap(err)
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].AverageScore == stats[j].AverageScore {
			return stats[i].Attempts > stats[j].Attempts
		}
		return stats[i].AverageScore < stats[j].AverageScore
	})

	items := make([]TopicPerformance, 0, len(stats))
	for i := range stats {
		score := roundedPtr(stats[i].AverageScore)
		correctness := ratioPtr(stats[i].AverageScore)
		avgTime := averageTime(stats[i].TotalTimeSpentSeconds, stats[i].Attempts)
		items = append(items, TopicPerformance{
			ID:                 stats[i].TopicKey,
			Name:               displayName(stats[i].TopicKey),
			Score:              score,
			QuestionsAnswered:  stats[i].Attempts,
			CorrectnessRate:    correctness,
			AverageTimeSeconds: avgTime,
			Trend:              "flat",
			Level:              level(score),
		})
	}

	return Topics{
		Items:  items,
		Weak:   weakTopics(items),
		Strong: strongTopics(items),
	}, nil
}

func (b *Builder) RecentActivity(
	ctx context.Context,
	userID string,
	limit int,
	cursor *string,
) (RecentActivity, error) {
	resp, err := b.portalContainer.Interview().ListDashboardSessions(ctx, &interview.ListDashboardSessionsRequest{
		UserID: userID,
		Statuses: []string{
			"in_progress",
			"completed",
			"abandoned",
			"scoring",
		},
		Limit:  limit,
		Cursor: cursor,
	})
	if err != nil {
		return RecentActivity{}, errx.Wrap(err)
	}

	items := make([]ActivityItem, 0, len(resp.Items))
	for i := range resp.Items {
		items = append(items, toActivityItem(resp.Items[i]))
	}

	return RecentActivity{
		Items:      items,
		NextCursor: resp.NextCursor,
	}, nil
}

func (b *Builder) Recommendations(ctx context.Context, userID string) (Recommendations, error) {
	profile, hasProfile, err := b.profile(ctx, userID)
	if err != nil {
		return Recommendations{}, errx.Wrap(err)
	}

	topics, err := b.Topics(ctx, userID, RangeAll)
	if err != nil {
		return Recommendations{}, errx.Wrap(err)
	}

	recommended := recommendedTopics(topics.Items)
	nextTopics := make([]Option, 0, len(recommended))
	for i := range recommended {
		nextTopics = append(nextTopics, Option{ID: recommended[i].ID, Name: recommended[i].Name})
	}

	var targetRole *Option
	var experienceLevel *Option
	if hasProfile {
		targetRole = optionPtr(profile.TargetRole)
		experienceLevel = optionPtr(profile.ExperienceLevel)
		if len(nextTopics) == 0 {
			preferred, prefErr := b.portalContainer.Candidate().ListTopicPreferencesByProfileID(ctx, profile.ID)
			if prefErr != nil {
				return Recommendations{}, errx.Wrap(prefErr)
			}
			for i := range preferred {
				nextTopics = append(
					nextTopics,
					Option{ID: preferred[i].TopicKey, Name: displayName(preferred[i].TopicKey)},
				)
			}
		}
	}

	difficulty := "medium"
	if len(recommended) == 0 {
		difficulty = "mixed"
	}

	return Recommendations{
		RecommendedTopics: recommended,
		NextInterview: NextInterview{
			TargetRole:               targetRole,
			ExperienceLevel:          experienceLevel,
			Topics:                   nextTopics,
			Difficulty:               difficulty,
			QuestionCount:            5,
			EstimatedDurationSeconds: 3600,
		},
	}, nil
}

func (b *Builder) Overview(ctx context.Context, userCtx *auth.UserContext, rangeValue string) (Overview, error) {
	user, err := b.User(ctx, userCtx)
	if err != nil {
		return Overview{}, errx.Wrap(err)
	}
	stats, err := b.Stats(ctx, userCtx.UserID, rangeValue)
	if err != nil {
		return Overview{}, errx.Wrap(err)
	}
	performance, err := b.Performance(ctx, userCtx.UserID, rangeValue, nil)
	if err != nil {
		return Overview{}, errx.Wrap(err)
	}
	topics, err := b.Topics(ctx, userCtx.UserID, rangeValue)
	if err != nil {
		return Overview{}, errx.Wrap(err)
	}
	recent, err := b.RecentActivity(ctx, userCtx.UserID, 10, nil)
	if err != nil {
		return Overview{}, errx.Wrap(err)
	}
	recommendations, err := b.Recommendations(ctx, userCtx.UserID)
	if err != nil {
		return Overview{}, errx.Wrap(err)
	}

	return Overview{
		User:            user,
		Stats:           stats,
		Performance:     performance,
		Topics:          topics,
		RecentActivity:  recent,
		Recommendations: recommendations,
	}, nil
}

func (b *Builder) profile(ctx context.Context, userID string) (*candidate.Profile, bool, error) {
	profile, err := b.portalContainer.Candidate().GetProfileByUserID(ctx, userID)
	if errx.IsCodeIn(err, candidate.CodeProfileNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, errx.Wrap(err)
	}
	return profile, true, nil
}

func (b *Builder) progressSummary(ctx context.Context, userID string) (*progress.Summary, bool, error) {
	summary, err := b.domainContainer.ProgressSummaryRepo().Get(ctx, progress.SummaryFilter{UserID: &userID})
	if errx.IsCodeIn(err, progress.CodeProgressSummaryNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, errx.Wrap(err)
	}
	return summary, true, nil
}

func (b *Builder) avatarURL(ctx context.Context, profileID int64) *string {
	assocType := "avatar"
	files, err := b.portalContainer.Filevault().ListByEntity(ctx, &filevault.ListByEntityRequest{
		EntityType: candidate.EntityTypeProfile,
		EntityID:   profileID,
		AssocType:  &assocType,
	})
	if err != nil || len(files) == 0 {
		return nil
	}
	url := fmt.Sprintf("/api/v1/files/download?id=%s", files[0].ID)
	return &url
}

func (b *Builder) currentAndPreviousSessions(
	ctx context.Context,
	userID string,
	rangeValue string,
	topicID *string,
) ([]interview.DashboardSession, []interview.DashboardSession, error) {
	from, to, prevFrom, prevTo := timeRange(rangeValue)
	current, err := b.completedSessions(ctx, userID, from, to, topicID)
	if err != nil {
		return nil, nil, errx.Wrap(err)
	}
	previous, err := b.completedSessions(ctx, userID, prevFrom, prevTo, topicID)
	if err != nil {
		return nil, nil, errx.Wrap(err)
	}
	return current, previous, nil
}

func (b *Builder) completedSessions(
	ctx context.Context,
	userID string,
	from *time.Time,
	to *time.Time,
	topicID *string,
) ([]interview.DashboardSession, error) {
	resp, err := b.portalContainer.Interview().ListDashboardSessions(ctx, &interview.ListDashboardSessionsRequest{
		UserID:        userID,
		Statuses:      []string{"completed"},
		StartedAtFrom: from,
		StartedAtTo:   to,
		TopicID:       topicID,
		Limit:         1000,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	return resp.Items, nil
}

type metrics struct {
	Interviews      int
	AverageScore    *int
	PracticeSeconds int64
}

func metricsFromSessions(sessions []interview.DashboardSession) metrics {
	var totalScore float64
	scored := 0
	out := metrics{Interviews: len(sessions)}
	for i := range sessions {
		out.PracticeSeconds += sessions[i].DurationSeconds
		if sessions[i].Score != nil {
			totalScore += *sessions[i].Score
			scored++
		}
	}
	if scored > 0 {
		avg := int(math.Round(totalScore / float64(scored)))
		out.AverageScore = &avg
	}
	return out
}

func timeRange(rangeValue string) (*time.Time, *time.Time, *time.Time, *time.Time) {
	now := time.Now().UTC()
	var duration time.Duration
	switch rangeValue {
	case Range30D:
		duration = 30 * 24 * time.Hour
	case Range90D:
		duration = 90 * 24 * time.Hour
	case RangeAll:
		return nil, &now, nil, nil
	default:
		duration = 7 * 24 * time.Hour
	}
	from := now.Add(-duration)
	prevFrom := from.Add(-duration)
	prevTo := from
	return &from, &now, &prevFrom, &prevTo
}

func performancePoints(sessions []interview.DashboardSession) []PerformancePoint {
	type bucket struct {
		date     time.Time
		count    int
		scoreSum float64
		scored   int
		seconds  int64
	}

	buckets := map[string]*bucket{}
	for i := range sessions {
		day := sessions[i].StartedAt.UTC().Truncate(24 * time.Hour)
		key := day.Format("2006-01-02")
		b, ok := buckets[key]
		if !ok {
			b = &bucket{date: day}
			buckets[key] = b
		}
		b.count++
		b.seconds += sessions[i].DurationSeconds
		if sessions[i].Score != nil {
			b.scoreSum += *sessions[i].Score
			b.scored++
		}
	}

	points := make([]PerformancePoint, 0, len(buckets))
	for key, b := range buckets {
		var avg *int
		if b.scored > 0 {
			v := int(math.Round(b.scoreSum / float64(b.scored)))
			avg = &v
		}
		points = append(points, PerformancePoint{
			Date:                key,
			Label:               b.date.Format("Jan 2"),
			AverageScore:        avg,
			InterviewsCompleted: b.count,
			PracticeSeconds:     b.seconds,
		})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Date < points[j].Date
	})
	return points
}

func toActivityItem(s interview.DashboardSession) ActivityItem {
	var score *int
	if s.Score != nil {
		v := int(math.Round(*s.Score))
		score = &v
	}
	var completedAt *string
	if s.CompletedAt != nil {
		v := s.CompletedAt.UTC().Format(time.RFC3339)
		completedAt = &v
	}
	topics := make([]Option, 0, len(s.Topics))
	for i := range s.Topics {
		topics = append(topics, Option{ID: s.Topics[i].ID, Name: s.Topics[i].Name})
	}
	return ActivityItem{
		SessionID:       s.ID,
		Title:           s.Title,
		Status:          s.Status,
		Score:           score,
		StartedAt:       s.StartedAt.UTC().Format(time.RFC3339),
		CompletedAt:     completedAt,
		DurationSeconds: s.DurationSeconds,
		QuestionCount:   s.QuestionCount,
		AnsweredCount:   s.AnsweredCount,
		Topics:          topics,
	}
}

func weakTopics(items []TopicPerformance) []HighlightTopic {
	out := []HighlightTopic{}
	for i := range items {
		if items[i].Score == nil || *items[i].Score >= 70 {
			continue
		}
		action := fmt.Sprintf("Practice more %s questions.", items[i].Name)
		out = append(out, HighlightTopic{
			ID:                items[i].ID,
			Name:              items[i].Name,
			Score:             items[i].Score,
			QuestionsAnswered: items[i].QuestionsAnswered,
			Reason:            "Lowest average score in the selected range.",
			RecommendedAction: &action,
		})
		if len(out) == 3 {
			return out
		}
	}
	return out
}

func strongTopics(items []TopicPerformance) []HighlightTopic {
	copied := append([]TopicPerformance(nil), items...)
	sort.Slice(copied, func(i, j int) bool {
		if copied[i].Score == nil {
			return false
		}
		if copied[j].Score == nil {
			return true
		}
		return *copied[i].Score > *copied[j].Score
	})

	out := []HighlightTopic{}
	for i := range copied {
		if copied[i].Score == nil || *copied[i].Score < 80 {
			continue
		}
		out = append(out, HighlightTopic{
			ID:                copied[i].ID,
			Name:              copied[i].Name,
			Score:             copied[i].Score,
			QuestionsAnswered: copied[i].QuestionsAnswered,
			Reason:            "Consistently high scores across recent sessions.",
		})
		if len(out) == 3 {
			return out
		}
	}
	return out
}

func recommendedTopics(items []TopicPerformance) []RecommendedTopic {
	weak := weakTopics(items)
	out := make([]RecommendedTopic, 0, len(weak))
	for i := range weak {
		priority := "medium"
		if weak[i].Score != nil && *weak[i].Score < 60 {
			priority = "high"
		}
		out = append(out, RecommendedTopic{
			ID:       weak[i].ID,
			Name:     weak[i].Name,
			Priority: priority,
			Reason:   fmt.Sprintf("Your %s score is below your overall average.", weak[i].Name),
		})
	}
	return out
}

func deltaPercent(current int, previous int) int {
	if previous == 0 {
		if current > 0 {
			return 100
		}
		return 0
	}
	return int(math.Round((float64(current-previous) / float64(previous)) * 100))
}

func deltaPercent64(current int64, previous int64) int {
	return deltaPercent(int(current), int(previous))
}

func deltaPercentPtr(current *int, previous *int) int {
	if current == nil || previous == nil {
		return 0
	}
	return deltaPercent(*current, *previous)
}

func deltaDirection(current int, previous int) string {
	switch {
	case previous == 0 && current > 0:
		return "new"
	case current > previous:
		return "up"
	case current < previous:
		return "down"
	default:
		return "flat"
	}
}

func deltaDirection64(current int64, previous int64) string {
	return deltaDirection(int(current), int(previous))
}

func deltaDirectionPtr(current *int, previous *int) string {
	if current == nil || previous == nil {
		if current != nil {
			return "new"
		}
		return "flat"
	}
	return deltaDirection(*current, *previous)
}

func displayName(id string) string {
	switch id {
	case "python":
		return "Python Backend"
	case "golang":
		return "Golang Backend"
	case "javascript":
		return "JavaScript"
	case "algorithms":
		return "Algorithms"
	case "system_design", "System Design":
		return "System Design"
	case "database_design", "Database Design":
		return "Database Design"
	case "security":
		return "Security"
	case "api_design":
		return "API Design"
	case "junior":
		return "Junior"
	case "mid":
		return "Mid-Level"
	case "senior":
		return "Senior"
	default:
		return id
	}
}

func optionPtr(id *string) *Option {
	if id == nil {
		return nil
	}
	return &Option{ID: *id, Name: displayName(*id)}
}

func roundedPtr(v float64) *int {
	out := int(math.Round(v))
	return &out
}

func ratioPtr(v float64) *float64 {
	out := math.Round((v/100)*100) / 100
	return &out
}

func averageTime(totalSeconds int64, attempts int) *int64 {
	if attempts == 0 {
		return nil
	}
	out := totalSeconds / int64(attempts)
	return &out
}

func level(score *int) string {
	if score == nil {
		return "stable"
	}
	switch {
	case *score >= 80:
		return "strong"
	case *score >= 70:
		return "stable"
	default:
		return "needs_practice"
	}
}
