package controllers

import (
	"context"

	"github.com/airdb/xadmin-api/app/data"
	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TeamworkServiceController responsible for the business logic of our TeamworkService
type TeamworkServiceController interface {
	teamworkv1.TeamworkServiceServer
}

type teamworkInfoControllerDeps struct {
	fx.In

	DB       data.LostRepo
	LostRepo data.LostRepo
	Logger   log.Logger
}

type teamworkInfoController struct {
	teamworkv1.UnimplementedTeamworkServiceServer

	deps teamworkInfoControllerDeps
	log  log.Fields
}

// CreateTeamworkServiceController is a constructor for Fx
func CreateTeamworkServiceController(deps teamworkInfoControllerDeps) TeamworkServiceController {
	return &teamworkInfoController{
		deps: deps,
		log:  deps.Logger.WithField("controller", "teamwork"),
	}
}

func (c *teamworkInfoController) ListOnduty(ctx context.Context, request *teamworkv1.ListOndutyRequest) (*teamworkv1.ListOndutyResponse, error) {
	c.log.Debug(ctx, "list onduty accepted")
	return &teamworkv1.ListOndutyResponse{
		// Return recent 3 weeks.
		Schedule: []*teamworkv1.Onduty{
			{
				Id:          1,
				Year:        2022,
				Week:        19,
				TeamName:    "airdb",
				OndutyEmail: "dean@airdb.net",
				CreatedAt:   "2020-01-01",
				CreatedBy:   "dean",
			},
			{
				Id:          2,
				Year:        2022,
				Week:        20,
				TeamName:    "airdb",
				OndutyEmail: "dean@airdb.net",
				CreatedAt:   "2020-01-01",
				CreatedBy:   "dean",
			},
			{
				Id:          3,
				Year:        2022,
				Week:        21,
				TeamName:    "airdb",
				OndutyEmail: "dean@airdb.net",
				CreatedAt:   "2020-01-01",
				CreatedBy:   "dean",
			},
		},
	}, nil
}

func (c *teamworkInfoController) ListTaskByProject(ctx context.Context, request *teamworkv1.ListTaskByProjectRequest) (*teamworkv1.ListTaskByProjectResponse, error) {
	c.log.Debug(ctx, "list task by project accepted")

	return &teamworkv1.ListTaskByProjectResponse{
		Project: []*teamworkv1.Project{
			{
				Id:               1,
				ProjectName:      "项目申报",
				ProjectMilestone: "phase 1: 完成ppt演示",
				ProjectStatus:    "进行中",
				TaskProcess: []*teamworkv1.TaskProcess{
					{
						Email:    "dean@airdb.net",
						ThisWeek: "完成ppp demo",
						NextWeek: "完成ppt演示",
					},
					{
						Email:    "lucy@airdb.net",
						ThisWeek: "完成 part1",
						NextWeek: "完成 part2",
					},
					{
						Email:    "lily@airdb.net",
						ThisWeek: "完成 part3",
						NextWeek: "完成 part4",
					},
				},
			},
			{
				Id:               2,
				ProjectName:      "sso + xadmin ",
				ProjectMilestone: "phase 1: sso 成功部署 \n phase 2: xadmin 开发框架 \n Phase 3: 前后端集成",
				ProjectStatus:    "进行中",
				TaskProcess: []*teamworkv1.TaskProcess{
					{
						Email:    "bumu@airdb.net",
						ThisWeek: "框架基本完成",
						NextWeek: "代码合并",
					},
					{
						Email:    "xiaoqiang",
						ThisWeek: "修复部分bug",
						NextWeek: "优化",
					},
					{
						Email:    "zhiye",
						ThisWeek: "开发环境搭建",
						NextWeek: "功能集成",
					},
				},
			},
			{
				Id:               3,
				ProjectName:      "论坛迁移",
				ProjectMilestone: "phase 1: CDN 迁移 \n phase 2: 用户迁移 \n Phase 3: 论坛前端开发",
				ProjectStatus:    "进行中",
				TaskProcess: []*teamworkv1.TaskProcess{
					{
						Email:    "bumu@airdb.net",
						ThisWeek: "sso 迁移代码优化",
						NextWeek: "sso 迁移代码优化",
					},
					{
						Email:    "serven",
						ThisWeek: "未开始",
						NextWeek: "cdn部署",
					},
					{
						Email:    "lingling",
						ThisWeek: "未开始",
						NextWeek: "论坛前端开发",
					},
				},
			},
		},
	}, nil
}

func (c *teamworkInfoController) ListTaskByUser(ctx context.Context, request *teamworkv1.ListTaskByUserRequest) (*teamworkv1.ListTaskByUserResponse, error) {
	c.log.Debug(ctx, "list task by user accepted")
	return nil, status.Errorf(codes.Unimplemented, "method ListTaskByUser not implemented")
}
