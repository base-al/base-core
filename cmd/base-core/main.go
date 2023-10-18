package main

import (
	basedb "base-core/db"
	"fmt"
	"net/http"
	"os"

	"base-core/core/authentication/auth"
	"base-core/core/authorization/roles"
	"base-core/core/orgs"
	"base-core/core/subscriptions"
	"base-core/core/users"
	swdocs "base-core/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/mattevans/postmark-go"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found")
	}
}

func Index(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Base Core API v1.0", "status": "OK"})
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(http.StatusNotFound).JSON(fiber.Map{
		"message": "Not found",
	})
}

func main() {
	jwtSecret := "123456"
	postmarkAPIKey := ""
	// s3ApiKey := ""
	// s3ApiSecret := ""
	// s3ApiEndpoint := ""
	// s3Region := ""
	// appUrl := ""
	env := ""
	// s3OwnerAccId := ""
	apiPort := 3000

	if os.Getenv("JWT_SECRET") != "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	if os.Getenv("BASE_CORE_ENV") != "" {
		env = os.Getenv("BASE_CORE_ENV")
	}
	// if os.Getenv("APP_URL") != "" {
	// 	uiAppUrl = os.Getenv("APP_URL")
	// }

	// Set swagger info
	swdocs.SwaggerInfo.Title = "Base Core API"
	swdocs.SwaggerInfo.Description = "This is Base API documentation in detail."
	swdocs.SwaggerInfo.Version = "1.0"
	swdocs.SwaggerInfo.BasePath = "/api"
	swdocs.SwaggerInfo.Schemes = []string{"http"}
	swdocs.SwaggerInfo.Host = "localhost:3000"
	if env == "stage" || env == "test" || env == "prod" {
		swdocs.SwaggerInfo.Schemes = []string{"https", "http"}
	}
	if env == "test" {
		swdocs.SwaggerInfo.Host = "test.base.al"
	}
	if env == "stage" {
		apiPort = 4000
		swdocs.SwaggerInfo.Host = "stage.base.al"
	}

	// Setup Postmark client
	postmarkclient := postmark.NewClient(&http.Client{
		Transport: &postmark.AuthTransport{Token: postmarkAPIKey},
	})

	// New fiber app
	app := fiber.New()
	app.Use(cors.New())

	defaultLogger := log.DefaultLogger()

	// Setup DB
	db, err := basedb.ConnectDB()
	if err != nil {
		fmt.Println("basedb.ConnectDB", err)
	}

	// Setup APIs
	userApi := users.NewHTTPTransport(users.NewUserAPI(db, jwtSecret, postmarkclient, defaultLogger), defaultLogger)
	authApi := auth.NewHTTPTransport(auth.NewAuthAPI(db, jwtSecret, postmarkclient, defaultLogger), defaultLogger)
	orgsApi := orgs.NewHTTPTransport(orgs.NewOrgAPI(db, defaultLogger), defaultLogger)

	// Migrate the schema
	db.AutoMigrate(
		&users.User{},
		&orgs.Org{},
		&orgs.OrgSettings{},
		&roles.Role{},
		&orgs.UserOrgRole{},
		&roles.Permission{},
		&subscriptions.Subscription{},
		&subscriptions.Feature{})

	// Routes
	app.Get("/", Index)
	apisRoute := app.Group("/api")
	apisRoute.Get("/swagger/*", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"dev": "basecode",
		},
	}), swagger.HandlerDefault)

	// User API Routes
	users.RegisterRoutes(apisRoute, userApi)
	auth.RegisterRoutes(apisRoute, authApi)
	orgs.RegisterRoutes(apisRoute, orgsApi)

	// Handle no route match
	app.Use(NotFound)

	exit := make(chan error)

	// Serve API
	log.Fatal(app.Listen(fmt.Sprintf(`:%d`, apiPort)))
	exit <- nil
}
