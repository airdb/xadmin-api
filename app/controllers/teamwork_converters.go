package controllers

import (
	"github.com/airdb/xadmin-api/app/data"
	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"github.com/airdb/xadmin-api/pkg/datatypes"
	"github.com/airdb/xadmin-api/pkg/idkit"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type teamworkConvert struct{}

func newTeamworkConvert() *teamworkConvert {
	return &teamworkConvert{}
}

// FromProtoIssueToModelIssue converts proto model to our data Entity
func (c teamworkConvert) FromProtoIssueToModelIssue(request *teamworkv1.Issue) *data.IssueEntity {
	if request == nil {
		return nil
	}
	return &data.IssueEntity{
		PrimaryKey: datatypes.PrimaryKey{
			Id: idkit.MustFromString(request.GetId()),
		},
		Title:   request.GetTitle(),
		Content: request.GetContent(),
	}
}

// FromProtoIssueToModelIssue converts proto model to our data Entity
func (c teamworkConvert) FromProtoCreateIssueToModelIssue(request *teamworkv1.CreateIssueRequest) *data.IssueEntity {
	if request == nil {
		return nil
	}
	return &data.IssueEntity{
		Title:   request.Issue.GetTitle(),
		Content: request.Issue.GetContent(),
	}
}

// FromModelIssueToProtoIssue converts our data Entity to proto model
func (c teamworkConvert) FromModelIssueToProtoIssue(in *data.IssueEntity) *teamworkv1.Issue {
	if in == nil {
		return nil
	}

	return &teamworkv1.Issue{
		Id: in.Id.String(),
		CreatedAt: func() *timestamppb.Timestamp {
			t, err := ptypes.TimestampProto(in.CreatedAt)
			if err != nil {
				return nil
			}
			return t
		}(),
		CreatedBy: in.CreatedBy,
		Title:     in.Title,
		Content:   in.Content,
	}
}
