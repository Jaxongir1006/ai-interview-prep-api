package listerrors

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/domain/alerterror"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/pagination"
	"github.com/rise-and-shine/pkg/sorter"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	pagination.Request

	Code        *string `query:"code"`
	Service     *string `query:"service"`
	Operation   *string `query:"operation"`
	Alerted     *bool   `query:"alerted"`
	CreatedFrom *string `query:"created_from"`
	CreatedTo   *string `query:"created_to"`
	Search      string  `query:"search"`
	Sort        string  `query:"sort"`
}

type UseCase = ucdef.UserAction[*Request, *pagination.Response[alerterror.Error]]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "list-errors" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*pagination.Response[alerterror.Error], error) {
	// Normalize pagination params
	in.Normalize()

	// Parse sort options
	sortOpts := sorter.MakeFromStr(in.Sort, "created_at", "code", "service", "operation")

	// Build filter
	filter := alerterror.Filter{
		Code:      in.Code,
		Service:   in.Service,
		Operation: in.Operation,
		Alerted:   in.Alerted,
		Search:    in.Search,
		Limit:     ptrInt(in.Limit()),
		Offset:    ptrInt(in.Offset()),
		SortOpts:  sortOpts,
	}

	if in.CreatedFrom != nil {
		from, err := parseTime(*in.CreatedFrom)
		if err != nil {
			return nil, errx.New(
				"invalid 'created_from' time format, expected RFC3339",
				errx.WithType(errx.T_Validation),
				errx.WithCode("INVALID_TIME_FORMAT"),
			)
		}
		filter.CreatedFrom = &from
	}
	if in.CreatedTo != nil {
		to, err := parseTime(*in.CreatedTo)
		if err != nil {
			return nil, errx.New(
				"invalid 'created_to' time format, expected RFC3339",
				errx.WithType(errx.T_Validation),
				errx.WithCode("INVALID_TIME_FORMAT"),
			)
		}
		filter.CreatedTo = &to
	}

	// List errors with filters, search, sorting, and pagination
	items, count, err := uc.domainContainer.AlertErrorRepo().ListWithCount(ctx, filter)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return paginated errors
	resp := pagination.NewResponse(items, count, in.Request)
	return &resp, nil
}

func ptrInt(v int) *int {
	return &v
}

func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
