package config

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/triliun/mcupload/backend/features/auth"
	"github.com/triliun/mcupload/backend/features/category"
	"github.com/triliun/mcupload/backend/features/minecraft_version"
	"github.com/triliun/mcupload/backend/features/pack_category"
	"github.com/triliun/mcupload/backend/features/resource_pack"
	"github.com/triliun/mcupload/backend/features/user"
	myMiddleware "github.com/triliun/mcupload/backend/middleware"
	authM "github.com/triliun/mcupload/backend/middleware/auth"
)

func SetupRouter(db *sqlx.DB) *chi.Mux {
	authMRepository := authM.NewPostgresRepository(db)
	authMService := authM.NewService(authMRepository)
	authMiddleware := authM.NewAuthMiddleware(authMService, os.Getenv("JWT_SECRET"))

	authRepository := auth.NewPostgresRepository(db)
	authService := auth.NewService(authRepository)
	authHandler := auth.NewHandler(authService, authMiddleware)

	userRepository := user.NewPostgresRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService, authMiddleware)

	minecraftVersionRepository := minecraft_version.NewPostgresRepository(db)
	minecraftVersionService := minecraft_version.NewService(minecraftVersionRepository)
	minecraftVersionHandler := minecraft_version.NewHandler(minecraftVersionService, authMiddleware)

	categoryRepository := category.NewPostgresRepository(db)
	categoryService := category.NewService(categoryRepository)
	categoryHandler := category.NewHandler(categoryService, authMiddleware)

	packRepository := resource_pack.NewPostgresRepository(db)
	packService := resource_pack.NewService(packRepository)
	packHandler := resource_pack.NewHandler(packService, authMiddleware)

	packCategoryRepository := pack_category.NewPostgresRepository(db)
	packCategoryService := pack_category.NewService(packCategoryRepository)
	packCategoryHandler := pack_category.NewHandler(packCategoryService, authMiddleware)

	// Setup new chi router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(myMiddleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)
	router.Use(myMiddleware.Compress())
	router.Use(myMiddleware.Database(db))

	router.Route("/api/v1", func(v1 chi.Router) {
		authHandler.RegisterRoutes(v1)
		userHandler.RegisterRoutes(v1)
		minecraftVersionHandler.RegisterRoutes(v1)
		categoryHandler.RegisterRoutes(v1)
		packHandler.RegisterRoutes(v1)
		packCategoryHandler.RegisterRoutes(v1)
	})

	return router
}
