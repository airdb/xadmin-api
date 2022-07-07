package bchm

import (
	"context"

	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/genproto/googleapis/api/httpbody"
)

func (w *serviceImpl) ListLosts(ctx context.Context, request *bchmv1.ListLostsRequest) (*bchmv1.ListLostsResponse, error) {
	w.log.WithField("request", request).Debug(ctx, "list lost request")
	if err := w.deps.Validations.ListLosts(ctx, request); err != nil {
		return nil, err
	}

	return w.deps.Controller.ListLosts(ctx, request)
}

func (w *serviceImpl) GetLost(ctx context.Context, request *bchmv1.GetLostRequest) (*bchmv1.GetLostResponse, error) {
	w.log.WithField("request", request).Debug(ctx, "get lost request")
	if err := w.deps.Validations.GetLost(ctx, request); err != nil {
		return nil, err
	}
	return w.deps.Controller.GetLost(ctx, request)
}

func (w *serviceImpl) ShareLostCallback(ctx context.Context, request *bchmv1.ShareLostCallbackRequest) (*bchmv1.ShareLostCallbackResponse, error) {
	return w.deps.Controller.ShareLostCallback(ctx, request)
}

func (w *serviceImpl) GetLostMpCode(ctx context.Context, request *bchmv1.GetLostMpCodeRequest) (*httpbody.HttpBody, error) {
	response, err := w.deps.Controller.GetLostMpCode(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (w *serviceImpl) CreateLost(ctx context.Context, request *bchmv1.CreateLostRequest) (*bchmv1.CreateLostResponse, error) {
	w.log.WithField("request", request).Debug(ctx, "create lost request")
	if err := w.deps.Validations.CreateLost(ctx, request); err != nil {
		return nil, err
	}
	return w.deps.Controller.CreateLost(ctx, request)
}

func (w *serviceImpl) UpdateLost(ctx context.Context, request *bchmv1.UpdateLostRequest) (result *bchmv1.UpdateLostResponse, err error) {
	w.log.WithField("request", request).Debug(ctx, "update lost request")
	err = w.deps.Validations.UpdateLost(ctx, request)
	if err == nil {
		result, err = w.deps.Controller.UpdateLost(ctx, request)
	}
	w.log.WithError(err).Debug(ctx, "update lost done")
	return
}

func (w *serviceImpl) UpdateLostAudited(ctx context.Context, request *bchmv1.UpdateLostAuditedRequest) (result *empty.Empty, err error) {
	w.log.WithField("request", request).Debug(ctx, "update lost audited request")
	result, err = w.deps.Controller.UpdateLostAudited(ctx, request)
	w.log.WithError(err).Debug(ctx, "update lost audited done")
	return
}

func (w *serviceImpl) UpdateLostDone(ctx context.Context, request *bchmv1.UpdateLostDoneRequest) (result *empty.Empty, err error) {
	w.log.WithField("request", request).Debug(ctx, "update lost done request")
	result, err = w.deps.Controller.UpdateLostDone(ctx, request)
	w.log.WithError(err).Debug(ctx, "update lost done done")
	return
}

func (w *serviceImpl) DeleteLost(ctx context.Context, request *bchmv1.DeleteLostRequest) (result *empty.Empty, err error) {
	w.log.WithField("request", request).Debug(ctx, "delete lost request")
	err = w.deps.Validations.DeleteLost(ctx, request)
	if err == nil {
		result, err = w.deps.Controller.DeleteLost(ctx, request)
	}
	w.log.WithError(err).Debug(ctx, "delete lost done")
	return
}
