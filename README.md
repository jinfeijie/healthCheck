# health check
这是一个简单的健康检查脚本

用于在服务不可用情况下的报警通知

前往[http://push.strcpy.cn/](http://push.strcpy.cn/)注册一个账号，获取到jwt token就可以开启无缝通知

另外配置支持热加载（注意校验配置文件，避免监控意外退出）。因此，你可以使用配置中心对配置文件进行变更。

```json
{
  "mod": "debug",
  "app_name": "MJ-Health",
  "check_interval": 10,  // 每次请求的间隔 10 * 100ms （单位100ms）
  "email_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODc2NTUyOTksImlkIjoxLCJvcmlnX2lhdCI6MTU4NTA2MzI5OSwidXNlcm5hbWUiOiJmaiBqaW4ifQ.iVUkeyowxQMhaU8CYEqdB6PPAIkiihBe9iVjS9rFYmc",
  "email_server": "http://push.strcpy.cn/mail/free_send",
  "site": [
    {
      "name": "site name",
      "domain": "https://domain.com/ping",  // 检测网址
      "notify": [
        "your@email.com"  // 通知的邮箱
      ],
      "notify_interval": 600,
      "notify_format": "%s <br> 网站：%s 异常<br>时间：%s",
      "result": "pong"  // 响应内容
    }
  ]
}
```