package api

import (
	docs "github.com/coven-discord-bot/internal/controller/api/docs"
	"github.com/coven-discord-bot/internal/usecase"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type API struct {
	Engine *gin.Engine
	*usecase.AbsenceUseCase
	*usecase.PlayerUseCase
	*usecase.StrikeUseCase
	*usecase.LootUseCase
	*usecase.RaidUseCase
}

//	@BasePath	/api/v1

// player
//  /:id/loots
//  /:id/strikes
//  /:id/absences
// /:id/raids
// player/?name=xxx
// player/

func (a *API) Init() {
	//logger, _ := zap.NewProduction()
	//
	//a.Engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	//a.Engine.Use(ginzap.RecoveryWithZap(logger, true))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := a.Engine.Group("/api/v1")
	{
		player := v1.Group("/player")
		{
			player.GET("/", a.GetPlayer)
			player.DELETE("/", a.RemovePlayer)
			player.PUT("/", a.PutPlayer)
			player.PATCH("/", a.PatchPlayer)
			id := player.Group("/id")
			{
				id.GET("/", a.GetPlayerByID)
			}
			loots := player.Group("/loot")
			{
				loots.GET("/", a.GetLootOnPlayer)
			}
			strikes := player.Group("/strike")
			{
				strikes.GET("/", a.GetStrikeOnPlayer)
			}
			absences := player.Group("/absence")
			{
				absences.GET("/", a.GetAbsenceOnPlayer)
			}
		}
		strike := v1.Group("/strike")
		{
			strike.GET("/", a.GetStrike)
			strike.POST("/", a.AddStrike)
			strike.DELETE("/", a.RemoveStrike)
			strike.PUT("/", a.PutStrike)
			strike.PATCH("/", a.PatchStrike)
		}
	}
	a.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}
