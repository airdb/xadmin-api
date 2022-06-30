package bchm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/airdb/xadmin-api/apps/data"
	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	"github.com/airdb/xadmin-api/pkg/cachekit"
	"github.com/airdb/xadmin-api/pkg/ipkit"
	"github.com/airdb/xadmin-api/pkg/querykit"
	"github.com/airdb/xadmin-api/pkg/wechatkit"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Controller responsible for the business logic of our BchmService
type Controller interface {
	bchmv1.ServiceServer
}

type controllerDeps struct {
	fx.In

	Config       cfg.Config
	Logger       log.Logger
	Cache        *cachekit.Cache
	LostRepo     data.LostRepo
	LostStatRepo data.LostStatRepo
}

type controller struct {
	bchmv1.UnimplementedServiceServer

	deps   controllerDeps
	log    log.Fields
	cache  *cachekit.Redis
	conver *Convert
}

// CreateController is a constructor for Fx
func CreateController(deps controllerDeps) Controller {
	return &controller{
		deps: deps,
		log:  deps.Logger.WithField("controller", "bchm"),
		cache: deps.Cache.Redis(
			deps.Config.Get(fmt.Sprintf("service.%s.redis.db")).Int(),
		),
		conver: newConvert(),
	}
}

func (c *controller) BuildFilterScope(q *bchmv1.ListLostsRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(q.GetKeyword()) > 0 {
			queryWord := "%" + q.GetKeyword() + "%"
			db = db.Where("(nickname like ?) OR (missed_address like ?)",
				queryWord,
				queryWord,
			)
		}

		if len(q.GetCategory()) > 0 {
			db = db.Where("category = ?", q.GetCategory())
		}

		return db
	}
}

func (c *controller) ListLosts(ctx context.Context, request *bchmv1.ListLostsRequest) (*bchmv1.ListLostsResponse, error) {
	c.log.Debug(ctx, "list losts accepted")

	total, filtered, err := c.deps.LostRepo.Count(ctx, request)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "list losts count error")
		return nil, errors.New("list losts count error")
	}
	if total == 0 {
		return nil, errors.New("losts is empty")
	}

	items, err := c.deps.LostRepo.List(ctx, request)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "list losts error")
		return nil, errors.New("list losts error")
	}

	return &bchmv1.ListLostsResponse{
		TotalSize:    total,
		FilteredSize: filtered,
		Losts: func() []*bchmv1.Lost {
			res := make([]*bchmv1.Lost, len(items))
			for i := 0; i < len(items); i++ {
				res[i] = c.conver.FromModelLostToProtoLost(items[i])
			}
			return res
		}(),
	}, nil
}

func (c *controller) GetLost(ctx context.Context, request *bchmv1.GetLostRequest) (*bchmv1.GetLostResponse, error) {
	c.log.Debug(ctx, "get lost accepted")

	item, err := c.deps.LostRepo.Get(ctx, uint(request.GetId()))
	if err != nil {
		c.log.WithError(err).Debug(ctx, "get lost error")
		return nil, errors.New("name not exist")
	}

	return &bchmv1.GetLostResponse{
		Lost: c.conver.FromModelLostToProtoLost(item),
	}, err
}

func (c *controller) ShareLostCallback(ctx context.Context, request *bchmv1.ShareLostCallbackRequest) (*bchmv1.ShareLostCallbackResponse, error) {
	shareKey := strings.Join([]string{
		request.GetShareKey(), ipkit.RemoteIp(ctx)}, ":")

	item, err := c.deps.LostRepo.Get(ctx, uint(request.GetId()))
	if err != nil {
		return nil, errors.New("the lost is not exist")
	}

	shareKeyRedisValue, err := c.cache.Get(shareKey)
	if err != nil {
		return nil, errors.New("get share key info failed")
	}

	var shareCount int
	if shareKeyRedisValue != "" {
		shareCount, err = strconv.Atoi(shareKeyRedisValue)
		if err != nil {
			return nil, errors.New("get share key value error")
		}
	}

	if shareCount >= 3 {
		return nil, errors.New("reached share limit")
	}

	shareCount++

	err = c.cache.Set(shareKey, strconv.Itoa(shareCount), time.Second*86400)
	if err != nil {
		return nil, errors.New("set share count failed")
	}

	if err = c.deps.LostStatRepo.IncreaseShare(ctx, uint(request.GetId())); err != nil {
		return nil, errors.New("increase share failed")
	}

	return &bchmv1.ShareLostCallbackResponse{
		Lost: c.conver.FromModelLostToProtoLost(item),
	}, nil
}

func (c *controller) GetLostMpCode(ctx context.Context, request *bchmv1.GetLostMpCodeRequest) (*bchmv1.GetLostMpCodeResponse, error) {
	wx := wechatkit.NewWechatMiniProgram(wechatkit.NewWechat())
	code, err := wx.CodeUnlimit(
		`pages/redirect/wxmpcode`,
		fmt.Sprintf("id=%d&s=bbhj.lost", request.GetId()),
	)
	c.log.WithField("code", code).Info(ctx, "new mp code generated")
	if err != nil {
		c.log.WithError(err).Error(ctx, "can not generate wechat mini programe code")
		return nil, errors.New("can not generate wechat mini programe code")
	}

	return &bchmv1.GetLostMpCodeResponse{
		Code: code,
	}, nil
}

func (c *controller) CreateLost(ctx context.Context, request *bchmv1.CreateLostRequest) (*bchmv1.CreateLostResponse, error) {
	c.log.Debug(ctx, "create lost accepted")

	item := c.conver.FromProtoCreateLostToModelLost(request)
	err := c.deps.LostRepo.Create(ctx, item)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "create lost item failed")
		return nil, errors.New("create lost item failed")
	}

	return &bchmv1.CreateLostResponse{
		Lost: c.conver.FromModelLostToProtoLost(item),
	}, err
}

func (c *controller) UpdateLost(ctx context.Context, request *bchmv1.UpdateLostRequest) (*bchmv1.UpdateLostResponse, error) {
	c.log.Debug(ctx, "update lost accepted")
	data := c.conver.FromProtoLostToModelLost(request.GetLost())

	fm := querykit.NewField(request.GetUpdateMask(), request.GetLost()).WithAction("update")

	err := c.deps.LostRepo.Update(ctx, data.ID, data, fm)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "update lost item failed")
		return nil, errors.New("update lost item failed")
	}

	item, err := c.deps.LostRepo.Get(ctx, uint(data.ID))
	if err != nil || item == nil {
		c.log.WithError(err).Debug(ctx, "update lost item not exist")
		return nil, errors.New("update lost item not exist")
	}

	return &bchmv1.UpdateLostResponse{
		Lost: c.conver.FromModelLostToProtoLost(item),
	}, err
}

func (c *controller) DeleteLost(ctx context.Context, request *bchmv1.DeleteLostRequest) (*empty.Empty, error) {
	c.log.Debug(ctx, "delete lost accepted")

	err := c.deps.LostRepo.Delete(ctx, uint(request.GetId()))
	if err != nil {
		c.log.WithError(err).Debug(ctx, "delete lost item failed")
		return nil, errors.New("delete lost item failed")
	}

	return &empty.Empty{}, nil
}
