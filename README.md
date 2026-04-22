# 游戏服务器项目

## 项目简介

这是一个基于Go语言开发的游戏服务器项目，提供了完整的游戏后端基础设施，包括用户认证、WebSocket实时通信、排行榜、匹配系统等核心功能。

## 技术栈

- **Go 1.25.0** - 后端开发语言
- **Gin** - Web框架
- **GORM** - ORM框架
- **MySQL** - 关系型数据库
- **Redis** - 缓存服务
- **WebSocket** - 实时通信
- **JWT** - 身份认证
- **bcrypt** - 密码加密
- **Kafka** - 消息队列
- **Zap** - 日志系统

## 项目结构

```
game-server/
├── api/                # HTTP API接口
│   └── user.go         # 用户认证接口
├── config/             # 配置文件
│   ├── dev.yaml        # 开发环境配置
│   ├── test.yaml       # 测试环境配置
│   └── prod.yaml       # 生产环境配置
├── internal/           # 内部模块
│   ├── antiCheat/      # 防作弊系统
│   ├── chat/           # 聊天系统
│   ├── jwt/            # JWT认证
│   ├── match/          # 匹配系统
│   ├── msg/            # 消息定义
│   ├── rank/           # 排行榜系统
│   └── ws/             # WebSocket系统
├── logs/               # 日志文件
├── pkg/                # 公共包
│   ├── config/         # 配置管理
│   ├── kafka/          # Kafka消息队列
│   ├── logger/         # 日志系统
│   ├── mysql/          # MySQL数据库
│   └── redis/          # Redis缓存
├── Dockerfile          # Docker构建文件
├── docker-compose.yml  # Docker Compose配置
├── go.mod              # Go模块依赖
├── go.sum              # 依赖版本锁定
├── main.go             # 主入口文件
└── game-server.exe     # 编译后的可执行文件
```

## 安装步骤

### 1. 环境准备

- 安装 Go 1.25.0 或更高版本
- 安装 MySQL 8.0
- 安装 Redis
- 安装 Kafka (可选，用于匹配系统)

### 2. 克隆项目

```bash
git clone <项目地址>
cd game-server
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 配置数据库

- 创建 MySQL 数据库：`game_server`
- 配置 `config/test.yaml` 文件中的数据库连接信息

### 5. 构建项目

```bash
go build -o game-server.exe main.go
```

## 运行方法

### 开发环境

```bash
go run main.go -env=dev
```

### 测试环境

```bash
go run main.go -env=test
```

### 生产环境

```bash
go run main.go -env=prod
```

## API接口说明

### 用户认证

- **注册**：`POST /api/v1/register`
  - 请求体：`{"nickname": "用户名", "password": "密码"}`
  - 响应：`{"code": 0, "msg": "注册成功"}`

- **登录**：`POST /api/v1/login`
  - 请求体：`{"nickname": "用户名", "password": "密码"}`
  - 响应：`{"code": 0, "msg": "登录成功", "token": "JWT令牌"}`

### 排行榜

- **增加积分**：`GET /api/rank/add?uid=用户ID&score=积分`
- **获取排行榜**：`GET /api/rank/list`
- **获取个人排名**：`GET /api/rank/user?uid=用户ID`

### WebSocket

- **连接**：`GET /ws?token=JWT令牌`
- **消息格式**：`{"cmd": 指令ID, "data": "消息内容", "from": 发送者ID, "to": 接收者ID}`
- **指令定义**：
  - 1001：世界聊天
  - 1002：私聊
  - 2000：发起匹配
  - 2001：匹配成功

## 配置说明

### 配置文件结构

```yaml
app:
  name: game-server
  mode: test
  port: 8081

log:
  level: info
  path: ./logs/server.log
  max_size: 100
  max_backups: 10
  max_age: 7
  compress: true

mysql:
  host: localhost
  port: 3306
  username: root
  password: 密码
  database: game_server

redis:
  addr: localhost:6379
  password: ""
  db: 0

kafka:
  addr: localhost:9092
```

### 环境变量

- `GIN_MODE`：Gin框架模式（debug/release）
- `GO_ENV`：运行环境（dev/test/prod）

## 核心功能

### 1. 用户认证系统
- 注册、登录功能
- JWT令牌生成和验证
- 密码bcrypt加密

### 2. WebSocket系统
- 实时消息通信
- 玩家在线状态管理
- 世界聊天和私聊

### 3. 排行榜系统
- 基于Redis的有序集合实现
- 实时排名更新
- 支持获取前100名

### 4. 匹配系统
- 基于积分的匹配算法
- Kafka消息通知
- 匹配成功实时通知

### 5. 防作弊系统
- 接口频率限制
- 积分异常检测

## 部署方式

### 本地部署

1. 启动MySQL、Redis、Kafka服务
2. 运行 `game-server.exe -env=prod`

### Docker部署

```bash
docker-compose up -d
```

## 监控与日志

- 日志文件：`./logs/server.log`
- 日志级别：info/debug/warn/error
- 支持日志文件轮转

## 注意事项

1. 确保MySQL、Redis服务正常运行
2. Kafka服务可选，不影响基本功能
3. 生产环境建议使用HTTPS
4. 定期备份数据库

## 开发指南

### 代码规范
- 使用Go 1.25.0语法
- 遵循Go代码规范
- 使用zap进行日志记录

### 扩展建议
- 增加游戏核心逻辑
- 添加更多API接口
- 实现游戏房间系统
- 增加支付系统

## 许可证

MIT License

## 联系方式

如有问题，请联系项目维护者。