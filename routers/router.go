package routers

import (
	admin "gin-synolux/controllers/admin"
	controllers "gin-synolux/controllers/api"

	"gin-synolux/middleware"

	"github.com/gin-gonic/gin"
)

// 初始化
func Init() *gin.Engine {
	router := gin.New()

	//设置静态资源路径
	router.Static("/uploads", "./uploads")
	//加载视图
	router.LoadHTMLGlob("views/*")

	//跨域解决, 使用路由前进行设置，否则会导致不生效
	router.Use(middleware.Cors())
	//设置多语言文件
	router.Use(middleware.SetLocale())

	//路由分组
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")

		//不需要登录
		public := v1.Group("")
		{
			public.GET("/captcha", new(controllers.CommonController).Captcha)
			public.POST("/login", new(controllers.UserController).Login)
			public.POST("/register", new(controllers.UserController).Register)
			public.GET("/article", new(controllers.ArticleController).Index)
			public.GET("/article/detail", new(controllers.ArticleController).Detail)
			public.GET("/ping", new(controllers.CommonController).Ping)
			public.GET("/weather", new(controllers.CommonController).Weather)
		}

		//需要登录
		auth := v1.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.POST("/article/save", new(controllers.ArticleController).Save)
			auth.POST("/article/delete", new(controllers.ArticleController).Delete)
			auth.POST("/article/enable", new(controllers.ArticleController).Enable)
			auth.POST("/article/disable", new(controllers.ArticleController).Disable)

			auth.POST("/upload", new(controllers.UploadController).Upload)
			auth.GET("/download", new(controllers.UploadController).Download)

			auth.GET("/chat_gpt", new(controllers.CommonController).ChatGPT)
			auth.GET("/ip", new(controllers.CommonController).Ip)
			
			auth.POST("/logout", new(controllers.UserController).Logout)
			auth.POST("/set_password", new(controllers.UserController).SetPassword)
			auth.GET("/get_userinfo", new(controllers.UserController).GetUserInfo)
			auth.POST("/profile", new(controllers.UserController).Profile)
			auth.POST("/feedback", new(controllers.FeedbackController).Save)

			auth.GET("/hardware", new(controllers.CommonController).Hardware)
			auth.GET("/gist", new(controllers.CommonController).Gist)
			auth.GET("/movie", new(controllers.MovieController).Index)
		}
	}
	admin_api := router.Group("/admin_api")
	{
		v1 := admin_api.Group("/v1")

		// 不需要登录
		public := v1.Group("")
		{
			public.GET("/captcha", new(admin.CommonController).Captcha)
		}

		// 需要 admin 权限
		auth := v1.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/article", new(admin.ArticleController).Index)
			auth.GET("/article/detail", new(admin.ArticleController).Detail)
			auth.POST("/article/save", new(admin.ArticleController).Save)
			auth.POST("/article/delete", new(admin.ArticleController).Delete)
			auth.POST("/article/enable", new(admin.ArticleController).Enable)
			auth.POST("/article/disable", new(admin.ArticleController).Disable)

			auth.GET("/ad", new(admin.AdController).Index)
			auth.GET("/ad/detail", new(admin.AdController).Detail)
			auth.POST("/ad/save", new(admin.AdController).Save)
			auth.POST("/ad/delete", new(admin.AdController).Delete)
			auth.POST("/ad/enable", new(admin.AdController).Enable)
			auth.POST("/ad/disable", new(admin.AdController).Disable)

			auth.POST("/upload", new(admin.UploadController).Upload)
			auth.GET("/download", new(admin.UploadController).Download)

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
