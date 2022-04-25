package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Greeter is a Greeter model.
type ArkOperatorName struct {
	Name string
}

type ArkOperatorInfo struct {
	Name       string
	TarList    []string
	Profession string
	Rarity     int32
}

type ArkRecruitTags struct {
	Tags []string
}

type ArkRecruitRecommendInfo struct {
	State                 string
	RecommendTags         []string
	RecommendOperatorInfo []*ArkOperatorInfo
}

// GreeterRepo is a Greater repo.
type MyKratosRepo interface {
	GetAllArkOperatorName(context.Context) ([]*ArkOperatorName, error)
	GetArkOperatorInfo(context.Context, *ArkOperatorName) (*ArkOperatorInfo, error)
	CheckOperatorTags(context.Context, *ArkOperatorName, *ArkRecruitTags) (bool, error)
	//Save(context.Context, *MyKratos) (*MyKratos, error)
	//Update(context.Context, *MyKratos) (*MyKratos, error)
	//FindByID(context.Context, int64) (*MyKratos, error)
	//ListByHello(context.Context, string) ([]*MyKratos, error)
	//ListAll(context.Context) ([]*MyKratos, error)
}

// GreeterUsecase is a Greeter usecase.
type MyKratosUsecase struct {
	repo MyKratosRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewMyKratosUsecase(repo MyKratosRepo, logger log.Logger) *MyKratosUsecase {
	return &MyKratosUsecase{repo: repo, log: log.NewHelper(logger)}
}

// GetArkOperatorInfo creates an ArkOperatorInfo, and returns an ArkOperatorInfo.
func (uc *MyKratosUsecase) GetArkOperatorInfo(ctx context.Context, aon *ArkOperatorName) (*ArkOperatorInfo, error) {

	aoi, err := uc.repo.GetArkOperatorInfo(ctx, aon)
	uc.log.WithContext(ctx).Infof("GetArkOperatorInfo: %v", err)
	return aoi, err
}

// CreateRecruitRecommendInfo creates an ArkRecruitRecommendInfo, and returns an ArkRecruitRecommendInfo.
func (uc *MyKratosUsecase) GetRecruitRecommendInfo(ctx context.Context, art *ArkRecruitTags) (*ArkRecruitRecommendInfo, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", art.Tags)
	arkRecruitRecommendInfo := &ArkRecruitRecommendInfo{}
	if len(art.Tags) <= 0 {
		arkRecruitRecommendInfo.State = "fail:请输入足够的tag"
		return arkRecruitRecommendInfo, nil
	}

	//aons, err := uc.repo.GetAllArkOperatorName(ctx)

	return nil, nil
}
