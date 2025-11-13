package main

import (
	"log"
	"net/http"
	"os"
    "io/ioutil"
	_ "github.com/Bilal-Cplusoft/sunready/docs"
	"github.com/Bilal-Cplusoft/sunready/internal/client"
	"github.com/Bilal-Cplusoft/sunready/internal/database"
	"github.com/Bilal-Cplusoft/sunready/internal/handler"
	"github.com/Bilal-Cplusoft/sunready/internal/middleware"
	"github.com/Bilal-Cplusoft/sunready/internal/repo"
	"github.com/Bilal-Cplusoft/sunready/internal/service"
	"github.com/go-chi/chi/v5"
	ChiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Sun Ready API
// @version 1.0
// @description API for Sun Ready project
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@sunready.com

// @license.name Sun Ready Private License
// @license.url INTERNAL

// @host localhost:8080

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	databaseURL, jwtSecret, port := os.Getenv("DATABASE_URL"), os.Getenv("JWT_SECRET"), os.Getenv("PORT")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	if port == "" {
		port = "8080"
	}
	db, err := database.New(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// seed hardware options for now
	sqlBytes, err := ioutil.ReadFile("db/init.sql")
	if err != nil {
		log.Fatalf("Failed to read init.sql: %v", err)
	}
	sqlContent := string(sqlBytes)
	if err := db.Exec(sqlContent).Error; err != nil {
		log.Fatalf("Failed to execute init.sql: %v", err)
	}

	userRepo := repo.NewUserRepo(db)
	quoteRepo := repo.NewQuoteRepo(db)
	leadRepo := repo.NewLeadRepo(db)
	houseRepo := repo.NewHouseRepo(db)
	hardwareRepo := repo.NewHardwareRepo(db)

	twilioClient, sendGridClient := client.InitializeTwilio(), client.InitializeSendGrid()
	lightFusionURL, lightFusionEmail, lightFusionPassword := os.Getenv("LIGHTFUSION_API"), os.Getenv("LIGHTFUSION_EMAIL"), os.Getenv("LIGHTFUSION_PASSWORD")
	lightFusionClient := client.NewLightFusionClient(lightFusionURL,lightFusionEmail, lightFusionPassword)

	authService := service.NewAuthService(userRepo, jwtSecret)
	userService := service.NewUserService(userRepo)
	quoteService := service.NewQuoteService(quoteRepo)
	leadService := service.NewLeadService(leadRepo, houseRepo,lightFusionClient,userRepo)


	authHandler := handler.NewAuthHandler(authService, sendGridClient)
	userHandler := handler.NewUserHandler(userService)
	quoteHandler := handler.NewQuoteHandler(quoteService)
	leadHandler := handler.NewLeadHandler(leadRepo, leadService, userRepo)
	otpHandler := handler.NewOtpHandler(twilioClient, sendGridClient)
	hardwareHandler := handler.NewHardwareHandler(hardwareRepo)

	r := chi.NewRouter()

	r.Use(ChiMiddleware.Logger)
	r.Use(ChiMiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Group(func(user chi.Router) {
		user.Use(middleware.AuthMiddleware(authService))
		user.Post("/api/leads", leadHandler.CreateLead)
		user.Get("/api/leads/{id}/mesh-files", leadHandler.GetMeshFiles)
		user.Get("/api/leads/{id}", leadHandler.GetLead)
		user.Put("/api/leads/{id}", leadHandler.UpdateLead)
		user.Post("/api/quote", quoteHandler.GetQuote)
	})
	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AdminMiddleware(authService))
		admin.Get("/admin/users/{id}", userHandler.GetByID)
		admin.Put("/admin/users/{id}", userHandler.Update)
		admin.Delete("/admin/users/{id}", userHandler.Delete)
		admin.Get("/admin/users", userHandler.List)
		admin.Post("/admin/hardware/panel",hardwareHandler.AddPanel)
		admin.Post("/admin/hardware/storage",hardwareHandler.AddStorage)
		admin.Post("/admin/hardware/inverter",hardwareHandler.AddInverter)
		admin.Get("/admin/leads", leadHandler.ListLeads)
		admin.Delete("/admin/leads/{id}", leadHandler.DeleteLead)
	})

	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	r.Get("/api/otp/send", otpHandler.SendOTP)
	r.Get("/api/otp/verify", otpHandler.VerifyOTP)

	r.Get("/api/hardware/panels", hardwareHandler.ListPanels)
	r.Get("/api/hardware/storages", hardwareHandler.ListStorages)
	r.Get("/api/hardware/inverters", hardwareHandler.ListInverters)

	fileServer := http.StripPrefix("/media/", http.FileServer(http.Dir("./media")))
	r.Handle("/media/*", fileServer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "ok",
			"project_name": "SunReady",
			"version": "v1.0.0"
		}`))
	})

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
