package util

import "os"

var (
	PRODUCTION = false
	// Constants
	MARGIN       = 0.05
	SHIPPING_FEE = 300

	// Stripe
	STRIPE_SECRET_API_KEY = os.Getenv("STRIPE_SECRET_API_KEY")
	STRIPE_API_KEY        = os.Getenv("STRIPE_API_KEY")

	// Strapi
	STRAPI_JWT = os.Getenv("STRAPI_JWT")
	STRAPI_URL = os.Getenv("STRAPI_URL")

	// Mail
	AUTH_EMAIL       = os.Getenv("AUTH_EMAIL")
	MAIL_FORM        = os.Getenv("MAIL_FORM")
	MAIL_AUTH_PASS   = os.Getenv("MAIL_AUTH_PASS")
	SMTP_SERVER_ADDR = os.Getenv("SMTP_SERVER_ADDR")
	SMTP_SERVER      = os.Getenv("SMTP_SERVER")
	ADMIN_EMAIL      = os.Getenv("ADMIN_EMAIL")

	// MeiliSearch
	MEILI_MASTER_KEY            = os.Getenv("MEILI_MASTER_KEY")
	MEILI_URL                   = os.Getenv("MEILI_URL")
	MEILI_ITEM_INDEX            = "items"
	MEILI_ITEM_INDEX_IDENTIFIER = "id"

	// Cloudfare R2
	R2_ENDPOINT          = os.Getenv("ENDPOINT")
	R2_ACCOUNT_ID        = os.Getenv("ACCOUNT_ID")
	R2_ACCESS_KEY_ID     = os.Getenv("ACCESS_KEY_ID")
	R2_ACCESS_KEY_SECRET = os.Getenv("ACCESS_KEY_SECRET")
	R2_BUCKET_NAME       = os.Getenv("BUCKET_NAME")
	IMAGES_URL           = os.Getenv("IMAGES_URL")
)
