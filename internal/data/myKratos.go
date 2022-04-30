package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"kratosTest/internal/biz"
	"strconv"
)

type myKratosRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewMyKratosRepo(data *Data, logger log.Logger) biz.MyKratosRepo {
	return &myKratosRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func (r *myKratosRepo) GetAllArkOperatorName(ctx context.Context) ([]*biz.ArkOperatorName, error) {
	names, err := r.data.rdb.SMembers(ctx, "name:chara").Result()
	if err != nil {
		return nil, err
	}
	aons := make([]*biz.ArkOperatorName, len(names))
	for index, name := range names {
		aons[index] = &biz.ArkOperatorName{}
		aons[index].Name = name
	}
	return aons, nil
}

func (r *myKratosRepo) GetArkOperatorInfo(ctx context.Context, name *biz.ArkOperatorName) (*biz.ArkOperatorInfo, error) {
	if val, err := r.data.rdb.SIsMember(ctx, "name:chara", name.Name).Result(); err != nil {
		return nil, err
	} else {
		if !val {
			fmt.Println("Cannot find the operator")
			return &biz.ArkOperatorInfo{
				Name: "不存在该干员",
			}, nil
		}
	}
	aoi := &biz.ArkOperatorInfo{
		Name: name.Name,
	}
	if val, err := r.data.rdb.Get(ctx, name.Name+":profession:chara").Result(); err != nil {
		return nil, err
	} else {
		aoi.Profession = val
	}

	if val, err := r.data.rdb.Get(ctx, name.Name+":rarity:chara").Result(); err != nil {
		return nil, err
	} else {
		if v, err1 := strconv.Atoi(val); err1 != nil { //nolint:gosec
			panic(err1)
		} else {
			aoi.Rarity = int32(v)
		}

	}
	val, err := r.data.rdb.SMembers(ctx, name.Name+":tagList:chara").Result()
	if err != nil {
		return nil, err
	}
	aoi.TarList = val
	return aoi, nil
}

func (r *myKratosRepo) CheckOperatorTags(ctx context.Context, name *biz.ArkOperatorName, tags *biz.ArkRecruitTags) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *myKratosRepo) GetRecruitTags(ctx context.Context, tagType string) (*biz.ArkRecruitTags, error) {
	if result, err := r.data.rdb.SMembers(ctx, tagType+":tagName:recruit").Result(); err != nil {
		return nil, err
	} else {
		return &biz.ArkRecruitTags{Tags: result}, nil
	}
}

// GetAllRecruitOperatorName 实现从redis数据库中获取公开招募干员名,返回所有公招干员名列表 只出现在公开招募干员map 其他公招干员map
func (r *myKratosRepo) GetAllRecruitOperatorName(ctx context.Context) ([]*biz.ArkOperatorName, map[string]bool, map[string]bool, error) {
	onlyRecruitOperator, err := r.data.rdb.SMembers(ctx, "only_recruit_operator").Result()
	if err != nil {
		return nil, nil, nil, err
	}
	recruitOperatorNames := make([]*biz.ArkOperatorName, 0)
	onlyRecruitOperatorMap := make(map[string]bool)
	for _, name := range onlyRecruitOperator {
		recruitOperatorNames = append(recruitOperatorNames, &biz.ArkOperatorName{
			Name: name,
		})
		onlyRecruitOperatorMap[name] = true
	}
	notOnlyRecruitOperator, err := r.data.rdb.SMembers(ctx, "not_only_recruit_operator").Result()
	if err != nil {
		return nil, nil, nil, err
	}
	notOnlyRecruitOperatorMap := make(map[string]bool)
	for _, name := range notOnlyRecruitOperator {
		recruitOperatorNames = append(recruitOperatorNames, &biz.ArkOperatorName{
			Name: name,
		})
		notOnlyRecruitOperatorMap[name] = true
	}
	return recruitOperatorNames, onlyRecruitOperatorMap, notOnlyRecruitOperatorMap, nil
}
