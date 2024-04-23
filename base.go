package main

import (
	dbhelper "base/db"
	"fmt"
	"net/http"
	"os"

	"base/middleware"
	"base/s3"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/mattevans/postmark-go"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	swdocs "base/swagger"

	// Core modules
	authService "base/core/auth"
	oauthGoogleService "base/core/oauth/google"
	userService "base/core/users"
	// App modules
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found")
	}
}

func Index(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Base Platform API v1", "status": "OK"})
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(http.StatusNotFound).JSON(fiber.Map{
		"message": "Not found",
	})
}

func main() {
	jwtSecret := "123456"
	postmarkAPIKey := ""
	s3ApiKey := ""
	s3ApiSecret := ""
	s3ApiEndpoint := ""
	s3Region := ""
	uiAppUrl := ""
	env := ""
	s3OwnerAccId := ""
	apiPort := 3000
	googleOauthClientId := ""
	googleOauthClientSecret := ""
	baseHostUrl := ""
	sendFromEmail := ""

	if os.Getenv("JWT_SECRET") != "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	if os.Getenv("POSTMARK_API_TOKEN") != "" {
		postmarkAPIKey = os.Getenv("POSTMARK_API_TOKEN")
	}
	if os.Getenv("S3_API_KEY") != "" {
		s3ApiKey = os.Getenv("S3_API_KEY")
	}
	if os.Getenv("S3_API_SECRET") != "" {
		s3ApiSecret = os.Getenv("S3_API_SECRET")
	}
	if os.Getenv("S3_API_ENDPOINT") != "" {
		s3ApiEndpoint = os.Getenv("S3_API_ENDPOINT")
	}
	if os.Getenv("S3_REGION") != "" {
		s3Region = os.Getenv("S3_REGION")
	}
	if os.Getenv("UI_APP_URL") != "" {
		uiAppUrl = os.Getenv("UI_APP_URL")
	}
	if os.Getenv("CAMPER_ENV") != "" {
		env = os.Getenv("CAMPER_ENV")
	}
	if os.Getenv("S3_OWNER_ACCOUNT_ID") != "" {
		s3OwnerAccId = os.Getenv("S3_OWNER_ACCOUNT_ID")
	}
	if os.Getenv("GOOGLE_OAUTH_CLIENT_ID") != "" {
		googleOauthClientId = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	}
	if os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET") != "" {
		googleOauthClientSecret = os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	}
	if os.Getenv("BASE_HOST_URL") != "" {
		baseHostUrl = os.Getenv("BASE_HOST_URL")
	}
	if os.Getenv("POSTMARK_EMAIL_FROM") != "" {
		sendFromEmail = os.Getenv("POSTMARK_EMAIL_FROM")
	}

	// Set swagger info
	swdocs.SwaggerInfo.Title = "Base API"
	swdocs.SwaggerInfo.Description = "This is Base API documentation in detail."
	swdocs.SwaggerInfo.Version = "1.0"
	swdocs.SwaggerInfo.BasePath = "/api"
	swdocs.SwaggerInfo.Schemes = []string{"http"}
	swdocs.SwaggerInfo.Host = "localhost:3000"
	if env == "stage" || env == "test" || env == "prod" {
		swdocs.SwaggerInfo.Schemes = []string{"https", "http"}
	}
	if env == "stage" {
		apiPort = 4000
		swdocs.SwaggerInfo.Host = "base-api.base.al"
	}
	if env == "test" {
		swdocs.SwaggerInfo.Host = "base-api.base.al"
	}

	// Boostrap new fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20 MB in bytes
	})
	app.Use(cors.New())
	app.Use(logger.New())

	// Setup Postmark client
	postmarkclient := postmark.NewClient(&http.Client{
		Transport: &postmark.AuthTransport{Token: postmarkAPIKey},
	})

	// Setup Logger
	defaultLogger := log.DefaultLogger()

	// Setup S3 client instance for AWS
	s3c := s3.New(&s3.S3Config{
		APIKey:      s3ApiKey,
		APISecret:   s3ApiSecret,
		APIEndpoint: s3ApiEndpoint,
		Region:      s3Region,
		Env:         env,
		OwnerAccID:  s3OwnerAccId,
	})

	googleOauth2Cfg := &oauth2.Config{
		ClientID:     googleOauthClientId,
		ClientSecret: googleOauthClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	// Setup DB
	db, err := dbhelper.ConnectDB()
	if err != nil {
		fmt.Println("data.ConnectDB", err)
	}

	// Migrate the schema
	// Core
	db.AutoMigrate(
		&userService.User{},
	)

	// Middlewares
	authMiddleware := middleware.Authentication(jwtSecret)

	// Routes
	app.Get("/", Index)
	apisRoute := app.Group("/api")
	apisRoute.Get("/swagger/*", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"dev": "basecode",
		},
	}), swagger.HandlerDefault)

	// Initialize API services
	authApiSvc := authService.NewAuthHTTPTransport(
		authService.NewAuthAPI(db, jwtSecret, postmarkclient, s3c, uiAppUrl, defaultLogger, googleOauth2Cfg, baseHostUrl),
		defaultLogger,
	)

	userAPISvc := userService.NewUserHTTPTransport(
		userService.NewUserAPI(db, s3c, jwtSecret, postmarkclient, uiAppUrl, defaultLogger, sendFromEmail),
		defaultLogger)

	oAuthAccountAPISvc := oauthGoogleService.NewOAuthAccountHTTPTransport(
		oauthGoogleService.NewOAuthAccountAPI(db, defaultLogger, googleOauth2Cfg, s3c, jwtSecret, baseHostUrl),
		defaultLogger, uiAppUrl,
	)

	userRoute := apisRoute.Group("/users/:userId", authMiddleware)
	// Register API routes
	authService.RegisterRoutes(apisRoute, authApiSvc)
	userService.RegisterRoutes(userRoute, userAPISvc, authMiddleware)
	oauthGoogleService.RegisterRoutes(apisRoute, oAuthAccountAPISvc, authMiddleware)

	// Handle no route match
	app.Use(NotFound)

	exit := make(chan error)

	// Serve API
	log.Fatal(app.Listen(fmt.Sprintf(`:%d`, apiPort)))
	exit <- nil
}
