package catalog

import "context"

type Repo interface {
	ListActiveTargetRoles(ctx context.Context) ([]TargetRole, error)
	ListActiveExperienceLevels(ctx context.Context) ([]ExperienceLevel, error)
	ListActiveTopics(ctx context.Context) ([]Topic, error)
	ListActiveTargetRoleTopics(ctx context.Context) ([]TargetRoleTopic, error)
}
