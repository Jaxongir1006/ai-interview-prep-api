package embassy

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/actionlog"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/statuschangelog"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/mask"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/samber/lo"
)

type embassy struct {
	domainContainer *domain.Container
}

func New(domainContainer *domain.Container) audit.Portal {
	return &embassy{
		domainContainer: domainContainer,
	}
}

func (e *embassy) Log(ctx context.Context, action audit.Action, opts ...audit.LogOption) error {
	cfg := audit.BuildLogConfig(opts)
	userID := lo.ToPtr(meta.Find(ctx, meta.ActorID))
	traceID := meta.Find(ctx, meta.TraceID)
	now := time.Now()

	var owned bool
	uow, err := e.domainContainer.UOWFactory().NewBorrowed(ctx)
	if err != nil {
		uow, err = e.domainContainer.UOWFactory().NewUOW(ctx)
		if err != nil {
			return errx.Wrap(err)
		}
		owned = true
		defer uow.DiscardUnapplied()
	}

	var groupKey *string
	if cfg.GroupKey != "" {
		groupKey = &cfg.GroupKey
	}

	al, err := uow.ActionLog().Create(ctx, &actionlog.ActionLog{
		UserID:         userID,
		Module:         action.Module,
		OperationID:    action.OperationID,
		RequestPayload: mask.StructToOrdMap(action.Payload),
		Tags:           cfg.Tags,
		GroupKey:       groupKey,
		IPAddress:      meta.Find(ctx, meta.IPAddress),
		UserAgent:      meta.Find(ctx, meta.UserAgent),
		TraceID:        traceID,
		CreatedAt:      now,
	})
	if err != nil {
		return errx.Wrap(err)
	}

	if len(cfg.StatusChanges) > 0 {
		scLogs := make([]statuschangelog.StatusChangeLog, len(cfg.StatusChanges))
		for i, sc := range cfg.StatusChanges {
			scLogs[i] = statuschangelog.StatusChangeLog{
				ActionLogID: al.ID,
				EntityType:  sc.EntityType,
				EntityID:    sc.EntityID,
				Status:      sc.Status,
				TraceID:     traceID,
				CreatedAt:   now,
			}
		}

		err = uow.StatusChangeLog().BulkCreate(ctx, scLogs)
		if err != nil {
			return errx.Wrap(err)
		}
	}

	if owned {
		err = uow.ApplyChanges()
		if err != nil {
			return errx.Wrap(err)
		}
	}

	return nil
}
