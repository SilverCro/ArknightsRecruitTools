package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "kratosTest/api/myKratos/v1"
)

// ArkOperatorName  is a ... model.
type ArkOperatorName struct {
	Name string
}

type ArkOperatorInfo struct {
	Name       string
	TarList    []string
	Profession string
	Rarity     int32
}

func (aoi *ArkOperatorInfo) ToArkOperatorInfo() *v1.ArkOperatorInfo {
	return &v1.ArkOperatorInfo{
		Name:       aoi.Name,
		TarList:    aoi.TarList,
		Profession: aoi.Profession,
		Rarity:     aoi.Rarity,
	}
}

type ArkRecruitTags struct {
	Tags []string
}

type ArkRecruitRecommendInfo struct {
	RecommendTags          []string
	RecommendOperatorInfos []*ArkOperatorInfo
}

func (arri *ArkRecruitRecommendInfo) ToRecruitRecommend() *v1.RecruitRecommend {
	recommendOperatorInfos := make([]*v1.ArkOperatorInfo, len(arri.RecommendOperatorInfos))
	for index, recommendOperatorInfo := range arri.RecommendOperatorInfos {
		recommendOperatorInfos[index] = recommendOperatorInfo.ToArkOperatorInfo()
	}
	return &v1.RecruitRecommend{
		RecommendTags:          arri.RecommendTags,
		RecommendOperatorInfos: recommendOperatorInfos,
	}
}

// MyKratosRepo  .
type MyKratosRepo interface {
	GetAllArkOperatorName(context.Context) ([]*ArkOperatorName, error)
	GetArkOperatorInfo(context.Context, *ArkOperatorName) (*ArkOperatorInfo, error)
	GetRecruitTags(context.Context, string) (*ArkRecruitTags, error)
	GetAllRecruitOperatorName(context.Context) ([]*ArkOperatorName, map[string]bool, map[string]bool, error)
	// Save(context.Context, *MyKratos) (*MyKratos, error)
	// Update(context.Context, *MyKratos) (*MyKratos, error)
	// FindByID(context.Context, int64) (*MyKratos, error)
	// ListByHello(context.Context, string) ([]*MyKratos, error)
	// ListAll(context.Context) ([]*MyKratos, error)
}

// MyKratosUsecase .
type MyKratosUsecase struct {
	repo MyKratosRepo
	log  *log.Helper
}

// NewMyKratosUsecase .
func NewMyKratosUsecase(repo MyKratosRepo, logger log.Logger) *MyKratosUsecase {
	return &MyKratosUsecase{repo: repo, log: log.NewHelper(logger)}
}

// GetArkOperatorInfo creates an ArkOperatorInfo, and returns an ArkOperatorInfo.
func (uc *MyKratosUsecase) GetArkOperatorInfo(ctx context.Context, aon *ArkOperatorName) (*ArkOperatorInfo, error) {
	aoi, err := uc.repo.GetArkOperatorInfo(ctx, aon)
	uc.log.WithContext(ctx).Infof("GetArkOperatorInfo: %v", aon)
	return aoi, err
}

// GetRecruitRecommendInfo creates an ArkRecruitRecommendInfo, and returns an ArkRecruitRecommendInfo.
func (uc *MyKratosUsecase) GetRecruitRecommendInfo(ctx context.Context, art *ArkRecruitTags) (string, []*ArkRecruitRecommendInfo, error) {
	uc.log.WithContext(ctx).Infof("ArkRecruitTags: %v", art.Tags)
	stateStr := ""
	arkRecruitRecommendInfo := make([]*ArkRecruitRecommendInfo, 0)
	if len(art.Tags) == 0 {
		stateStr = "fail 请输入足够的tag"
		return stateStr, arkRecruitRecommendInfo, nil
	}
	if len(art.Tags) > 5 {
		stateStr = "fail 输入的tag数量超过5个"
		return stateStr, arkRecruitRecommendInfo, nil
	}
	tagTypes := []string{"trait", "profession", "position", "rarity"}
	_, allTagMap, err := getRecruitTagMap(ctx, tagTypes, uc)
	if err != nil {
		return "fail 获取RecruitTags失败", nil, err
	}
	checkRepetitionMap := make(map[string]bool)
	for _, tag := range art.Tags {
		if !allTagMap[tag] {
			stateStr = "fail 输入无效的tag"
			return stateStr, arkRecruitRecommendInfo, nil
		}
		if checkRepetitionMap[tag] {
			stateStr = "fail 输入重复的tag"
			return stateStr, arkRecruitRecommendInfo, nil
		}
		checkRepetitionMap[tag] = true
	}
	recruitOpNames, _, _, err := uc.repo.GetAllRecruitOperatorName(ctx)
	if err != nil {
		return "fail 获取RecruitOperatorNames失败", nil, err
	}
	tagUsed := &ArkRecruitTags{}
	tagUsed.Tags = make([]string, 0)
	err = tagCombination(ctx, uc, 0, 0, art, recruitOpNames, tagUsed, &arkRecruitRecommendInfo)
	if err != nil {
		return "fail 获取推荐tag失败!", nil, err
	}
	if len(arkRecruitRecommendInfo) != 0 {
		return "success 获取推荐tag成功", arkRecruitRecommendInfo, nil
	}
	return "fail 没有推荐的tag选择", nil, nil
}

func tagCombination(ctx context.Context, uc *MyKratosUsecase, tagCount int, tagNumberNow int, art *ArkRecruitTags, recruitOpNames []*ArkOperatorName, tagUsed *ArkRecruitTags, arkRecruitRecommendInfo *[]*ArkRecruitRecommendInfo) error {
	if tagCount == 3 || tagNumberNow == len(art.Tags) {
		if tagCount != 0 {
			eligibleOpList := make([]*ArkOperatorInfo, 0)
			for _, name := range recruitOpNames {
				flag, info, err := CheckOperatorTags(ctx, uc, &ArkOperatorName{Name: name.Name}, tagUsed)
				if err != nil {
					return err
				}
				seniorExperiencedFlag := false
				for _, val := range tagUsed.Tags {
					if val == "高级资深干员" {
						seniorExperiencedFlag = true
					}
				}
				if flag && (info.Rarity != 6 || seniorExperiencedFlag) {
					eligibleOpList = append(eligibleOpList, info)
				}
				uc.log.WithContext(ctx).Infof("checkName:%s" + info.Name)
				uc.log.WithContext(ctx).Infof("checkResult:%t", flag)
			}
			flag := checkRarityOnly(eligibleOpList)
			// uc.log.WithContext(ctx).Infof("len(eligibleOpList):", len(eligibleOpList))
			// uc.log.WithContext(ctx).Infof("checkRarityOnlyResult:", flag)
			if flag {
				recommendTags := make([]string, len(tagUsed.Tags))
				copy(recommendTags, tagUsed.Tags)
				*arkRecruitRecommendInfo = append(*arkRecruitRecommendInfo, &ArkRecruitRecommendInfo{
					RecommendTags:          recommendTags,
					RecommendOperatorInfos: eligibleOpList,
				})
			}
			uc.log.WithContext(ctx).Infof("len(arkRecruitRecommendInfo):%d", len(*arkRecruitRecommendInfo))
		}
		return nil
	}

	err := tagCombination(ctx, uc, tagCount, tagNumberNow+1, art, recruitOpNames, tagUsed, arkRecruitRecommendInfo)
	if err != nil {
		return err
	}
	tagUsed.Tags = append(tagUsed.Tags, art.Tags[tagNumberNow])
	err = tagCombination(ctx, uc, tagCount+1, tagNumberNow+1, art, recruitOpNames, tagUsed, arkRecruitRecommendInfo)
	if err != nil {
		return err
	}
	tagUsed.Tags = tagUsed.Tags[:len(tagUsed.Tags)-1]
	return nil
}

func getRecruitTagMap(ctx context.Context, tagTypes []string, uc *MyKratosUsecase) (map[string]map[string]bool, map[string]bool, error) {
	recruitTagMap := make(map[string]map[string]bool)
	allTagMap := make(map[string]bool)
	for _, tagType := range tagTypes {
		recruitTags, err := uc.repo.GetRecruitTags(ctx, tagType)
		if err != nil {
			return nil, nil, err
		}
		m := make(map[string]bool)
		for _, tag := range recruitTags.Tags {
			m[tag] = true
			allTagMap[tag] = true
		}
		recruitTagMap[tagType] = m
	}
	return recruitTagMap, allTagMap, nil
}

func CheckOperatorTags(ctx context.Context, uc *MyKratosUsecase, aon *ArkOperatorName, arts *ArkRecruitTags) (bool, *ArkOperatorInfo, error) {
	arkOperatorInfo, err := uc.repo.GetArkOperatorInfo(ctx, aon)

	if err != nil {
		uc.log.WithContext(ctx).Infof("aonName:%s", aon.Name)
		uc.log.WithContext(ctx).Infof("err:%v", err)
		return false, nil, err
	}
	tagCount := 0
	for _, tag := range arts.Tags {
		flag := 0
		for _, t := range arkOperatorInfo.TarList {
			if t == tag {
				flag = 1
				break
			}
		}
		if flag == 1 {
			tagCount++
		}
	}
	return tagCount == len(arts.Tags), arkOperatorInfo, nil
}

func checkRarityOnly(eligibleOpList []*ArkOperatorInfo) bool {
	if len(eligibleOpList) == 0 {
		return false
	}
	oneStarCount := 0
	aboveFourStarCount := 0
	for _, info := range eligibleOpList {
		if info.Rarity == 1 {
			oneStarCount++
		}
		if info.Rarity >= 4 {
			aboveFourStarCount++
		}
	}
	return oneStarCount+aboveFourStarCount == len(eligibleOpList)
}
