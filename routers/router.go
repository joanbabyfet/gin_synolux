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
		{
			v1.GET("/home_article", new(controllers.ArticleController).HomeArticle) //首页文章
			v1.GET("/article", new(controllers.ArticleController).Index)
			v1.GET("/article/detail", new(controllers.ArticleController).Detail)
			v1.POST("/article/save", new(controllers.ArticleController).Save)
			v1.POST("/article/delete", new(controllers.ArticleController).Delete)
			v1.POST("/article/enable", new(controllers.ArticleController).Enable)
			v1.POST("/article/disable", new(controllers.ArticleController).Disable)
			v1.POST("/upload", new(controllers.UploadController).Upload)
			v1.GET("/download", new(controllers.UploadController).Download)
			v1.GET("/chat_gpt", new(controllers.CommonController).ChatGPT)
			v1.GET("/ip", new(controllers.CommonController).Ip)
			v1.GET("/ping", new(controllers.CommonController).Ping)
			v1.GET("/captcha", new(controllers.CommonController).Captcha) //获取验证码
			v1.GET("/test", new(controllers.TestController).Test)
			v1.POST("/login", new(controllers.UserController).Login) //登录
			v1.POST("/logout", new(controllers.UserController).Logout)
			v1.POST("/set_password", new(controllers.UserController).SetPassword)
			v1.GET("/get_userinfo", new(controllers.UserController).GetUserInfo) //获取用户信息
			v1.POST("/register", new(controllers.UserController).Register)
			v1.POST("/profile", new(controllers.UserController).Profile)
			v1.POST("/feedback", new(controllers.FeedbackController).Save)
			v1.GET("/weather", new(controllers.CommonController).Weather)   //获取天气信息
			v1.GET("/hardware", new(controllers.CommonController).Hardware) //获取系统信息
			v1.GET("/gist", new(controllers.CommonController).Gist)         //获取Gist信息
		}
	}
	admin_api := router.Group("/admin_api")
	{
		v1 := admin_api.Group("/v1")
		{
			v1.GET("/article", new(admin.ArticleController).Index)
			v1.GET("/article/detail", new(admin.ArticleController).Detail)
			v1.POST("/article/save", new(admin.ArticleController).Save)
			v1.POST("/article/delete", new(admin.ArticleController).Delete)
			v1.POST("/article/enable", new(admin.ArticleController).Enable)
			v1.POST("/article/disable", new(admin.ArticleController).Disable)
			v1.GET("/ad", new(admin.AdController).Index)
			v1.GET("/ad/detail", new(admin.AdController).Detail)
			v1.POST("/ad/save", new(admin.AdController).Save)
			v1.POST("/ad/delete", new(admin.AdController).Delete)
			v1.POST("/ad/enable", new(admin.AdController).Enable)
			v1.POST("/ad/disable", new(admin.AdController).Disable)
			v1.POST("/upload", new(admin.UploadController).Upload)
			v1.GET("/download", new(admin.UploadController).Download)
			v1.GET("/chat_gpt", new(admin.CommonController).ChatGPT)
			v1.GET("/ip", new(admin.CommonController).Ip)
			v1.GET("/ping", new(admin.CommonController).Ping)
			v1.GET("/captcha", new(admin.CommonController).Captcha) //获取验证码
			v1.GET("/test", new(admin.TestController).Test)
			v1.POST("/send_msg", new(admin.CommonController).SendMsg)
		}
	}
	return router
}
