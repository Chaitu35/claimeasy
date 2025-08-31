package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chaitu35/claimeasy/internal/audit"
	"github.com/chaitu35/claimeasy/internal/config"
	"github.com/chaitu35/claimeasy/internal/database"
	"github.com/chaitu35/claimeasy/internal/permission"
	"github.com/chaitu35/claimeasy/internal/role"
	"github.com/chaitu35/claimeasy/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	seed := flag.Bool("seed", false, "Run database seeders")
	flag.Parse()



	// 2️⃣ Connect to DB
	database.Connection()
	db := database.DB

	// 2.1️⃣ Run migrations
	database.Migrate(db)

		if *seed {
		database.Seed(db)
		log.Println("🌱 Database seeding completed")
	}

	// 3️⃣ Init audit storage & logger
	auditStorage := audit.NewDBStorage(db)
	logger := audit.NewLogger(auditStorage)

	

	// 5️⃣ Build Permission module, Role module, User module
	permRepo := permission.NewRepository()
	permService := permission.NewService(permRepo, logger)
	permHandler := permission.NewHandler(permService, db)

	roleRepo := role.NewRepository()
	roleService := role.NewService(roleRepo,  logger)
	roleHandler := role.NewHandler(roleService, db)

	userRepo := user.NewUserRepository()
	userService := user.NewService(userRepo, db, logger,permRepo)
	userHandler := user.NewUserHandler(userService, db)
	

	// 6️⃣ Setup router
	r := gin.Default()

	// 7️⃣ Register routes
	api := r.Group("/api/v1")
	user.RegisterUserRoutes(api, userHandler)
	permission.RegisterPermissionRoutes(api, permHandler)
	role.RegisterRoleRoutes(api, roleHandler)

	// 8️⃣ Register Audit routes
	auditHandler := audit.NewHandler(logger)
	audit.RegisterRoutes(api, auditHandler)

	// 9️⃣ Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("Server starting on port %s...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for Ctrl+C or termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Context with timeout for cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
