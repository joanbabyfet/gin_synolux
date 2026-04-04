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

		// =========================
		// 前台登录用户
		// =========================
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(common.RoleUser))
		{
			//用户
			auth.POST("/logout", apiUser.Logout)
			auth.GET("/get_userinfo", apiUser.GetUserInfo) //获取用户信息
			auth.POST("/profile", apiUser.Profile)
			auth.POST("/set_password", apiUser.SetPassword)
		}
	}

	

	// =========================
	// 后台 API
	// =========================
	adminAPI := router.Group("/admin_api/v1")
	{
		adminAPI.POST("/login", adminAdmin.Login)
		adminAPI.POST("/register", adminAdmin.Register)
		adminAPI.GET("/captcha", new(admin.CommonController).Captcha)

		// =========================
		// 后台 API（必须登录）
		// =========================
		auth := adminAPI.Group("")
		auth.Use(middleware.AuthMiddleware(common.RoleAdmin))
		{
			// 文章管理（完整权限）
			auth.GET("/article", adminArticle.Index)
			auth.GET("/article/detail", adminArticle.Detail)
			auth.POST("/article/save", adminArticle.Save)
			auth.POST("/article/delete", adminArticle.Delete)
			auth.POST("/article/enable", adminArticle.Enable)
			auth.POST("/article/disable", adminArticle.Disable)

			// 广告管理（完整权限）
			auth.GET("/ad", adminAd.Index)
			auth.GET("/ad/detail", adminAd.Detail)
			auth.POST("/ad/save", adminAd.Save)
			auth.POST("/ad/delete", adminAd.Delete)
			auth.POST("/ad/enable", adminAd.Enable)
			auth.POST("/ad/disable", adminAd.Disable)
			
			//视频管理
			auth.GET("/movie", adminMovie.Index)
			auth.GET("/movie/detail", adminMovie.Detail)
			auth.POST("/movie/save", adminMovie.Save)
			auth.POST("/movie/delete", adminMovie.Delete)
			auth.POST("/movie/enable", adminMovie.Enable)
			auth.POST("/movie/disable", adminMovie.Disable)
			
			//用户
			auth.POST("/logout", adminAdmin.Logout)
			auth.GET("/get_userinfo", adminAdmin.GetUserInfo) //获取用户信息
			auth.POST("/profile", adminAdmin.Profile)
			auth.POST("/set_password", adminAdmin.SetPassword)

			// 其他
			auth.POST("/upload", adminUpload.Upload)
			auth.GET("/download", adminUpload.Download)

			auth.GET("/chat_gpt", new(admin.CommonController).ChatGPT)
			auth.GET("/ip", new(admin.CommonController).Ip)
			auth.GET("/ping", new(admin.CommonController).Ping)

			auth.GET("/test", new(admin.TestController).Test)
			auth.GET("/queue", new(admin.TestController).Queue)
			auth.POST("/send_msg", new(admin.CommonController).SendMsg)
		}
	}

	return router
}