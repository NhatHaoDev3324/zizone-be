package word

import (
	"github.com/NhatHaoDev3324/zizone-be/constant"
	"github.com/NhatHaoDev3324/zizone-be/internal/middleware"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/word/handler"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/word/repository"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/word/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func WordRoutes(r *gin.RouterGroup, db *gorm.DB, redis *redis.Client) {
	repo := repository.NewWordRepository(db, redis)
	svc := service.NewWordService(repo)
	h := handler.NewWordHandler(svc)

	wordGroup := r.Group("/word")
	{
		wordGroup.GET("/list", h.GetListWord)
		wordGroup.GET("/detail/:id", h.GetWordByID)

		adminWordGroup := wordGroup.Group("/")
		adminWordGroup.Use(middleware.RequireAuth(), middleware.RequireRole(constant.RoleAdmin))
		{
			adminWordGroup.POST("/create", h.CreateWord)
			adminWordGroup.PUT("/update", h.UpdateWord)
			adminWordGroup.DELETE("/delete/:id", h.DeleteWord)
		}
	}

}
