package routers

import (
	"gin-synolux/common"
	admin "gin-synolux/controllers/admin"
	controllers "gin-synolux/controllers/api"
	"gin-synolux/db"
	"gin-synolux/service"

	"gin-synolux/middleware"

	"github.com/gin-gonic/gin"
)

// 初始化
func Init() *gin.Engine {
	router := gin.New()

	// ===== 依赖注入 =====
	articleService := service.NewArticleService(db.DB.Self)
	adService := service.NewAdService(db.DB.Self)
	feedbackService := service.NewFeedbackService(db.DB.Self)
	movieService := service.NewMovieService(db.DB.Self)
	userService := service.NewUserService(db.DB.Self)
	uploadService := service.NewUploadService(db.DB.Self)
	adminService := service.NewAdminService(db.DB.Self)

	//后台
	adminAd := admin.NewAdController(adService)
	adminArticle := admin.NewArticleController(articleService)
	adminMovie := admin.NewMovieController(movieService)
	adminUpload := admin.NewUploadController(uploadService)
	adminAdmin := admin.NewAdminController(adminService)

	//api
	apiArticle := controllers.NewArticleController(articleService)
	apiFeedback := controllers.NewFeedbackController(feedbackService)
	apiUser := controllers.NewUserController(userService)
	apiUpload := admin.NewUploadController(uploadService)

	// ===== 中间件 =====
	router.Use(middleware.Cors())
	router.Use(middleware.SetLocale())

	router.Static("/uploads", "./uploads")
	router.LoadHTMLGlob("views/*")

	// =========================
	// 前台 API
	// =========================
	api := router.Group("/api/v1")
	{
		// 公共接口
		api.GET("/captcha", new(controllers.CommonController).Captcha)
		api.GET("/ping", new(controllers.CommonController).Ping)
		api.POST("/upload", apiUpload.Upload)
		api.GET("/download", apiUpload.Download)
		api.POST("/feedback", apiFeedback.Save)
		api.GET("/chat_gpt", new(controllers.CommonController).ChatGPT)
		api.GET("/ip", new(controllers.CommonController).Ip)

		//用户
		api.POST("/login", apiUser.Login)
		api.POST("/register", apiUser.Register)

		// 文章（只读）
		api.GET("/article", apiArticle.Index)
		api.GET("/article/detail", apiArticle.Detail)
		api.GET("/home_article", apiArticle.HomeArticle)
	}

	// =========================
	// 前台登录用户
	// =========================
	apiAuth := router.Group("/api/v1")
	apiAuth.Use(middleware.AuthMiddleware(common.RoleUser))
	{
		//用户
		apiAuth.POST("/logout", apiUser.Logout)
		apiAuth.GET("/get_userinfo", apiUser.GetUserInfo) //获取用户信息
		apiAuth.POST("/profile", apiUser.Profile)
		apiAuth.POST("/set_password", apiUser.SetPassword)
	}

	// =========================
	// 后台 API
	// =========================
	adminAPI := router.Group("/admin_api/v1")
	{
		adminAPI.POST("/login", adminAdmin.Login)
		adminAPI.POST("/register", adminAdmin.Register)
		adminAPI.GET("/captcha", new(admin.CommonController).Captcha)
	}

	// =========================
	// 后台 API（必须登录）
	// =========================
	adminAPIAuth := router.Group("/admin_api/v1")
	adminAPIAuth.Use(middleware.AuthMiddleware(common.RoleAdmin))
	{
		// 文章管理（完整权限）
		adminAPIAuth.GET("/article", adminArticle.Index)
		adminAPIAuth.GET("/article/detail", adminArticle.Detail)
		adminAPIAuth.POST("/article/save", adminArticle.Save)
		adminAPIAuth.POST("/article/delete", adminArticle.Delete)
		adminAPIAuth.POST("/article/enable", adminArticle.Enable)
		adminAPIAuth.POST("/article/disable", adminArticle.Disable)

		// 广告管理（完整权限）
		adminAPIAuth.GET("/ad", adminAd.Index)
		adminAPIAuth.GET("/ad/detail", adminAd.Detail)
		adminAPIAuth.POST("/ad/save", adminAd.Save)
		adminAPIAuth.POST("/ad/delete", adminAd.Delete)
		adminAPIAuth.POST("/ad/enable", adminAd.Enable)
		adminAPIAuth.POST("/ad/disable", adminAd.Disable)
		
		//视频管理
		adminAPIAuth.GET("/movie", adminMovie.Index)
		adminAPIAuth.GET("/movie/detail", adminMovie.Detail)
		adminAPIAuth.POST("/movie/save", adminMovie.Save)
		adminAPIAuth.POST("/movie/delete", adminMovie.Delete)
		adminAPIAuth.POST("/movie/enable", adminMovie.Enable)
		adminAPIAuth.POST("/movie/disable", adminMovie.Disable)
		
		//用户
		adminAPIAuth.POST("/logout", adminAdmin.Logout)
		adminAPIAuth.GET("/get_userinfo", adminAdmin.GetUserInfo) //获取用户信息
		adminAPIAuth.POST("/profile", adminAdmin.Profile)
		adminAPIAuth.POST("/set_password", adminAdmin.SetPassword)

		// 其他
		adminAPIAuth.POST("/upload", adminUpload.Upload)
		adminAPIAuth.GET("/download", adminUpload.Download)

		adminAPIAuth.GET("/chat_gpt", new(admin.CommonController).ChatGPT)
		adminAPIAuth.GET("/ip", new(admin.CommonController).Ip)
		adminAPIAuth.GET("/ping", new(admin.CommonController).Ping)

		adminAPIAuth.GET("/test", new(admin.TestController).Test)
		adminAPIAuth.GET("/queue", new(admin.TestController).Queue)
		adminAPIAuth.POST("/send_msg", new(admin.CommonController).SendMsg)
	}

	return router
}