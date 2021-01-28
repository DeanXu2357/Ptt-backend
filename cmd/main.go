package main

import (
	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/delivery/http"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"

	"github.com/PichuChen/go-bbs"
	_ "github.com/PichuChen/go-bbs/pttbbs"
)

func main() {
	var logger = logging.NewLogger()

	logger.Informationalf("server start")

	globalConfig, err := config.NewDefaultConfig()
	if err != nil {
		logger.Errorf("failed to get config: %v", err)
		return
	}

	db, err := bbs.Open("pttbbs", globalConfig.BBSHome)
	if err != nil {
		logger.Errorf("open bbs db error: %v", err)
		return
	}

	boardRepo, err := repository.NewBoardRepository(db)
	if err != nil {
		logger.Errorf("failed to create board repository: %s\n", err)
		return
	}

	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		logger.Errorf("failed to create user repository: %s\n", err)
		return
	}

	userUsecase := usecase.NewUserUsecase(userRepo)
	boardUsecase := usecase.NewBoardUsecase(boardRepo)

	httpDelivery := http.NewHTTPDelivery(globalConfig, userUsecase, boardUsecase)
	if err := httpDelivery.Run(); err != nil {
		logger.Errorf("run http delivery error: %s\n", err)
	}
}
