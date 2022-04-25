package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// 实现从公招字符串提取出以逗号分隔的干员姓名 参数s为输入公招字符串 格式如" 干员a / 干员b / 干员c "
func GetStringList(s string) {
	ss := strings.Split(s, "/")
	result := ""
	for _, val := range ss {
		if val != "" {
			result += "\"" + strings.TrimSpace(val) + "\"" + ","
		}
	}
	fmt.Println(result)
}

// 用于读取json的公招干员结构体
type RecruitOperator struct {
	OnlyRecruitOperator    []string `json:"only_recruit_operator"`
	NotOnlyRecruitOperator []string `json:"not_only_recruit_operator"`
}

// 从公招干员json文件中读取不同干员姓名 参数fileName为公招干员姓名json文件路径例如（arkdata/recruit_operator.json） 返回值是读取出的干员结构体
func readJsonRecruitData(fileName string) *RecruitOperator {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("error reading json file")
		return nil
	}
	var recruitOperator RecruitOperator
	json.Unmarshal(jsonData, &recruitOperator)
	//fmt.Println(recruitOperator)
	return &recruitOperator
}

// 用于将读出的公招干员姓名写入redis数据库中 参数rdb为redis客户端,recruitOperator为公招干员数据
func writeRecruitOperator2Redis(rdb *redis.Client, recruitOperator *RecruitOperator) {
	ctx := context.Background()
	fmt.Println(rdb)
	// Redis<110.40.221.138:6880 db:0>
	for _, val := range recruitOperator.OnlyRecruitOperator {
		// 执行命令
		err := rdb.SAdd(ctx, "only_recruit_operator", val).Err()
		if err != nil {
			panic(err)
		}
	}
	for _, val := range recruitOperator.NotOnlyRecruitOperator {
		// 执行命令
		err := rdb.SAdd(ctx, "not_only_recruit_operator", val).Err()
		if err != nil {
			panic(err)
		}
	}

	val1, err1 := rdb.SCard(ctx, "only_recruit_operator").Result()
	val2, err2 := rdb.SCard(ctx, "not_only_recruit_operator").Result()
	if err1 != nil {
		panic(err1)
	}
	if err2 != nil {
		panic(err2)
	}
	fmt.Println("SCard only_recruit_operator : ", val1)     // set value
	fmt.Println("SCard not_only_recruit_operator : ", val2) // set value
}

// 用于从redis中读取公招干员数据 参数rdb为redis客户端
func readRecruitOperatorFromRedis(rdb *redis.Client) {
	ctx := context.Background()
	fmt.Println(rdb)
	// Redis<110.40.221.138:6880 db:0>
	val1, err1 := rdb.SMembers(ctx, "only_recruit_operator").Result()
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(val1)
}

// 用于存储从json文件中读出的干员信息结构体
type ArkCharaData struct {
	name       string
	tagList    []string
	profession string
	position   string
	rarity     int
}

// 用于从json文件中读取干员数据 参数fileName为干员数据json文件路径例如（arkdata/character_table.json） 返回值为[]*ArkCharaData干员数据数组
func readJsonCharData(fileName string) []*ArkCharaData {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("error reading json file")
		return nil
	}
	fmt.Println(len(jsonData))
	var p interface{}
	json.Unmarshal(jsonData, &p)
	m := p.(map[string]interface{})
	arkCharaDataList := make([]*ArkCharaData, 0)
	for key := range m {
		chara := m[key]
		m1 := chara.(map[string]interface{})
		name := m1["name"].(string)
		fmt.Println(name)
		tagTempList := make([]interface{}, 0)
		if m1["tagList"] != nil {
			tagTempList = m1["tagList"].([]interface{})
		}
		tagList := make([]string, len(tagTempList))
		for index, val := range tagTempList {
			tagList[index] = val.(string)
		}
		profession := m1["profession"].(string)
		position := m1["position"].(string)
		rarityTemp := m1["rarity"].(float64)
		rarity := int(rarityTemp)
		arkCharaData := &ArkCharaData{
			name:       name,
			tagList:    tagList,
			position:   position,
			rarity:     rarity,
			profession: profession,
		}
		arkCharaDataList = append(arkCharaDataList, arkCharaData)
	}
	return arkCharaDataList
}

// 用于将读出的干员数据存入redis数据库中 参数rdb为redis客户端,arkCharaDataList为读出的干员数据数组
func writeCharaData2Redis(rdb *redis.Client, arkCharaDataList []*ArkCharaData) {
	ctx := context.Background()
	fmt.Println(rdb)
	// Redis<localhost:6379 db:0>

	for _, val := range arkCharaDataList {
		err := rdb.SAdd(ctx, "name:chara", val.name).Err()
		if err != nil {
			panic(err)
		}
		err = rdb.Set(ctx, val.name+":rarity:chara", val.rarity+1, 0).Err()
		if err != nil {
			panic(err)
		}
		err = rdb.Set(ctx, val.name+":profession:chara", val.profession, 0).Err()
		if err != nil {
			panic(err)
		}
		for _, v := range val.tagList {
			err = rdb.SAdd(ctx, val.name+":tagList:chara", v).Err()
			if err != nil {
				panic(err)
			}
		}
		if val.position != "" {
			if val.position == "MELEE" {
				err = rdb.SAdd(ctx, val.name+":tagList:chara", "近战位").Err()
				if err != nil {
					panic(err)
				}
			} else if val.position == "RANGED" {
				err = rdb.SAdd(ctx, val.name+":tagList:chara", "远程位").Err()
				if err != nil {
					panic(err)
				}
			}
		}
	}
	val1, err1 := rdb.Keys(ctx, "*").Result()
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(val1)

}

// 修正干员职业数据 参数rdb为redis客户端
func modifyCharaprofession2Redis(rdb *redis.Client) {
	ctx := context.Background()
	fmt.Println(rdb)
	nameList := []string{}
	if val, err := rdb.SMembers(ctx, "name:chara").Result(); err != nil {
		panic(err)
	} else {
		nameList = val
	}
	m := make(map[string]string)
	m["SUPPORT"] = "辅助干员"
	m["WARRIOR"] = "近卫干员"
	m["SNIPER"] = "狙击干员"
	m["MEDIC"] = "医疗干员"
	m["PIONEER"] = "先锋干员"
	m["TANK"] = "重装干员"
	m["SPECIAL"] = "特种干员"
	m["CASTER"] = "术士干员"
	for _, val := range nameList {
		//if err := rdb.SRem(ctx, val+":tagList:chara", 0).Err(); err != nil {
		//	panic(err)
		//}
		if v, err := rdb.Get(ctx, val+":profession:chara").Result(); err != nil {
			panic(err)
		} else {
			if vv, ok := m[v]; ok {
				if err = rdb.Set(ctx, val+":profession:chara", vv, 0).Err(); err != nil {
					panic(err)
				}
				if err = rdb.SAdd(ctx, val+":tagList:chara", vv).Err(); err != nil {
					panic(err)
				}
			}

		}
	}

}

// 为干员tarList添加稀有度标签 参数rdb为redis客户端
func appendCharaTagList2Redis(rdb *redis.Client) {
	ctx := context.Background()
	fmt.Println(rdb)
	nameList := []string{}
	if val, err := rdb.SMembers(ctx, "name:chara").Result(); err != nil {
		panic(err)
	} else {
		nameList = val
	}
	m := make(map[int]string)
	m[1] = "机械支援"
	m[2] = "新手"
	m[5] = "资深干员"
	m[6] = "高级资深干员"
	for _, val := range nameList {
		//if err := rdb.SRem(ctx, val+":tagList:chara", 0).Err(); err != nil {
		//	panic(err)
		//}
		if v, err := rdb.Get(ctx, val+":rarity:chara").Result(); err != nil {
			panic(err)
		} else {
			stars := 3
			if vv, err1 := strconv.Atoi(v); err1 != nil {
				panic(err1)
			} else {
				stars = vv
			}
			if vv, ok := m[stars]; ok {
				if err = rdb.SAdd(ctx, val+":tagList:chara", vv).Err(); err != nil {
					panic(err)
				}
			}

		}
	}

}

// 用于存储从redis数据库中读出的干员数据
type ArkCharaReadData struct {
	name       string
	tagList    []string
	profession string
	rarity     int
}

// 从redis数据库中读取某一干员的数据 参数rdb为redis客户端,searchName为查询干员的名字
func readOneCharaDataFromRedis(rdb *redis.Client, searchName string) *ArkCharaReadData {
	ctx := context.Background()
	fmt.Println(rdb)
	// Redis<110.40.221.138:6880 db:0>

	if val, err := rdb.SIsMember(ctx, "name:chara", searchName).Result(); err != nil {
		panic(err)
	} else {
		if !val {
			fmt.Println("Cannot find the operator")
			return nil
		}
	}

	arkCharaReadData := &ArkCharaReadData{
		name: searchName,
	}

	if val, err := rdb.Get(ctx, searchName+":profession:chara").Result(); err != nil {
		panic(err)
	} else {
		arkCharaReadData.profession = val
	}

	if val, err := rdb.Get(ctx, searchName+":rarity:chara").Result(); err != nil {
		panic(err)
	} else {
		if v, err1 := strconv.Atoi(val); err1 != nil {
			panic(err1)
		} else {
			arkCharaReadData.rarity = v
		}

	}

	if val, err := rdb.SMembers(ctx, searchName+":tagList:chara").Result(); err != nil {
		panic(err)
	} else {
		arkCharaReadData.tagList = val
	}
	return arkCharaReadData
}

// 将tag数据存入redis数据库中 参数rdb为redis客户端
//func writeTags2Redis(rdb *redis.Client)  {
//	m := make(map[string]string)
//	m[]
//	tagStr := "控场 爆发 治疗 支援 费用回复 输出 生存 群攻 防护 减速 削弱 快速复活 位移 召唤 支援机械 " +
//	"近卫干员 狙击干员 重装干员 医疗干员 辅助干员 术士干员 特种干员 先锋干员 " +
//	"近战位 远程位 " +
//	"新手 资深干员 高级资深干员"
//	tagStrs := strings.Split(tagStr, " ")
//	for _,val := tagStrs
//
//}
func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "110.40.221.138:6880",
		Password: "19681112gzp8", // no password set
		DB:       0,              // use default DB
	})

	//modifyCharaprofession2Redis(rdb)
	//appendCharaTagList2Redis(rdb)
	var input string
	for {
		fmt.Scanf("%s\n", &input)
		if input == "end" {
			break
		}
		fmt.Println(input)
		arkCharaReadData := readOneCharaDataFromRedis(rdb, input)
		fmt.Println(arkCharaReadData)
	}

	//arkCharaReadData := readOneCharaDataFromRedis(rdb, "陈")
	//fmt.Println(arkCharaReadData)
}
