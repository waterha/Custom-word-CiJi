# 词迹 (CiJi) - Word Cultivator

<p align="center">
  <a href="https://github.com/waterha/CiJi" target="_blank">
    <img src="https://img.shields.io/badge/GitHub-waterha/CiJi-blue?logo=github" alt="GitHub">
  </a>
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js" alt="Vue 3">
  <img src="https://img.shields.io/badge/MySQL-9.0-4479A1?logo=mysql" alt="MySQL">
  <img src="https://img.shields.io/badge/Redis-7.x-DC382D?logo=redis" alt="Redis">
  <img src="https://img.shields.io/badge/Docker-✓-2496ED?logo=docker" alt="Docker">
</p>

英语单词在线学习平台，支持管理员管理单词库、用户学习跟踪、学习数据统计可视化。

[GitHub 仓库](https://github.com/waterha) | [在线 Demo](http://8.134.108.251)

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + Vite + Pinia + Vue Router + Axios + ECharts |
| 后端 | Go 1.22 + Gin + GORM + JWT |
| 数据库 | MySQL 9.0 |
| 缓存 | Redis 7.x |
| 部署 | Docker Compose + Nginx |

## 功能特性

- **用户认证** — 注册 / 登录，JWT 身份验证（7 天过期）
- **单词学习** — 随机展示单词，标记"认识/不认识"，进度持久化
- **自定义单词** — 用户可添加个人单词库独立学习
- **单词搜索** — 支持全文搜索，Redis 缓存加速
- **管理员后台** — 单词 CRUD 管理
- **监控面板** — 今日访问量、时段访问量折线图、每日注册量、错词排行（ECharts）
- **限流与缓存** — 接口限流、Redis 缓存、Nginx 代理缓存降级
- **健康检查** — MySQL/Redis 连接状态检测

## 快速开始（本地部署）

### 方式一：Docker Compose 部署（推荐）

**前置条件：** 安装 [Docker](https://docs.docker.com/get-docker/) 和 [Docker Compose](https://docs.docker.com/compose/install/)

```bash
# 1. 克隆项目
git clone https://github.com/waterha/CiJi.git && cd CiJi

# 2. 构建并启动所有服务
docker-compose up -d

# 3. 查看日志确认启动正常
docker-compose logs -f
```

服务启动后：
- **前端页面**：http://localhost （通过 Nginx 访问）
- **后端 API**：http://localhost:8080

### 方式二：本地手动部署（开发环境）

#### 前置条件

| 依赖 | 版本要求 | 验证命令 |
|------|---------|---------|
| Go | >= 1.22 | `go version` |
| Node.js | >= 18 | `node -v` |
| MySQL | >= 9.0 | `mysql --version` |
| Redis | >= 7.x | `redis-cli --version` |

#### 步骤

**1. 启动 MySQL 和 Redis**

确保 MySQL 和 Redis 服务已运行，创建数据库：

```sql
CREATE DATABASE IF NOT EXISTS xiuyanzhe CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

**2. 配置环境变量**

根据本地服务地址修改 `docker-compose.yml` 中的环境变量，或直接设置系统环境变量：

```bash
# 数据库连接（本地 MySQL 通常用 localhost）
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=123456
export MYSQL_DATABASE=xiuyanzhe

# Redis 连接
export REDIS_HOST=localhost
export REDIS_PORT=6379
```

**3. 启动后端**

```bash
# 在项目根目录
go mod tidy
go run main.go
```

后端默认监听 `:8080`，启动后会自动：
- 创建数据库表（AutoMigrate）
- 初始化管理员账号（admin）
- 添加示例单词
- 建立 MySQL 全文索引

**4. 启动前端（开发模式）**

```bash
cd frontend
npm install
npm run dev
```

前端开发服务器默认监听 `:5173`，已配置 `/api` 代理到 `localhost:8080`。

**访问地址：**
- **开发模式**：http://localhost:5173
- **生产构建**：`npm run build` → 输出到 `frontend/dist/`

### 方式三：Docker 多阶段构建（仅后端容器）

```bash
# 构建镜像
docker build -t xiuyanzhe:latest .

# 运行（需先启动 MySQL 和 Redis 容器）
docker run -d --name xiuyanzhe-app \
  -p 8080:8080 \
  -e MYSQL_HOST=host.docker.internal \
  -e MYSQL_PORT=3306 \
  -e MYSQL_USER=root \
  -e MYSQL_PASSWORD=123456 \
  -e MYSQL_DATABASE=xiuyanzhe \
  -e REDIS_HOST=host.docker.internal \
  -e REDIS_PORT=6379 \
  xiuyanzhe:latest
```

## 初始账号

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 管理员 | admin | admin |

普通用户可通过注册页面自行注册。

## 项目结构

```
xiuyanzhe/
├── main.go                 # 程序入口
├── go.mod / go.sum         # Go 依赖管理
├── Dockerfile              # Docker 多阶段构建（前端+后端）
├── docker-compose.yml      # Docker 编排（app + mysql + redis + nginx）
├── nginx.conf              # Nginx 反向代理配置（限流/缓存/安全头）
├── .env                    # 环境变量
│
├── controllers/            # 请求处理器
│   ├── auth.go             # 登录/注册
│   ├── learn.go            # 学习流程
│   ├── word.go             # 单词管理
│   ├── search.go           # 全文搜索
│   ├── stats.go            # 学习统计
│   ├── monitor.go          # 监控面板
│   └── custom_word.go      # 自定义单词
│
├── models/
│   └── models.go           # GORM 数据模型（User/Word/Progress/WrongWord等）
│
├── database/
│   └── database.go         # MySQL 连接池 + AutoMigrate + 初始化数据
│
├── redis/
│   ├── redis.go            # Redis 客户端
│   ├── cache.go            # 缓存操作
│   └── queue.go            # 异步队列
│
├── middleware/
│   ├── middleware.go        # JWT 认证 / CORS / 管理权限
│   ├── ratelimit.go        # 接口限流
│   ├── cache.go            # 响应缓存
│   ├── db.go               # 数据库中间件
│   └── circuitbreaker.go   # 熔断器
│
├── routes/
│   └── routes.go           # 路由注册
│
└── frontend/               # Vue 3 前端
    ├── src/
    │   ├── views/          # 页面组件（Login/Learn/Admin/Monitor等）
    │   ├── stores/         # Pinia 状态管理
    │   ├── router/         # 前端路由
    │   ├── api/            # Axios 封装
    │   ├── App.vue         # 根组件
    │   └── main.js         # 入口
    ├── index.html
    ├── vite.config.js      # Vite 配置（含 API 代理）
    └── package.json
```

## API 概览

### 公开接口
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/register` | 用户注册 |
| POST | `/api/login` | 用户登录 |
| GET | `/api/health` | 健康检查 |

### 学习接口（需认证）
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/learn/next` | 下一个学习单词 |
| POST | `/api/learn/answer` | 提交认识/不认识 |
| GET | `/api/learn/progress` | 学习进度 |
| GET | `/api/learn/search` | 单词搜索 |
| GET | `/api/learn/stats` | 学习统计 |

### 自定义单词（需认证）
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/custom/words` | 自定义单词列表 |
| POST | `/api/custom/words` | 添加自定义单词 |
| GET | `/api/custom/words/:id` | 查看自定义单词 |
| PUT | `/api/custom/words/:id` | 更新自定义单词 |
| DELETE | `/api/custom/words/:id` | 删除自定义单词 |

### 管理员接口（需 admin 权限）
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/words` | 单词列表 |
| POST | `/api/admin/words` | 新增单词 |
| PUT | `/api/admin/words/:id` | 更新单词 |
| DELETE | `/api/admin/words/:id` | 删除单词 |

### 监控接口（需 admin 权限）
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/monitor/overview` | 监控概览 |
| GET | `/api/monitor/hourly-visits` | 24 小时访问量 |
| GET | `/api/monitor/daily-registrations` | 7 天注册量 |

## 部署架构

```
                        ┌──────────────┐
                        │   Nginx:80   │
                        │ 限流/缓存/SSL │
                        └──────┬───────┘
                               │
                        ┌──────┴───────┐
                        │  Go App:8080 │
                        │  Gin + GORM  │
                        └──┬───────┬───┘
                           │       │
                    ┌──────┴┐  ┌───┴──────┐
                    │ MySQL │  │  Redis   │
                    │  9.0  │  │   7.x    │
                    └───────┘  └──────────┘
```

## 注意事项

1. **MySQL 连接**：本地开发时 `MYSQL_HOST` 需改为 `localhost`（默认值 `mysql` 是容器名）
2. **Redis 连接**：同理，`REDIS_HOST` 本地需改为 `localhost`
3. **首次启动**：自动建表 + 创建管理员 + 插入 10 个示例单词
4. **学习进度**：存储在 Redis + MySQL，通过异步队列批量持久化
5. **监控数据**：用户访问时自动记录 VisitLog，监控面板读取真实数据

## License

MIT
