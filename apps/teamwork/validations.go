package teamwork

import (
	"context"

	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
)

type TeamworkServiceValidations interface {
	ListOnduty(context.Context, *teamworkv1.ListOndutyRequest) error
	ListTaskByProject(context.Context, *teamworkv1.ListTaskByProjectRequest) error
	ListTaskByUser(context.Context, *teamworkv1.ListTaskByUserRequest) error
}

type teamworkServiceValidations struct {
}

func CreateTeamworkServiceValidations() TeamworkServiceValidations {
	return new(teamworkServiceValidations)
}

func (w *teamworkServiceValidations) ListOnduty(context.Context, *teamworkv1.ListOndutyRequest) error {
	return nil
}

func (w *teamworkServiceValidations) ListTaskByProject(context.Context, *teamworkv1.ListTaskByProjectRequest) error {
	return nil
}

func (w *teamworkServiceValidations) ListTaskByUser(context.Context, *teamworkv1.ListTaskByUserRequest) error {
	return nil
}
