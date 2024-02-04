#  明日方舟公开招募辅助工具（开发中）

## 项目描述
是本项目基于[kratos微服务框架](https://github.com/go-kratos/kratos) 开发，利用redis存储明日方舟干员和公开招募数据，实现公开招募tag选择推荐


## 目前完成的工作
#### 1.redis数据库数据上传
#### 2.干员信息数据获取接口
接口访问方式：`get`
接口地址 `{{APIURL}}/api/ark/tools/operator/{{干员名}}`

#### 3.公开招募推荐接口

接口访问方式：`post`  
接口地址 `{{APIURL}}/api/ark/tools/recruit/recommend`  
接口request body格式`application/json`  
request body样例  

```json
{
    "tags" : [
        "近卫干员",
        "减速"
    ]
}
```

注：tags字段中包换的tag数量应该在5个以内且不应出现非法的tag，否则将返回推荐失败  
response body样例 

```json
{
    "status": "success 获取推荐tag成功",
    "recruitRecommends": [
        {
            "recommendTags": [
                "近卫干员",
                "减速"
            ],
            "recommendOperatorInfos": [
                {
                    "name": "霜叶",
                    "tarList": [
                        "近卫干员",
                        "输出",
                        "减速",
                        "近战位"
                    ],
                    "profession": "近卫干员",
                    "rarity": 4
                }
            ]
        }
    ]
}
```

注：status字段表示能否根据请求中的tag找到必定招募到四星以上或者一星的干员（除招募结果中tag被划到的以外），recruitRecommends字段表示推荐使用tag组合以及该组合能招募到的干员信息（可能有多种组合请根据需要选取），recommendTags字段为推荐的tag组合，recommendOperatorInfos字段为能够招募到的干员信息

## 待完成的工作
#### ~~1.公开招募推荐接口~~

#### 2.网站前端界面
## ~~在线测试地址~~（已到期）


