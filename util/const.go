package util

import "os"

var (
	STRIPE_SECRET_API_KEY = os.Getenv("STRIPE_SECRET_API_KEY")
	MARGIN                = 0.05
	STRAPI_JWT            = os.Getenv("STRAPI_JWT")
	STRIPE_API_KEY        = os.Getenv("STRIPE_API_KEY")
	STRAPI_URL            = os.Getenv("STRAPI_URL")
	AUTH_EMAIL            = os.Getenv("AUTH_EMAIL")
	MAIL_FORM             = os.Getenv("MAIL_FORM")
	MAIL_AUTH_PASS        = os.Getenv("MAIL_AUTH_PASS")
	SMTP_SERVER_ADDR      = os.Getenv("SMTP_SERVER_ADDR")
	SMTP_SERVER           = os.Getenv("SMTP_SERVER")
	SHIPPING_FEE          = 300
	ADMIN_EMAIL           = os.Getenv("ADMIN_EMAIL")
)
