为我的 **GoMall Lite** 项目最适合加一套“分层清晰、够真实但不复杂”的日志系统。

建议从 4 个层面加日志：

```txt
1. 后端 Go 日志：记录请求、业务错误、数据库错误
2. Gin 中间件日志：记录每一次 API 请求
3. GORM SQL 日志：记录慢 SQL / 数据库错误
4. Docker 日志：容器部署后方便查看日志
```

前端也可以加轻量日志，但重点应该放在后端。

---

# 一、推荐日志方案

后端推荐用 Go 标准库：

```txt
log/slog
```

优点：

```txt
Go 标准库自带
支持结构化日志
不需要额外依赖
适合教学项目和中小项目
Docker 里也方便查看
```

日志格式建议用 JSON：

```json
{
  "time": "2026-06-24T10:00:00Z",
  "level": "INFO",
  "msg": "request completed",
  "method": "GET",
  "path": "/api/products",
  "status": 200,
  "latency_ms": 12
}
```

---

# 二、后端目录新增 logger 包

你的后端可以新增：

```txt
backend
└─ internal
   ├─ logger
   │  └─ logger.go
   ├─ middleware
   │  ├─ auth.go
   │  ├─ cors.go
   │  ├─ request_logger.go
   │  └─ recovery.go
```

---

# 三、新增全局日志初始化

新建：

```txt
backend/internal/logger/logger.go
```

```go
package logger

import (
	"log/slog"
	"os"
	"strings"
)

var Log *slog.Logger

func Init() {
	level := parseLevel(os.Getenv("LOG_LEVEL"))
	format := os.Getenv("LOG_FORMAT")

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler

	if strings.ToLower(format) == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
```

---

# 四、在 main.go 初始化日志

例如：

```go
package main

import (
	"gomall-lite/internal/logger"
	"gomall-lite/internal/model"
	"gomall-lite/internal/router"
)

func main() {
	logger.Init()

	logger.Log.Info("starting gomall-lite api")

	model.InitDB()

	r := router.SetupRouter()

	if err := r.Run(":8080"); err != nil {
		logger.Log.Error("server start failed", "error", err)
	}
}
```

启动时会输出：

```json
{"time":"2026-06-24T10:00:00Z","level":"INFO","msg":"starting gomall-lite api"}
```

---

# 五、添加 Gin 请求日志中间件

新建：

```txt
backend/internal/middleware/request_logger.go
```

```go
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"gomall-lite/internal/logger"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		logger.Log.Info(
			"request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"status", c.Writer.Status(),
			"latency_ms", latency.Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)
	}
}
```

然后在 router 里使用。

```go
func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())

	api := r.Group("/api")

	// routes...

	return r
}
```

这样每次请求都会自动记录日志：

```json
{
  "level": "INFO",
  "msg": "request completed",
  "method": "GET",
  "path": "/api/products",
  "status": 200,
  "latency_ms": 18,
  "client_ip": "172.18.0.1"
}
```

---

# 六、添加 panic recovery 日志

新建：

```txt
backend/internal/middleware/recovery.go
```

```go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gomall-lite/internal/logger"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Log.Error(
			"panic recovered",
			"error", recovered,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP(),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务器内部错误",
			"data":    nil,
		})
	})
}
```

这样后端即使 panic，也不会直接崩掉，而且能记录错误原因。

---

# 七、在 Service 层添加业务日志

Service 层是最值得打日志的地方。

例如登录：

```go
func Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	logger.Log.Info("user login attempt", "username", req.Username)

	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		logger.Log.Warn("user login failed: user not found", "username", req.Username)
		return nil, errors.New("用户名或密码错误")
	}

	if user.Password != req.Password {
		logger.Log.Warn("user login failed: wrong password", "username", req.Username)
		return nil, errors.New("用户名或密码错误")
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		logger.Log.Error("generate token failed", "user_id", user.ID, "error", err)
		return nil, err
	}

	logger.Log.Info("user login success", "user_id", user.ID, "username", user.Username)

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
		},
	}, nil
}
```

注意：**不要打印密码。**

不要这样：

```go
logger.Log.Info("login", "password", req.Password)
```

---

# 八、购物车 Service 日志示例

```go
func AddToCart(userID uint, productID uint, count int) error {
	logger.Log.Info(
		"add to cart",
		"user_id", userID,
		"product_id", productID,
		"count", count,
	)

	product, err := model.GetProductByID(productID)
	if err != nil {
		logger.Log.Warn(
			"add to cart failed: product not found",
			"user_id", userID,
			"product_id", productID,
		)
		return errors.New("商品不存在")
	}

	if product.Stock < count {
		logger.Log.Warn(
			"add to cart failed: stock not enough",
			"user_id", userID,
			"product_id", productID,
			"stock", product.Stock,
			"count", count,
		)
		return errors.New("库存不足")
	}

	err = model.UpsertCartItem(userID, productID, count)
	if err != nil {
		logger.Log.Error(
			"add to cart db failed",
			"user_id", userID,
			"product_id", productID,
			"error", err,
		)
		return err
	}

	logger.Log.Info(
		"add to cart success",
		"user_id", userID,
		"product_id", productID,
	)

	return nil
}
```

---

# 九、订单 Service 日志示例

订单是核心链路，建议重点记录。

```go
func CreateOrder(userID uint, req dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	logger.Log.Info("create order start", "user_id", userID, "address_id", req.AddressID)

	order, err := model.CreateOrderWithItems(userID, req)
	if err != nil {
		logger.Log.Error(
			"create order failed",
			"user_id", userID,
			"address_id", req.AddressID,
			"error", err,
		)
		return nil, err
	}

	logger.Log.Info(
		"create order success",
		"user_id", userID,
		"order_id", order.ID,
		"order_no", order.OrderNo,
		"total_amount", order.TotalAmount,
	)

	return &dto.CreateOrderResponse{
		OrderNo: order.OrderNo,
	}, nil
}
```

---

# 十、给请求加 request_id

实际项目中最好每个请求都有一个 `request_id`，方便排查问题。

新增：

```txt
backend/internal/middleware/request_id.go
```

```go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
```

安装依赖：

```bash
cd backend
go get github.com/google/uuid
```

在 router 中：

```go
r.Use(middleware.RequestID())
r.Use(middleware.RequestLogger())
```

然后改造 `RequestLogger`：

```go
requestID, _ := c.Get("request_id")

logger.Log.Info(
	"request completed",
	"request_id", requestID,
	"method", c.Request.Method,
	"path", c.Request.URL.Path,
	"status", c.Writer.Status(),
	"latency_ms", latency.Milliseconds(),
)
```

这样你能通过 `request_id` 串起一整次请求链路。

---

# 十一、GORM 添加 SQL 日志

你的项目用 MySQL + GORM，建议打开慢 SQL 日志。

在数据库初始化的地方，比如：

```txt
backend/internal/model/db.go
```

可以这样配置：

```go
package model

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_DSN")

	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	DB = db
}
```

推荐：

```txt
开发环境：GORM LogLevel = Info
测试环境：GORM LogLevel = Warn
生产环境：GORM LogLevel = Error / Warn
```

可以用环境变量控制：

```go
func parseGormLogLevel() gormlogger.LogLevel {
	switch os.Getenv("GORM_LOG_LEVEL") {
	case "info":
		return gormlogger.Info
	case "error":
		return gormlogger.Error
	case "silent":
		return gormlogger.Silent
	default:
		return gormlogger.Warn
	}
}
```

---

# 十二、环境变量配置

在 `.env.example` 里增加：

```env
LOG_LEVEL=info
LOG_FORMAT=json
GORM_LOG_LEVEL=warn
```

在 `docker-compose.yml` 的 backend 服务中增加：

```yaml
backend:
  build: ./backend
  environment:
    DB_DSN: gomall:gomall123@tcp(mysql:3306)/gomall_lite?charset=utf8mb4&parseTime=True&loc=Local
    JWT_SECRET: gomall-lite-secret
    LOG_LEVEL: info
    LOG_FORMAT: json
    GORM_LOG_LEVEL: warn
```

---

# 十三、Docker 中怎么看日志

启动项目：

```bash
docker compose up -d
```

查看后端日志：

```bash
docker compose logs -f backend
```

查看最近 100 行：

```bash
docker compose logs --tail=100 backend
```

查看 MySQL 日志：

```bash
docker compose logs -f mysql
```

---

# 十四、Docker 日志轮转

避免日志把磁盘打满，可以在 `docker-compose.yml` 里加：

```yaml
backend:
  build: ./backend
  logging:
    driver: "json-file"
    options:
      max-size: "10m"
      max-file: "3"
```

前端和 MySQL 也可以加：

```yaml
frontend:
  build: ./frontend
  logging:
    driver: "json-file"
    options:
      max-size: "10m"
      max-file: "3"

mysql:
  image: mysql:8.0
  logging:
    driver: "json-file"
    options:
      max-size: "10m"
      max-file: "3"
```

---

# 十五、前端 Axios 日志

前端不要打太多日志，主要在开发环境记录请求即可。

修改：

```txt
frontend/src/utils/request.js
```

```js
import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  if (import.meta.env.DEV) {
    console.log('[request]', config.method?.toUpperCase(), config.url, config.params || config.data)
  }

  return config
})

request.interceptors.response.use(
  (response) => {
    if (import.meta.env.DEV) {
      console.log('[response]', response.config.url, response.data)
    }

    return response.data
  },
  (error) => {
    if (import.meta.env.DEV) {
      console.error('[response error]', error)
    }

    return Promise.reject(error)
  }
)

export default request
```

注意：前端也不要打印密码、token 等敏感信息。

可以过滤：

```js
function maskSensitive(data) {
  if (!data || typeof data !== 'object') return data

  const copy = { ...data }

  if (copy.password) copy.password = '******'
  if (copy.token) copy.token = '******'

  return copy
}
```

---

# 十六、哪些地方应该打日志？

## 必须打

```txt
服务启动成功 / 失败
数据库连接成功 / 失败
登录成功 / 失败
JWT 解析失败
商品不存在
库存不足
加入购物车失败
创建订单成功 / 失败
支付成功 / 失败
取消订单成功 / 失败
接口 panic
数据库错误
```

## 不建议打

```txt
用户密码
完整 token
身份证 / 银行卡 / 敏感隐私
过大的请求体
每一行普通循环数据
```

---

# 十七、日志级别怎么用？

```txt
DEBUG：开发调试信息，比如参数、临时状态
INFO：正常业务流程，比如登录成功、创建订单成功
WARN：业务异常，比如库存不足、密码错误、资源不存在
ERROR：系统异常，比如数据库错误、panic、token 生成失败
```

例子：

```go
logger.Log.Debug("query params", "category", category, "keyword", keyword)

logger.Log.Info("order created", "order_no", orderNo, "user_id", userID)

logger.Log.Warn("stock not enough", "product_id", productID, "stock", stock)

logger.Log.Error("database error", "error", err)
```

---


# 十八、日志格式规范

如果用 slog，则优先用结构化字段：

logger.Log.Info(
	"create order success",
	"order_id", order.ID,
	"order_no", order.OrderNo,
	"total_amount", order.TotalAmount,
)

如果是普通格式化日志，就用：

logger.Log.Info(fmt.Sprintf("create order success order_id:%d order_no:%s total_amount:%d", order.ID, order.OrderNo, order.TotalAmount))

项目里建议统一规则：

普通 fmt/log 输出：用 xxx:%d / xxx:%s 格式化
结构化 slog 输出：用 key-value 字段
不要使用字符串 + 变量直接拼接
```

# 十九、最终推荐落地顺序

你可以按这个顺序加：

```txt
1. backend/internal/logger/logger.go
2. main.go 初始化 logger
3. middleware/RequestLogger
4. middleware/Recovery
5. docker-compose.yml 增加 LOG_LEVEL / LOG_FORMAT
6. Service 层关键业务加日志
7. GORM 慢 SQL 日志
8. Docker logging 日志轮转
9. 前端 Axios 开发环境日志
10. request_id 链路追踪
```

