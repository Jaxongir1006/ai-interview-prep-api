package postgres

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/catalog"

	"github.com/code19m/errx"
	"github.com/uptrace/bun"
)

type catalogRepo struct {
	idb bun.IDB
}

func NewCatalogRepo(idb bun.IDB) catalog.Repo {
	return &catalogRepo{idb: idb}
}

func (r *catalogRepo) ListActiveTargetRoles(ctx context.Context) ([]catalog.TargetRole, error) {
	var items []catalog.TargetRole
	err := r.idb.NewSelect().
		Model(&items).
		ModelTableExpr(schemaName+".interview_target_roles AS itr").
		Where("is_active = TRUE").
		Order("display_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	return items, nil
}

func (r *catalogRepo) ListActiveExperienceLevels(ctx context.Context) ([]catalog.ExperienceLevel, error) {
	var items []catalog.ExperienceLevel
	err := r.idb.NewSelect().
		Model(&items).
		ModelTableExpr(schemaName+".interview_experience_levels AS iel").
		Where("is_active = TRUE").
		Order("display_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	return items, nil
}

func (r *catalogRepo) ListActiveTopics(ctx context.Context) ([]catalog.Topic, error) {
	var items []catalog.Topic
	err := r.idb.NewSelect().
		Model(&items).
		ModelTableExpr(schemaName+".interview_topics AS it").
		Where("is_active = TRUE").
		Order("display_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	return items, nil
}

func (r *catalogRepo) ListActiveTargetRoleTopics(ctx context.Context) ([]catalog.TargetRoleTopic, error) {
	var items []catalog.TargetRoleTopic
	err := r.idb.NewSelect().
		Model(&items).
		ModelTableExpr(schemaName+".interview_target_role_topics AS itrtt").
		Where("is_active = TRUE").
		Order("display_order ASC", "id ASC").
		Scan(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	return items, nil
}
