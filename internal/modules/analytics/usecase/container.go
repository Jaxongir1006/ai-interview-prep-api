package usecase

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/getoverview"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/getperformancetrend"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/getrecentactivity"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/getrecommendations"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/getstats"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/gettopics"
)

type Container struct {
	getDashboardOverview        getoverview.UseCase
	getDashboardStats           getstats.UseCase
	getPerformanceTrend         getperformancetrend.UseCase
	getDashboardTopics          gettopics.UseCase
	getRecentActivity           getrecentactivity.UseCase
	getDashboardRecommendations getrecommendations.UseCase
}

func NewContainer(
	getDashboardOverview getoverview.UseCase,
	getDashboardStats getstats.UseCase,
	getPerformanceTrend getperformancetrend.UseCase,
	getDashboardTopics gettopics.UseCase,
	getRecentActivity getrecentactivity.UseCase,
	getDashboardRecommendations getrecommendations.UseCase,
) *Container {
	return &Container{
		getDashboardOverview:        getDashboardOverview,
		getDashboardStats:           getDashboardStats,
		getPerformanceTrend:         getPerformanceTrend,
		getDashboardTopics:          getDashboardTopics,
		getRecentActivity:           getRecentActivity,
		getDashboardRecommendations: getDashboardRecommendations,
	}
}

func (c *Container) GetDashboardOverview() getoverview.UseCase {
	return c.getDashboardOverview
}

func (c *Container) GetDashboardStats() getstats.UseCase {
	return c.getDashboardStats
}

func (c *Container) GetPerformanceTrend() getperformancetrend.UseCase {
	return c.getPerformanceTrend
}

func (c *Container) GetDashboardTopics() gettopics.UseCase {
	return c.getDashboardTopics
}

func (c *Container) GetRecentActivity() getrecentactivity.UseCase {
	return c.getRecentActivity
}

func (c *Container) GetDashboardRecommendations() getrecommendations.UseCase {
	return c.getDashboardRecommendations
}
