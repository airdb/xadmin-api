package bchm

import (
	"context"
	"fmt"

	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	"github.com/airdb/xadmin-api/pkg/authkit"
)

type BchmServiceValidations interface {
	ListLosts(ctx context.Context, request *bchmv1.ListLostsRequest) error
	GetLost(ctx context.Context, request *bchmv1.GetLostRequest) error
	CreateLost(ctx context.Context, request *bchmv1.CreateLostRequest) error
	UpdateLost(ctx context.Context, request *bchmv1.UpdateLostRequest) error
	DeleteLost(ctx context.Context, request *bchmv1.DeleteLostRequest) error
}

type bchmServiceValidations struct {
}

func CreateBchmServiceValidations() BchmServiceValidations {
	return new(bchmServiceValidations)
}

func (w *bchmServiceValidations) ListLosts(ctx context.Context, request *bchmv1.ListLostsRequest) error {
	if request.GetPageSize() == 0 {
		request.PageSize = 10
	}
	if request.GetPageOffset() < 0 {
		request.PageOffset = 0
	}

	return nil
}

func (w *bchmServiceValidations) GetLost(ctx context.Context, request *bchmv1.GetLostRequest) error {
	return nil
}

func (w *bchmServiceValidations) CreateLost(ctx context.Context, request *bchmv1.CreateLostRequest) error {
	user := authkit.FromContextUser(ctx)
	if user == nil || !user.IsAdmin {
		return fmt.Errorf("您无权限执行该操作")
	}
	if len(request.GetName()) == 0 {
		return fmt.Errorf("请输入 姓名")
	}

	if !request.GetBirthedAt().IsValid() {
		return fmt.Errorf("请输入 出生日期")
	}

	if !request.GetMissedAt().IsValid() {
		return fmt.Errorf("请输入 失踪时间")
	}

	if len(request.GetMissedAddr()) == 0 {
		return fmt.Errorf("请输入 失踪地点")
	}

	if len(request.GetMissedHeight()) == 0 {
		return fmt.Errorf("请输入 失踪时身高")
	}

	if len(request.GetCharacter()) == 0 {
		return fmt.Errorf("请输入 特征")
	}

	if len(request.GetDetails()) == 0 {
		return fmt.Errorf("请输入 失踪详情")
	}

	if len(request.GetCategory()) == 0 {
		return fmt.Errorf("请输入 寻亲类型")
	}

	if len(request.GetDataFrom()) == 0 {
		return fmt.Errorf("请输入 信息来源")
	}

	if len(request.GetFollower()) == 0 {
		return fmt.Errorf("请输入 跟进志愿者")
	}

	return nil
}

func (w *bchmServiceValidations) UpdateLost(ctx context.Context, request *bchmv1.UpdateLostRequest) error {
	user := authkit.FromContextUser(ctx)
	if user == nil || !user.IsAdmin {
		return fmt.Errorf("您无权限执行该操作")
	}
	return nil
}

func (w *bchmServiceValidations) DeleteLost(ctx context.Context, request *bchmv1.DeleteLostRequest) error {
	user := authkit.FromContextUser(ctx)
	if user == nil || !user.IsAdmin {
		return fmt.Errorf("您无权限执行该操作")
	}
	return nil
}
