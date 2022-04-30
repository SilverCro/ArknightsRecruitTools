package service

import (
	"context"
	v1 "kratosTest/api/myKratos/v1"

	"kratosTest/internal/biz"
)

// GreeterService is a greeter service.
type MyKratosService struct {
	v1.UnimplementedMyKratosServer
	gu *biz.GreeterUsecase
	mu *biz.MyKratosUsecase
}

// NewMyKratosService new a greeter service.
func NewMyKratosService(mu *biz.MyKratosUsecase, gu *biz.GreeterUsecase) *MyKratosService {
	return &MyKratosService{mu: mu, gu: gu}
}

// SayHello implements helloworld.GreeterServer.
func (s *MyKratosService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.gu.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

func (s *MyKratosService) GetArkOperatorInfo(ctx context.Context, in *v1.GetArkOperatorInfoRequest) (*v1.GetArkOperatorInfoReply, error) {
	aoi, err := s.mu.GetArkOperatorInfo(ctx, &biz.ArkOperatorName{Name: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.GetArkOperatorInfoReply{
		ArkOperatorInfo: &v1.GetArkOperatorInfoReply_ArkOperatorInfo{
			Name:       aoi.Name,
			TarList:    aoi.TarList,
			Profession: aoi.Profession,
			Rarity:     aoi.Rarity,
		},
	}, nil

}

func (s *MyKratosService) ArkRecruitRecommendTool(ctx context.Context, in *v1.ArkRecruitRecommendRequest) (*v1.ArkRecruitRecommendReply, error) {
	//g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	//if err != nil {
	//	return nil, err
	//}
	rs, arris, err := s.mu.GetRecruitRecommendInfo(ctx, &biz.ArkRecruitTags{Tags: in.Tags})
	if err != nil {
		return nil, err
	}
	recruitRecommends := make([]*v1.RecruitRecommend, len(arris))
	for index, val := range arris {
		recruitRecommends[index] = val.ToRecruitRecommend()
	}
	return &v1.ArkRecruitRecommendReply{
		Status:            rs,
		RecruitRecommends: recruitRecommends,
	}, nil
}
