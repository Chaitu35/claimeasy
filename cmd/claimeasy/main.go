package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	//"strconv"
	"syscall"
	"time"

	config "github.com/chaitu35/claimeasy/configs"
	"github.com/chaitu35/claimeasy/internal/audit"
	"github.com/chaitu35/claimeasy/internal/auth"
	"github.com/chaitu35/claimeasy/internal/clinics"
	"github.com/chaitu35/claimeasy/internal/database"
	"github.com/chaitu35/claimeasy/internal/permission"
	"github.com/chaitu35/claimeasy/internal/role"
	"github.com/chaitu35/claimeasy/internal/user"
	notification "github.com/chaitu35/claimeasy/pkg/notifications"
	security "github.com/chaitu35/claimeasy/pkg/security"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	seed := flag.Bool("seed", false, "Run database seeders")
	flag.Parse()

	// ----------------------
	// Connect to DB
	// ----------------------
	database.Connection()
	db := database.DB
	database.Migrate(db)

	if *seed {
		database.Seed(db)
		log.Println("🌱 Database seeding completed")
	}

	// ----------------------
	// Init audit
	// ----------------------
	auditStorage := audit.NewDBStorage(db)
	logger := audit.NewLogger(auditStorage)

	// ----------------------
	// JWT Service
	// ----------------------
	accessExpiry := time.Minute * time.Duration(cfg.AccessTokenExpiryMin)
	refreshExpiry := time.Hour * 24 * time.Duration(cfg.RefreshTokenExpiryDays)
	jwtService := security.NewJWTService(cfg.JWTSecret, "ClaimEasy", accessExpiry, refreshExpiry)

	// ----------------------
	// Mailer
	// ----------------------
	mailer := notification.ProvideMailer(cfg)

	// ----------------------
	// Services
	// ----------------------
	clinicRepo := clinics.NewRepository(db)
    clinicService := clinics.NewService(clinicRepo, logger,db)
	permRepo := permission.NewRepository()
	permService := permission.NewService(permRepo, logger)
	roleRepo := role.NewRepository()
	roleService := role.NewService(roleRepo, logger)
	userRepo := user.NewUserRepository()
	userService := user.NewService(userRepo, db, logger, permRepo,clinicService)
	authService := auth.NewService(jwtService, permService,*userService, logger,mailer)

			// ----------------------
	// Handlers
	// ----------------------
	clinicHandler:= clinics.NewHandler(clinicService)
	permHandler := permission.NewHandler(permService, db)
	roleHandler := role.NewHandler(roleService, db)
	userHandler := user.NewUserHandler(userService, db)
	authHandler := auth.NewHandler(authService, userService,db)

	// ----------------------
	// Router
	// ----------------------
	r := gin.Default()
	api := r.Group("/api/v1")
	user.RegisterUserRoutes(api, userHandler)
	permission.RegisterPermissionRoutes(api, permHandler)
	role.RegisterRoleRoutes(api, roleHandler)
	auth.RegisterAuthRoutes(api, authHandler)
	clinics.RegisterClinicRoutes(api, clinicHandler)

	// Audit routes
	auditHandler := audit.NewHandler(logger)
	audit.RegisterRoutes(api, auditHandler)

	// ----------------------
	// Start server
	// ----------------------
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// ----------------------
	// Graceful shutdown
	// ----------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
