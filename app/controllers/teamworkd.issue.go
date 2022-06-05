package controllers

import (
	"context"
	"errors"

	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"github.com/airdb/xadmin-api/pkg/idkit"
	"github.com/airdb/xadmin-api/pkg/querykit"
	"github.com/golang/protobuf/ptypes/empty"
)

func (c *teamworkController) ListIssues(ctx context.Context, request *teamworkv1.ListIssuesRequest) (*teamworkv1.ListIssuesResponse, error) {
	c.log.Debug(ctx, "list issues accepted")

	total, filtered, err := c.deps.IssueRepo.Count(ctx, request)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "list issues count error")
		return nil, errors.New("list issues count error")
	}
	if total == 0 {
		return nil, errors.New("issues is empty")
	}

	items, err := c.deps.IssueRepo.List(ctx, request)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "list issues error")
		return nil, errors.New("list issues error")
	}

	return &teamworkv1.ListIssuesResponse{
		TotalSize:    total,
		FilteredSize: filtered,
		Issues: func() []*teamworkv1.Issue {
			res := make([]*teamworkv1.Issue, len(items))
			for i := 0; i < len(items); i++ {
				res[i] = c.conver.FromModelIssueToProtoIssue(items[i])
			}
			return res
		}(),
	}, nil
}

func (c *teamworkController) GetIssue(ctx context.Context, request *teamworkv1.GetIssueRequest) (*teamworkv1.GetIssueResponse, error) {
	c.log.Debug(ctx, "get issue accepted")

	item, err := c.deps.IssueRepo.Get(ctx, idkit.MustFromString(request.GetId()))
	if err != nil {
		c.log.WithError(err).Debug(ctx, "get issue error")
		return nil, errors.New("name not exist")
	}

	return &teamworkv1.GetIssueResponse{
		Issue: c.conver.FromModelIssueToProtoIssue(item),
	}, err
}

func (c *teamworkController) CreateIssue(ctx context.Context, request *teamworkv1.CreateIssueRequest) (*teamworkv1.CreateIssueResponse, error) {
	c.log.Debug(ctx, "create issue accepted")

	item := c.conver.FromProtoCreateIssueToModelIssue(request)
	err := c.deps.IssueRepo.Create(ctx, item)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "create issue item failed")
		return nil, errors.New("create issue item failed")
	}

	return &teamworkv1.CreateIssueResponse{
		Issue: c.conver.FromModelIssueToProtoIssue(item),
	}, err
}

func (c *teamworkController) UpdateIssue(ctx context.Context, request *teamworkv1.UpdateIssueRequest) (*teamworkv1.UpdateIssueResponse, error) {
	c.log.Debug(ctx, "update issue accepted")
	data := c.conver.FromProtoIssueToModelIssue(request.GetIssue())

	fm := querykit.NewField(request.GetUpdateMask(), request.GetIssue()).WithAction("update")

	err := c.deps.IssueRepo.Update(ctx, data.Id, data, fm)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "update issue item failed")
		return nil, errors.New("update issue item failed")
	}

	item, err := c.deps.IssueRepo.Get(ctx, data.Id)
	if err != nil || item == nil {
		c.log.WithError(err).Debug(ctx, "update issue item not exist")
		return nil, errors.New("update issue item not exist")
	}

	return &teamworkv1.UpdateIssueResponse{
		Issue: c.conver.FromModelIssueToProtoIssue(item),
	}, err
}

func (c *teamworkController) DeleteIssue(ctx context.Context, request *teamworkv1.DeleteIssueRequest) (*empty.Empty, error) {
	c.log.Debug(ctx, "delete issue accepted")

	err := c.deps.IssueRepo.Delete(ctx, idkit.MustFromString(request.GetId()))
	if err != nil {
		c.log.WithError(err).Debug(ctx, "delete issue item failed")
		return nil, errors.New("delete issue item failed")
	}

	return &empty.Empty{}, nil
}
