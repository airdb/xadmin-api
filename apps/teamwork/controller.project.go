package teamwork

import (
	"context"
	"errors"
	"fmt"

	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"github.com/airdb/xadmin-api/pkg/idkit"
	"github.com/airdb/xadmin-api/pkg/querykit"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/samber/lo"
)

func (c *teamworkController) ListProjects(ctx context.Context, request *teamworkv1.ListProjectsRequest) (*teamworkv1.ListProjectsResponse, error) {
	c.log.Debug(ctx, "list projects accepted")

	total, filtered, err := c.deps.ProjectRepo.Count(ctx, request)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "list projects count error")
		return nil, errors.New("list projects count error")
	}
	if total == 0 {
		return nil, errors.New("projects is empty")
	}

	items, err := c.deps.ProjectRepo.List(ctx, request)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "list projects error")
		return nil, errors.New("list projects error")
	}

	return &teamworkv1.ListProjectsResponse{
		TotalSize:    total,
		FilteredSize: filtered,
		Projects: func() []*teamworkv1.Project {
			res := make([]*teamworkv1.Project, len(items))
			for i := 0; i < len(items); i++ {
				res[i] = c.conver.FromModelProjectToProtoProject(items[i])
			}
			return res
		}(),
	}, nil
}

func (c *teamworkController) GetProject(ctx context.Context, request *teamworkv1.GetProjectRequest) (*teamworkv1.GetProjectResponse, error) {
	c.log.Debug(ctx, "get project accepted")

	item, err := c.deps.ProjectRepo.Get(ctx, idkit.MustFromString(request.GetId()))
	if err != nil {
		c.log.WithError(err).Debug(ctx, "get project error")
		return nil, errors.New("name not exist")
	}

	return &teamworkv1.GetProjectResponse{
		Project: c.conver.FromModelProjectToProtoProject(item),
	}, err
}

func (c *teamworkController) CreateProject(ctx context.Context, request *teamworkv1.CreateProjectRequest) (*teamworkv1.CreateProjectResponse, error) {
	c.log.Debug(ctx, "create project accepted")

	item := c.conver.FromProtoCreateProjectToModelProject(request)
	err := c.deps.ProjectRepo.Create(ctx, item)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "create project item failed")
		return nil, errors.New("create project item failed")
	}

	return &teamworkv1.CreateProjectResponse{
		Project: c.conver.FromModelProjectToProtoProject(item),
	}, err
}

func (c *teamworkController) UpdateProject(ctx context.Context, request *teamworkv1.UpdateProjectRequest) (*teamworkv1.UpdateProjectResponse, error) {
	c.log.Debug(ctx, "update project accepted")
	data := c.conver.FromProtoProjectToModelProject(request.GetProject())

	fm := querykit.NewField(request.GetUpdateMask(), request.GetProject()).WithAction("update")

	err := c.deps.ProjectRepo.Update(ctx, data.Id, data, fm)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "update project item failed")
		return nil, errors.New("update project item failed")
	}

	item, err := c.deps.ProjectRepo.Get(ctx, data.Id)
	if err != nil || item == nil {
		c.log.WithError(err).Debug(ctx, "update project item not exist")
		return nil, errors.New("update project item not exist")
	}

	return &teamworkv1.UpdateProjectResponse{
		Project: c.conver.FromModelProjectToProtoProject(item),
	}, err
}

func (c *teamworkController) DeleteProject(ctx context.Context, request *teamworkv1.DeleteProjectRequest) (*empty.Empty, error) {
	c.log.Debug(ctx, "delete project accepted")

	err := c.deps.ProjectRepo.Delete(ctx, idkit.MustFromString(request.GetId()))
	if err != nil {
		c.log.WithError(err).Debug(ctx, "delete project item failed")
		return nil, errors.New("delete project item failed")
	}

	return &empty.Empty{}, nil
}

func (c *teamworkController) AssignProjectIssues(ctx context.Context, request *teamworkv1.AssignProjectIssuesRequest) (*empty.Empty, error) {
	c.log.Debug(ctx, "assign project issues accepted")

	project, err := c.deps.ProjectRepo.Get(ctx, idkit.MustFromString(request.GetId()))
	if err != nil || project == nil {
		c.log.WithError(err).Debug(ctx, "get project error")
		return nil, errors.New("project not exist")
	}

	issues, err := c.deps.IssueRepo.FindByIds(ctx, request.GetIssueIds())
	if err != nil {
		c.log.WithError(err).Debug(ctx, "assign project issues failed")
		return nil, errors.New("assign project issues failed")
	}

	if len(issues) != len(request.GetIssueIds()) {
		ids := make([]string, len(issues))
		for _, issue := range issues {
			if !issue.ProjectId.IsNil() {
				continue
			}
			ids = append(ids, issue.Id.String())
		}
		l, _ := lo.Difference(request.GetIssueIds(), ids)
		return nil, fmt.Errorf("issues %s not exist or already assigned", l)
	}

	ids := make([]idkit.Id, len(issues))
	for i := 0; i < len(issues); i++ {
		ids[i] = issues[i].Id
	}
	err = c.deps.IssueRepo.AssignToProject(ctx, ids, project.Id)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
