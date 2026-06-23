package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gomall-lite-api/config"
	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/middleware"
	"gomall-lite-api/internal/service"
)

func SetupRouter(cfg config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	userSvc := service.NewUserService(cfg)
	productSvc := service.NewProductService()
	cartSvc := service.NewCartService()
	addressSvc := service.NewAddressService()
	orderSvc := service.NewOrderService()

	api := r.Group("/api")
	{
		api.POST("/register", func(c *gin.Context) {
			var req dto.RegisterRequest
			if !bind(c, &req) {
				return
			}
			data, err := userSvc.Register(req)
			respond(c, data, err, "注册成功")
		})

		api.POST("/login", func(c *gin.Context) {
			var req dto.LoginRequest
			if !bind(c, &req) {
				return
			}
			data, err := userSvc.Login(req)
			respond(c, data, err, "登录成功")
		})

		api.GET("/products", func(c *gin.Context) {
			category := c.Query("category")
			keyword := c.Query("keyword")
			data, err := productSvc.List(category, keyword)
			respond(c, data, err, "success")
		})

		api.GET("/products/:id", func(c *gin.Context) {
			id, ok := paramID(c)
			if !ok {
				return
			}
			data, err := productSvc.Detail(id)
			respond(c, data, err, "success")
		})

		auth := api.Group("")
		auth.Use(middleware.Auth(cfg))
		{
			auth.GET("/user/info", func(c *gin.Context) {
				data, err := userSvc.GetUserInfo(currentUserID(c))
				respond(c, data, err, "success")
			})

			auth.GET("/cart", func(c *gin.Context) {
				data, err := cartSvc.List(currentUserID(c))
				respond(c, data, err, "success")
			})

			auth.POST("/cart", func(c *gin.Context) {
				var req dto.AddCartRequest
				if !bind(c, &req) {
					return
				}
				data, err := cartSvc.Add(currentUserID(c), req)
				respond(c, data, err, "加入购物车成功")
			})

			auth.PUT("/cart/:id", func(c *gin.Context) {
				id, ok := paramID(c)
				if !ok {
					return
				}
				var req dto.UpdateCartRequest
				if !bind(c, &req) {
					return
				}
				data, err := cartSvc.Update(currentUserID(c), id, req)
				respond(c, data, err, "更新成功")
			})

			auth.DELETE("/cart/:id", func(c *gin.Context) {
				id, ok := paramID(c)
				if !ok {
					return
				}
				data, err := cartSvc.Remove(currentUserID(c), id)
				respond(c, data, err, "删除成功")
			})

			auth.DELETE("/cart", func(c *gin.Context) {
				err := cartSvc.Clear(currentUserID(c))
				respond(c, gin.H{}, err, "清空成功")
			})

			auth.GET("/addresses", func(c *gin.Context) {
				data, err := addressSvc.List(currentUserID(c))
				respond(c, data, err, "success")
			})

			auth.POST("/addresses", func(c *gin.Context) {
				var req dto.AddressRequest
				if !bind(c, &req) {
					return
				}
				data, err := addressSvc.Create(currentUserID(c), req)
				respond(c, data, err, "新增地址成功")
			})

			auth.PUT("/addresses/:id", func(c *gin.Context) {
				id, ok := paramID(c)
				if !ok {
					return
				}
				var req dto.AddressRequest
				if !bind(c, &req) {
					return
				}
				data, err := addressSvc.Update(currentUserID(c), id, req)
				respond(c, data, err, "更新地址成功")
			})

			auth.DELETE("/addresses/:id", func(c *gin.Context) {
				id, ok := paramID(c)
				if !ok {
					return
				}
				data, err := addressSvc.Delete(currentUserID(c), id)
				respond(c, data, err, "删除地址成功")
			})

			auth.PUT("/addresses/:id/default", func(c *gin.Context) {
				id, ok := paramID(c)
				if !ok {
					return
				}
				data, err := addressSvc.SetDefault(currentUserID(c), id)
				respond(c, data, err, "设置默认地址成功")
			})

			auth.POST("/orders", func(c *gin.Context) {
				var req dto.CreateOrderRequest
				if !bind(c, &req) {
					return
				}
				data, err := orderSvc.Create(currentUserID(c), req)
				respond(c, data, err, "下单成功")
			})

			auth.GET("/orders", func(c *gin.Context) {
				data, err := orderSvc.List(currentUserID(c))
				respond(c, data, err, "success")
			})

			auth.GET("/orders/:id", func(c *gin.Context) {
				id, ok := paramID(c)
				if !ok {
					return
				}
				data, err := orderSvc.Detail(currentUserID(c), id)
				respond(c, data, err, "success")
			})

			auth.PUT("/orders/:id/pay", func(c *gin.Context) {
				id, ok := paramID(c)
				if !ok {
					return
				}
				data, err := orderSvc.Pay(currentUserID(c), id)
				respond(c, data, err, "支付成功")
			})

			auth.PUT("/orders/:id/cancel", func(c *gin.Context) {
				id, ok := paramID(c)
				if !ok {
					return
				}
				data, err := orderSvc.Cancel(currentUserID(c), id)
				respond(c, data, err, "取消订单成功")
			})
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}

func bind(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return false
	}
	return true
}

func paramID(c *gin.Context) (uint, bool) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "ID 参数错误"))
		return 0, false
	}
	return uint(id64), true
}

func currentUserID(c *gin.Context) uint {
	value, exists := c.Get("userID")
	if !exists {
		return 0
	}
	id, _ := value.(uint)
	return id
}

func respond(c *gin.Context, data interface{}, err error, message string) {
	if err != nil {
		if appErr, ok := err.(*service.AppError); ok {
			status := http.StatusBadRequest
			if appErr.Code == 401 {
				status = http.StatusUnauthorized
			}
			if appErr.Code == 404 {
				status = http.StatusNotFound
			}
			c.JSON(status, dto.Fail(appErr.Code, appErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "服务器错误"))
		return
	}
	c.JSON(http.StatusOK, dto.MessageOK(message, data))
}
