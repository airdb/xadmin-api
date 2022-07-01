package bchm

import (
	"github.com/airdb/xadmin-api/apps/data"
	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Convert struct{}

func newConvert() *Convert {
	return &Convert{}
}

// FromProtoLostToModelLost converts proto model to our data Entity
func (c Convert) FromProtoLostToModelLost(request *bchmv1.Lost) *data.LostEntity {
	if request == nil {
		return nil
	}
	return &data.LostEntity{
		ID: uint(request.GetId()),

		Nickname:  request.GetName(),
		Gender:    uint(request.Gender),
		BirthedAt: request.GetBirthedAt().AsTime(),

		MissedAt:       request.GetMissedAt().AsTime(),
		MissedCountry:  request.GetMissedCountry(),
		MissedProvince: request.GetMissedProvince(),
		MissedCity:     request.GetMissedCity(),
		MissedAddress:  request.GetMissedAddr(),
		Height:         request.GetMissedHeight(),

		Characters: request.GetCharacter(),
		Details:    request.GetDetails(),

		Category: request.GetCategory(),
		DataFrom: request.GetDataFrom(),
		Follower: request.GetFollower(),
	}
}

// FromProtoLostToModelLost converts proto model to our data Entity
func (c Convert) FromProtoCreateLostToModelLost(request *bchmv1.CreateLostRequest) *data.LostEntity {
	if request == nil {
		return nil
	}
	return &data.LostEntity{
		Nickname:  request.GetName(),
		Gender:    uint(request.Gender),
		BirthedAt: request.GetBirthedAt().AsTime(),

		MissedAt:       request.GetMissedAt().AsTime(),
		MissedCountry:  request.GetMissedCountry(),
		MissedProvince: request.GetMissedProvince(),
		MissedCity:     request.GetMissedCity(),
		MissedAddress:  request.GetMissedAddr(),
		Height:         request.GetMissedHeight(),

		Characters: request.GetCharacter(),
		Details:    request.GetDetails(),

		Category: request.GetCategory(),
		DataFrom: request.GetDataFrom(),
		Follower: request.GetFollower(),
	}
}

// FromModelLostToProtoLost converts our data Entity to proto model
func (c Convert) FromModelLostToProtoLost(in *data.LostEntity, files []*data.FileEntity) *bchmv1.Lost {
	if in == nil {
		return nil
	}

	return &bchmv1.Lost{
		Id: int32(in.ID),

		Name:   in.Nickname,
		Gender: uint32(in.Gender),
		BirthedAt: func() *timestamppb.Timestamp {
			t, err := ptypes.TimestampProto(in.BirthedAt)
			if err != nil {
				return nil
			}
			return t
		}(),
		Carousel: func() []*bchmv1.Lost_Carousel {
			if files == nil {
				return nil
			}
			items := make([]*bchmv1.Lost_Carousel, len(files))
			for k, v := range files {
				items[k] = &bchmv1.Lost_Carousel{
					Id:    int32(v.ID),
					Title: "",
					Url:   v.URL,
				}
			}
			return items
		}(),
		MissedAt: func() *timestamppb.Timestamp {
			t, err := ptypes.TimestampProto(in.BirthedAt)
			if err != nil {
				return nil
			}
			return t
		}(),
		MissedCountry:  in.MissedCountry,
		MissedProvince: in.MissedProvince,
		MissedCity:     in.MissedCity,
		MissedAddr:     in.MissedAddress,
		MissedHeight:   in.Height,

		Character: in.Characters,
		Details:   in.Details,

		Category: in.Category,
		DataFrom: in.DataFrom,
		Follower: in.Follower,
	}
}
