package main

import (
	"net/http"
	"time"

	"github.com/charisworks/charisworks-service-go/handler"
	"github.com/charisworks/charisworks-service-go/meilisearch"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	CORS(r)
	h := handler.NewHandler(r)
	meilisearch.InitMeilisearch()
	h.SetupStripeEventHandler()
	h.SetupStrapiEventHandler()
	h.SetupCheckoutEventHandler()
	h.SetupItmeHandler()
	http.ListenAndServe(":8080", r)
}
func CORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		// アクセス許可するオリジン
		AllowOrigins: []string{
			"https://beta.charis.works",
			"https://search.charis.works",
			"http://next:3000",
			"https://api.charis.works",
			"https://charis.works",
		},
		// アクセス許可するHTTPメソッド
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"PATCH",
			"DELETE",
		},
		// 許可するHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Content-Type",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Authorization",
			"Access-Control-Allow-Credentials",
		},

		// cookieなどの情報を必要とするかどうか
		AllowCredentials: true,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	}))
}
