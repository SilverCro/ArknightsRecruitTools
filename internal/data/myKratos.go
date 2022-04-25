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
	//TODO implement me

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
	//TODO implement me
	r.log.Info("message","data:GetArkOperatorInfo start")
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
	r.log.Info("message","data:GetArkOperatorInfo doing")
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
		if v, err1 := strconv.Atoi(val); err1 != nil {
			panic(err1)
		} else {
			aoi.Rarity = int32(v)
		}

	}
	if val, err := r.data.rdb.SMembers(ctx, name.Name+":tagList:chara").Result(); err != nil {
		return nil, err
	} else {
		aoi.TarList = val
	}
	return aoi, nil
}

func (r *myKratosRepo) CheckOperatorTags(ctx context.Context, name *biz.ArkOperatorName, tags *biz.ArkRecruitTags) (bool, error) {
	//TODO implement me
	panic("implement me")
}
