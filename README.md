#  明日方舟公开招募辅助工具（开发中）

## 项目描述
是本项目基于[kratos微服务框架](https://github.com/go-kratos/kratos) 开发，利用redis存储明日方舟干员和公开招募数据，实现公开招募tag选择推荐


## 目前完成的工作
#### 1.redis数据库数据上传
#### 2.干员信息数据获取接口
接口访问方式：`get`
接口地址 `{{APIURL}}/api/ark/tools/operator/{{干员名}}`

## 待完成的工作
#### 1.公开招募推荐接口
#### 2.网站前端界面
## 在线测试地址
sayHello测试接口：[`http://110.40.221.138:8000/helloworld/{{name}}`](http://110.40.221.138:8000/helloworld/kratos) 

干员信息接口：[`http://110.40.221.138:8000/api/ark/tools/operator/{{干员名}}`](http://110.40.221.138:8000/api/ark/tools/operator/%E9%93%B6%E7%81%B0)

