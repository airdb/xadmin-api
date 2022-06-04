package services

import (
	"context"

	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (w *teamworkImpl) CreateProject(ctx context.Context, request *teamworkv1.CreateProjectRequest) (*teamworkv1.CreateProjectResponse, error) {
	return w.deps.Controller.CreateProject(ctx, request)
}

func (w *teamworkImpl) GetProject(ctx context.Context, request *teamworkv1.GetProjectRequest) (*teamworkv1.GetProjectResponse, error) {
	return w.deps.Controller.GetProject(ctx, request)
}

func (w *teamworkImpl) ListProjects(ctx context.Context, request *teamworkv1.ListProjectsRequest) (*teamworkv1.ListProjectsResponse, error) {
	return w.deps.Controller.ListProjects(ctx, request)
}

func (w *teamworkImpl) UpdateProject(ctx context.Context, request *teamworkv1.UpdateProjectRequest) (*teamworkv1.UpdateProjectResponse, error) {
	return w.deps.Controller.UpdateProject(ctx, request)
}

func (w *teamworkImpl) DeleteProject(ctx context.Context, request *teamworkv1.DeleteProjectRequest) (*emptypb.Empty, error) {
	return w.deps.Controller.DeleteProject(ctx, request)
}

func (w *teamworkImpl) AssignProjectIssues(ctx context.Context, request *teamworkv1.AssignProjectIssuesRequest) (*emptypb.Empty, error) {
	return w.deps.Controller.AssignProjectIssues(ctx, request)
}
