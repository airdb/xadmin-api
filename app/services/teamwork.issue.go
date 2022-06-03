package services

import (
	"context"

	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (w *teamworkImpl) CreateIssue(ctx context.Context, request *teamworkv1.CreateIssueRequest) (*teamworkv1.CreateIssueResponse, error) {
	return w.deps.Controller.CreateIssue(ctx, request)
}

func (w *teamworkImpl) GetIssue(ctx context.Context, request *teamworkv1.GetIssueRequest) (*teamworkv1.GetIssueResponse, error) {
	return w.deps.Controller.GetIssue(ctx, request)
}

func (w *teamworkImpl) ListIssues(ctx context.Context, request *teamworkv1.ListIssuesRequest) (*teamworkv1.ListIssuesResponse, error) {
	return w.deps.Controller.ListIssues(ctx, request)
}

func (w *teamworkImpl) DeleteIssue(ctx context.Context, request *teamworkv1.DeleteIssueRequest) (*emptypb.Empty, error) {
	return w.deps.Controller.DeleteIssue(ctx, request)
}

func (w *teamworkImpl) UpdateIssue(ctx context.Context, request *teamworkv1.UpdateIssueRequest) (*teamworkv1.UpdateIssueResponse, error) {
	return w.deps.Controller.UpdateIssue(ctx, request)
}
