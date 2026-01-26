# Shop SaaS Base Project

这是一个基于 **Gin + Uber Fx + GORM** 构建的现代化电商 SaaS 基础项目。集成了微服务开发所需的各类基础设施组件，采用模块化设计，开箱即用。

## 🚀 技术栈

### 核心框架
- **Web 框架**: [Gin](https://github.com/gin-gonic/gin)
- **依赖注入**: [Uber Fx](https://go.uber.org/fx)
- **ORM**: [GORM](https://gorm.io/) (默认 MySQL)
- **配置管理**: [Viper](https://github.com/spf13/viper)
- **日志系统**: [Zap](https://go.uber.org/zap)

### 基础设施 & 中间件
- **缓存**: [Redis](https://github.com/redis/go-redis)
- **消息队列**: [RabbitMQ](https://github.com/rabbitmq/amqp091-go)
- **搜索引擎**: [Elasticsearch](https://github.com/elastic/go-elasticsearch)
- **实时通信**: [WebSocket](https://github.com/gorilla/websocket)
- **定时任务**: [Cron](https://github.com/robfig/cron)
- **ID 生成**: [Snowflake](https://github.com/bwmarrin/snowflake) (Twitter 雪花算法)

## 📂 目录结构

```
├── cmd
│   └── server              # 程序入口
├── configs                 # 配置文件
├── internal
│   ├── bootstrap           # Fx 应用组装与生命周期管理
│   ├── config              # 配置加载逻辑
│   ├── cron                # 定时任务管理
│   ├── database            # 数据库连接 (GORM)
│   ├── handler             # HTTP 控制层
│   ├── infra               # 基础设施客户端
│   │   ├── elasticsearch   # ES 客户端
│   │   ├── rabbitmq        # RabbitMQ 客户端
│   │   └── redis           # Redis 客户端
│   ├── middleware          # Gin 中间件 (CORS, Auth 等)
│   ├── model               # 数据模型
│   ├── repository          # 数据访问层 (DAO)
│   ├── router              # 路由定义 (SaaS, Admin, Mall)
│   ├── server              # HTTP Server 配置
│   ├── service             # 业务逻辑层
│   └── websocket           # WebSocket Hub & Client
└── pkg
    ├── idgen               # 分布式唯一 ID 生成器
    ├── logger              # 日志工具
    └── utils               # 通用工具 (Crypto, Random 等)
```

## 🧩 模块划分

路由层已预置了三个核心业务模块，分别应对不同的业务场景：

1.  **SaaS 管理端** (`/api/saas`)
    *   面向平台超级管理员。
    *   用于管理租户、计费套餐、系统全局配置等。
2.  **电商后台** (`/api/admin`)
    *   面向商家/租户管理员。
    *   用于管理商品、订单、会员、营销活动、店铺装修等。
3.  **电商前台** (`/api/mall`)
    *   面向 C 端消费者 (App/小程序/H5)。
    *   提供商品浏览、购物车、下单支付、个人中心等接口。

## 🛠️ 快速开始

### 1. 环境准备

确保本地或开发环境已安装以下服务：
*   MySQL
*   Redis
*   RabbitMQ
*   Elasticsearch (可选)

### 2. 配置文件

修改 `configs/config.yaml`，配置相关连接信息：

```yaml
server:
  port: ":8080"
  mode: "debug"
  node_id: 1 # Snowflake 节点 ID (0-1023)

database:
  dsn: "root:password@tcp(127.0.0.1:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"

elasticsearch:
  addresses: 
    - "http://localhost:9200"
```

### 3. 运行项目

```bash
go mod tidy
go run cmd/server/main.go
```

### 4. 接口测试

*   **HTTP API**:
    *   SaaS 健康检查: `GET http://localhost:8080/api/saas/health`
    *   注册用户: `POST http://localhost:8080/api/mall/register`
*   **WebSocket**:
    *   连接地址: `ws://localhost:8080/ws`

## 📖 开发指南

### 添加新 API

1.  **定义 Model**: 在 `internal/model` 中定义数据结构。
2.  **Repository**: 在 `internal/repository` 实现数据访问接口。
3.  **Service**: 在 `internal/service` 实现业务逻辑。
4.  **Handler**: 在 `internal/handler` 处理 HTTP 请求。
5.  **注册**:
    *   在 `internal/bootstrap/app.go` 的 `fx.Provide` 中注册新的 Repo, Service, Handler。
    *   在 `internal/router/router.go` 的对应模块 (`registerMallRoutes` 等) 中添加路由。

### 使用工具组件

*   **生成唯一 ID**:
    ```go
    // 注入 idgen.IDGenerator
    id := idGen.GenerateID() // int64
    ```
*   **定时任务**:
    在 `internal/cron/cron.go` 的 `RegisterJobs` 中添加任务。
*   **WebSocket 广播**:
    注入 `*websocket.Hub` 并调用相关方法（需自行实现广播接口或通过 channel 发送）。

## 📄 License

MIT
