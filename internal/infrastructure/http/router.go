package router

import (
	"github.com/gin-gonic/gin"
	"github.com/portaria-keys/internal/controller"
	"github.com/portaria-keys/internal/infrastructure/http/middleware"
	"github.com/portaria-keys/internal/usecase"

	// swagger embed files
	_ "github.com/portaria-keys/docs"
)

// Setup configura o router Gin e registra as rotas da API.
func Setup(userUseCase *usecase.UserUseCase, keyUseCase *usecase.KeyUseCase, reservationUseCase *usecase.ReservationUseCase) *gin.Engine {
	r := gin.New()

	// Global Middlewares
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	// Controllers
	authController := controller.NewAuthController(userUseCase)
	keyController := controller.NewKeyController(keyUseCase)
	reservationController := controller.NewReservationController(reservationUseCase)
	adminController := controller.NewAdminController(userUseCase)

	// Public Routes
	public := r.Group("/api/v1")
	{
		public.POST("/auth/register", authController.Register)
		public.POST("/auth/login", authController.Login)
	}

	// Authenticated Routes
	auth := r.Group("/api/v1")
	auth.Use(middleware.AuthMiddleware())
	{
		// Key Routes
		auth.GET("/keys", keyController.GetAllKeys)
		auth.GET("/keys/:id", keyController.GetKeyByID)

		// Reservation Routes
		auth.POST("/reservations", reservationController.CreateReservation)
		auth.GET("/reservations", reservationController.GetUserReservations)
		public.PUT("/reservations/:id/return", reservationController.ReturnKey)
		auth.GET("/reservations/history", reservationController.GetReservationHistory)
	}

	// Admin Routes
	admin := r.Group("/api/v1")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// Key Admin Routes
		admin.POST("/keys", keyController.CreateKey)
		admin.PUT("/keys/:id", keyController.UpdateKey)
		admin.DELETE("/keys/:id", keyController.DeleteKey)

		// Admin Reservation Routes
		admin.GET("/admin/reservations", reservationController.GetAllReservations)
		admin.PUT("/admin/reservations/:id/extend", reservationController.ExtendReservation)

		// Admin User Routes
		admin.POST("/admin/users/:id/block", adminController.BlockUser)
		admin.POST("/admin/users/:id/unblock", adminController.UnblockUser)
	}

	// Serve Swagger UI static files
	r.StaticFile("/swagger/doc.json", "./docs/swagger.json")
	r.StaticFile("/swagger/doc.yaml", "./docs/swagger.yaml")
	// You would typically serve an index.html for Swagger UI here, but for simplicity,
	// we're just serving the spec files directly. A full Swagger UI setup would involve
	// downloading the Swagger UI dist files and serving them statically.

	return r
}