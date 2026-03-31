package routers

import (
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

	//后台
	adminAd := admin.NewAdController(adService)
	adminArticle := admin.NewArticleController(articleService)
	adminMovie := admin.NewMovieController(movieService)
	adminUpload := admin.NewUploadController(uploadService)

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
	apiAuth.Use(middleware.AuthMiddleware())
	{
		//用户
		apiAuth.POST("/logout", apiUser.Logout)
		apiAuth.GET("/get_userinfo", apiUser.GetUserInfo) //获取用户信息
		apiAuth.POST("/profile", apiUser.Profile)
		apiAuth.POST("/set_password", apiUser.SetPassword)

		apiAuth.POST("/upload", apiUpload.Upload)
		apiAuth.GET("/download", apiUpload.Download)

		apiAuth.POST("/feedback", apiFeedback.Save)

		apiAuth.GET("/chat_gpt", new(controllers.CommonController).ChatGPT)
		apiAuth.GET("/ip", new(controllers.CommonController).Ip)
	}

	// =========================
	// 后台 API（必须登录）
	// =========================
	adminAPI := router.Group("/admin_api/v1")
	adminAPI.Use(middleware.AuthMiddleware())
	{
		// 文章管理（完整权限）
		adminAPI.GET("/article", adminArticle.Index)
		adminAPI.GET("/article/detail", adminArticle.Detail)
		adminAPI.POST("/article/save", adminArticle.Save)
		adminAPI.POST("/article/delete", adminArticle.Delete)
		adminAPI.POST("/article/enable", adminArticle.Enable)
		adminAPI.POST("/article/disable", adminArticle.Disable)

		// 广告管理（完整权限）
		adminAPI.GET("/ad", adminAd.Index)
		adminAPI.GET("/ad/detail", adminAd.Detail)
		adminAPI.POST("/ad/save", adminAd.Save)
		adminAPI.POST("/ad/delete", adminAd.Delete)
		adminAPI.POST("/ad/enable", adminAd.Enable)
		adminAPI.POST("/ad/disable", adminAd.Disable)
		
		//视频管理
		adminAPI.GET("/movie", adminMovie.Index)
		adminAPI.GET("/movie/detail", adminMovie.Detail)
		adminAPI.POST("/movie/save", adminMovie.Save)
		adminAPI.POST("/movie/delete", adminMovie.Delete)
		adminAPI.POST("/movie/enable", adminMovie.Enable)
		adminAPI.POST("/movie/disable", adminMovie.Disable)

		// 其他
		adminAPI.POST("/upload", adminUpload.Upload)
		adminAPI.GET("/download", adminUpload.Download)

		adminAPI.GET("/chat_gpt", new(admin.CommonController).ChatGPT)
		adminAPI.GET("/ip", new(admin.CommonController).Ip)
		adminAPI.GET("/ping", new(admin.CommonController).Ping)

		adminAPI.GET("/test", new(admin.TestController).Test)
		adminAPI.GET("/queue", new(admin.TestController).Queue)
		adminAPI.POST("/send_msg", new(admin.CommonController).SendMsg)
	}

	return router
}