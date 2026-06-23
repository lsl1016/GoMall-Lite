# GoMall Lite

GoMall Lite 是一个基于 Vue3 + Go + MySQL + Docker 的轻量级购物车商城项目。它覆盖登录认证、商品浏览、购物车管理、地址管理、订单提交、订单中心等核心电商流程，适合作为前后端分离项目的学习实战案例。

## 技术栈

前端：Vue 3、Vite、Vue Router、Pinia、Axios、原生 CSS。

后端：Go、Gin、GORM、MySQL、JWT、bcrypt。

部署：Docker、Docker Compose、Nginx。

## 项目结构

```txt
gomall-lite
├─ backend                # Go 后端 API
├─ frontend               # Vue3 前端
├─ database/init.sql      # MySQL 建表、索引、演示数据
├─ backend/cmd/cli        # Cobra 命令行工具
├─ scripts/init-db.cmd    # 调用 Go CLI 的 Windows 入口
├─ docker-compose.yml     # 一键启动前端 + 后端
└─ README.md
```

## 一键启动

确保本机已经安装 Docker 和 Docker Compose。

```bash
docker-compose up -d --build
```

首次创建 MySQL 数据卷时，`database/init.sql` 会自动建表、创建索引并写入演示数据。
如果 MySQL 容器或数据卷已经存在，可以手动重新执行初始化 SQL：

```bash
cd backend
go run ./cmd/cli db init
```

Windows 下也可以从项目根目录运行：

```powershell
.\scripts\init-db.cmd
```

如需清空当前 Docker 数据卷后重新初始化：

```powershell
.\scripts\init-db.cmd --reset
```

或使用 Go CLI：

```bash
cd backend
go run ./cmd/cli db init --reset
```

如果构建时无法连接 Docker Hub，例如拉取 `golang:1.22-alpine` 失败，可以在项目根目录创建 `.env` 并配置镜像源：

```env
MYSQL_IMAGE=docker.m.daocloud.io/library/mysql:8.0
GO_IMAGE=docker.m.daocloud.io/library/golang:1.22-alpine
ALPINE_IMAGE=docker.m.daocloud.io/library/alpine:3.20
NODE_IMAGE=docker.m.daocloud.io/library/node:20-alpine
NGINX_IMAGE=docker.m.daocloud.io/library/nginx:1.27-alpine
```

访问地址：

```txt
前端：http://localhost:5173
后端：http://localhost:8080
MySQL：localhost:3306
```

默认账号：

```txt
用户名：admin
密码：123456
```

## 后端分层

```txt
Router 层：注册路由、绑定参数、调用 Service、返回 JSON
Service 层：处理业务逻辑、封装 DTO、事务处理
Model 层：GORM 模型定义、MySQL CRUD 封装
DTO 层：请求 DTO、响应 DTO、统一返回结构
Middleware：CORS、JWT Auth
```

## 核心接口

```txt
POST   /api/register
POST   /api/login
GET    /api/user/info
GET    /api/products
GET    /api/products/:id
GET    /api/cart
POST   /api/cart
PUT    /api/cart/:id
DELETE /api/cart/:id
DELETE /api/cart
GET    /api/addresses
POST   /api/addresses
PUT    /api/addresses/:id
DELETE /api/addresses/:id
PUT    /api/addresses/:id/default
POST   /api/orders
GET    /api/orders
GET    /api/orders/:id
PUT    /api/orders/:id/pay
PUT    /api/orders/:id/cancel
```

## 本地开发

### 后端

```bash
cd backend
go mod tidy
go run ./cmd
```

默认读取环境变量。没有环境变量时会使用本地 MySQL 默认配置。

### 前端

```bash
cd frontend
npm install
npm run dev
```

前端开发服务器默认代理 `/api` 到 `http://localhost:8080`。

## 说明

用于学习vue框架
