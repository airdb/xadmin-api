package teamwork

import (
	"github.com/airdb/xadmin-api/apps/data"
	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"github.com/airdb/xadmin-api/pkg/datatypes"
	"github.com/airdb/xadmin-api/pkg/idkit"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type teamworkConvert struct{}

func newTeamworkConvert() *teamworkConvert {
	return &teamworkConvert{}
}

// Project Convert Start

// FromProtoProjectToModelProject converts proto model to our data Entity
func (c teamworkConvert) FromProtoProjectToModelProject(request *teamworkv1.Project) *data.ProjectEntity {
	if request == nil {
		return nil
	}
	return &data.ProjectEntity{
		PrimaryKey: datatypes.PrimaryKey{
			Id: idkit.MustFromString(request.GetId()),
		},
		Title:     request.GetTitle(),
		Milestone: request.GetMilestone(),
		Status:    request.GetStatus(),
	}
}

// FromProtoProjectToModelProject converts proto model to our data Entity
func (c teamworkConvert) FromProtoCreateProjectToModelProject(request *teamworkv1.CreateProjectRequest) *data.ProjectEntity {
	if request == nil {
		return nil
	}
	return &data.ProjectEntity{
		Title:     request.Project.GetTitle(),
		Milestone: request.Project.GetMilestone(),
		Status:    request.Project.GetStatus(),
	}
}

// FromModelProjectToProtoProject converts our data Entity to proto model
func (c teamworkConvert) FromModelProjectToProtoProject(in *data.ProjectEntity) *teamworkv1.Project {
	if in == nil {
		return nil
	}

	return &teamworkv1.Project{
		Id:        in.Id.String(),
		CreatedAt: timestamppb.New(in.CreatedAt),
		CreatedBy: in.CreatedBy,
		Title:     in.Title,
		Milestone: in.Milestone,
		Status:    in.Status,
	}
}

// Project Convert End

// Issue Convert Start

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
		Id:        in.Id.String(),
		CreatedAt: timestamppb.New(in.CreatedAt),
		CreatedBy: in.CreatedBy,
		Title:     in.Title,
		Content:   in.Content,
		ProjectId: func() string {
			if in.ProjectId.IsNil() {
				return ""
			}
			return in.ProjectId.String()
		}(),
	}
}

// Issue Convert End
