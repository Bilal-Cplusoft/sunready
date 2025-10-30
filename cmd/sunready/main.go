package main

import (
	"log"
	"net/http"
	"os"
    "github.com/Bilal-Cplusoft/sunready/internal/middleware"
	"github.com/Bilal-Cplusoft/sunready/internal/client"
	"github.com/Bilal-Cplusoft/sunready/internal/database"
	"github.com/Bilal-Cplusoft/sunready/internal/handler"
	"github.com/Bilal-Cplusoft/sunready/internal/repo"
	"github.com/Bilal-Cplusoft/sunready/internal/service"
	"github.com/go-chi/chi/v5"
	ChiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/Bilal-Cplusoft/sunready/docs"
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

	userRepo := repo.NewUserRepo(db)
	customerRepo := repo.NewCustomerRepo(db)
	projectRepo := repo.NewProjectRepo(db)
	quoteRepo := repo.NewQuoteRepo(db)
	leadRepo := repo.NewLeadRepo(db)
	houseRepo := repo.NewHouseRepo(db)

	twilioClient,sendGridClient := client.InitializeTwilio(),client.InitializeSendGrid()

	authService := service.NewAuthService(userRepo, jwtSecret)
	userService := service.NewUserService(userRepo)
	customerService := service.NewCustomerService(customerRepo)
	projectService := service.NewProjectService(projectRepo)
	quoteService := service.NewQuoteService(quoteRepo)
	leadService := service.NewLeadService(leadRepo,houseRepo)
	authHandler := handler.NewAuthHandler(authService,sendGridClient)
	userHandler := handler.NewUserHandler(userService)
	customerHandler := handler.NewCustomerHandler(customerService)
	projectHandler := handler.NewProjectHandler(projectService)
	quoteHandler := handler.NewQuoteHandler(quoteService)
	leadHandler := handler.NewLeadHandler(leadRepo,leadService,userRepo)
	otpHandler := handler.NewOtpHandler(twilioClient,sendGridClient)

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
     user.Post("/api/customers", customerHandler.CreateCustomer)
   	 user.Get("/api/customers/stats", customerHandler.GetCustomerStats)
	 user.Get("/api/customers/{id}", customerHandler.GetCustomer)
	 user.Put("/api/customers/{id}", customerHandler.UpdateCustomer)
	 user.Delete("/api/customers/{id}", customerHandler.DeleteCustomer)
	 user.Patch("/api/customers/{id}/status", customerHandler.UpdateCustomerStatus)
	 user.Post("/api/projects", projectHandler.Create)
	 user.Get("/api/projects/{id}", projectHandler.GetByID)
	 user.Put("/api/projects/{id}", projectHandler.Update)
	 user.Delete("/api/projects/{id}", projectHandler.Delete)
    })
	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AdminMiddleware(authService))
		admin.Get("/api/users/{id}", userHandler.GetByID)
		admin.Put("/api/users/{id}", userHandler.Update)
		admin.Delete("/api/users/{id}", userHandler.Delete)
		admin.Get("/api/users", userHandler.List)
		admin.Get("/api/customers", customerHandler.ListCustomers)
	})

	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)
	r.Get("/api/projects/user", projectHandler.ListByUser)


	r.Post("/api/quote", quoteHandler.GetQuote)

	r.Get("/api/leads", leadHandler.ListLeads)
	r.Get("/api/leads/{id}", leadHandler.GetLead)
	r.Put("/api/leads/{id}", leadHandler.UpdateLead)
	r.Delete("/api/leads/{id}", leadHandler.DeleteLead)

	r.Get("/api/otp/send",otpHandler.SendOTP)
	r.Get("/api/otp/verify",otpHandler.VerifyOTP)

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
