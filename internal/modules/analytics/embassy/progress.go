package embassy

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/progress"
	analyticsportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/analytics"

	"github.com/code19m/errx"
)

func (e *embassy) GetProgressSummaryByUserID(
	ctx context.Context,
	userID string,
) (*analyticsportal.ProgressSummary, error) {
	s, err := e.domainContainer.ProgressSummaryRepo().Get(ctx, progress.SummaryFilter{
		UserID: &userID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &analyticsportal.ProgressSummary{
		ID:                    s.ID,
		UserID:                s.UserID,
		CurrentStreak:         s.CurrentStreak,
		LongestStreak:         s.LongestStreak,
		TotalInterviewsTaken:  s.TotalInterviewsTaken,
		TotalTimeSpentSeconds: s.TotalTimeSpentSeconds,
		AverageScore:          s.AverageScore,
		LastInterviewAt:       s.LastInterviewAt,
		CreatedAt:             s.CreatedAt,
		UpdatedAt:             s.UpdatedAt,
	}, nil
}
