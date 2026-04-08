package uow

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/profile"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/topicpreference"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase"
)

type Factory = uowbase.Factory[UnitOfWork]

type UnitOfWork interface {
	uowbase.UnitOfWork

	Profile() profile.Repo
	TopicPreference() topicpreference.Repo
}
