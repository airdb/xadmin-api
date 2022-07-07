package bchm

import (
	"context"
	"errors"
	"fmt"
	"net/url"
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
	"github.com/google/uuid"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"go.uber.org/fx"
	"google.golang.org/genproto/googleapis/api/httpbody"
	fmpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	"gorm.io/gorm"
)

const (
	LOST_WXMP_CODE_FILENAME = `wxmp_code.jpg`
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
	FileRepo     data.FileRepo
	LostRepo     data.LostRepo
	LostStatRepo data.LostStatRepo
}

type controller struct {
	bchmv1.UnimplementedServiceServer

	deps   controllerDeps
	log    log.Fields
	cache  *cachekit.Redis
	wxmp   *wechatkit.WechatMiniProgram
	conver *Convert
}

// CreateController is a constructor for Fx
func CreateController(deps controllerDeps) Controller {
	return &controller{
		deps: deps,
		log:  deps.Logger.WithField("controller", "bchm"),
		cache: deps.Cache.Redis(
			deps.Config.Get(fmt.Sprintf("service.%s.redis.db", "bchm")).Int(),
		),
		wxmp: wechatkit.NewWechatMiniProgram(&config.Config{
			AppID:     deps.Config.Get("xadmin.wechat.app_id").String(),
			AppSecret: deps.Config.Get("xadmin.wechat.app_secret").String(),
		}, cachekit.RedisPool(
			fmt.Sprintf("%s:%d",
				deps.Config.Get("xadmin.cache.redis.host").String(),
				deps.Config.Get("xadmin.cache.redis.port").Int(),
			),
			deps.Config.Get("xadmin.cache.redis.password").String(),
			deps.Config.Get("xadmin.wechat.redis_cache_db").Int(),
		), deps.Logger),
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
				res[i] = c.conver.FromModelLostToProtoLost(items[i], nil)
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
		return nil, errors.New("lost not exist")
	}
	files, err := c.deps.FileRepo.GetLostByID(ctx, item.ID)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "get lost file error")
		return nil, errors.New("lost file not exist")
	}

	return &bchmv1.GetLostResponse{
		Lost: c.conver.FromModelLostToProtoLost(item, files),
		WxMore: &bchmv1.LostWxMore{
			ShareAppMessage: &bchmv1.LostWxMore_ShareAppMessage{
				ShareKey: uuid.New().String(),
				Title:    item.Title,
				ImageUrl: item.AvatarURL,
			},
			ShareTimeline: &bchmv1.LostWxMore_ShareTimeline{
				ShareKey: uuid.New().String(),
				Title:    item.Title,
				Query: func() string {
					query := url.Values{}
					query.Add("lost_id", strconv.Itoa(int(item.ID)))
					return query.Encode()
				}(),
				ImageUrl: item.AvatarURL,
			},
			CodeUnlimit: &bchmv1.LostWxMore_CodeUnlimit{
				Url: fmt.Sprintf(`/v1/bchm/lost/%d/%s`, item.ID, LOST_WXMP_CODE_FILENAME),
			},
		},
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
		Lost: c.conver.FromModelLostToProtoLost(item, nil),
	}, nil
}

func (c *controller) GetLostMpCode(ctx context.Context, request *bchmv1.GetLostMpCodeRequest) (*httpbody.HttpBody, error) {
	code, err := c.wxmp.CodeUnlimit(
		`pages/redirect/wxmpcode`,
		fmt.Sprintf("id=%d&s=bbhj.lost", request.GetId()),
	)
	c.log.WithField("code length", len(code)).Info(ctx, "new mp code generated")
	if err != nil {
		c.log.WithError(err).Error(ctx, "can not generate wechat mini programe code")
		return nil, errors.New("can not generate wechat mini programe code")
	}

	return &httpbody.HttpBody{
		ContentType: "image/jpeg",
		Data:        code,
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

	for k, v := range request.GetImages() {
		err = c.deps.FileRepo.Create(ctx, &data.FileEntity{
			Type:     "lost",
			SortID:   k,
			ParentID: item.ID,
			URL:      v,
		})
		if err != nil {
			c.log.WithError(err).Debug(ctx, "create lost images failed")
			return nil, errors.New("can not create file")
		}
	}

	return &bchmv1.CreateLostResponse{
		Lost: c.conver.FromModelLostToProtoLost(item, nil),
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
		Lost: c.conver.FromModelLostToProtoLost(item, nil),
	}, err
}

func (c *controller) UpdateLostAudited(ctx context.Context, request *bchmv1.UpdateLostAuditedRequest) (*empty.Empty, error) {
	c.log.Debug(ctx, "update lost audited accepted")
	data := &data.LostEntity{
		ID:      uint(request.GetId()),
		Audited: request.GetValue(),
	}

	fm := querykit.NewField(&fmpb.FieldMask{
		Paths: []string{
			"audited",
		},
	}, nil)

	err := c.deps.LostRepo.Update(ctx, data.ID, data, fm)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "update lost audited failed")
		return nil, errors.New("update lost audited failed")
	}

	item, err := c.deps.LostRepo.Get(ctx, uint(data.ID))
	if err != nil || item == nil {
		c.log.WithError(err).Debug(ctx, "update lost item not exist")
		return nil, errors.New("update lost item not exist")
	}

	return &empty.Empty{}, err
}

func (c *controller) UpdateLostDone(ctx context.Context, request *bchmv1.UpdateLostDoneRequest) (*empty.Empty, error) {
	c.log.Debug(ctx, "update lost done accepted")
	data := &data.LostEntity{
		ID:   uint(request.GetId()),
		Done: request.GetValue(),
	}

	fm := querykit.NewField(&fmpb.FieldMask{
		Paths: []string{
			"done",
		},
	}, nil)

	err := c.deps.LostRepo.Update(ctx, data.ID, data, fm)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "update lost done failed")
		return nil, errors.New("update lost done failed")
	}

	item, err := c.deps.LostRepo.Get(ctx, uint(data.ID))
	if err != nil || item == nil {
		c.log.WithError(err).Debug(ctx, "update lost item not exist")
		return nil, errors.New("update lost item not exist")
	}

	return &empty.Empty{}, err
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
