// 跨域解决
package middleware

import (
	"github.com/BurntSushi/toml"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// 多语言初始化
func SetLocale() gin.HandlerFunc {
	return ginI18n.Localize(
		ginI18n.WithBundle(&ginI18n.BundleCfg{
			RootPath:         "./conf/lang",
			AcceptLanguage:   []language.Tag{language.Chinese, language.English},
			DefaultLanguage:  language.Chinese, //默认语言
			FormatBundleFile: "toml",
			UnmarshalFunc:    toml.Unmarshal,
		}),
		ginI18n.WithGetLngHandle(
			func(ctx *gin.Context, defaultLng string) string {
				lng := ctx.Request.Header.Get("language")
				if lng == "" {
					return defaultLng
				}
				return lng
			},
		),
	)
}
