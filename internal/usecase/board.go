package usecase

import (
	"context"
	"fmt"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

type boardUsecase struct {
	logger    logging.Logger
	boardRepo repository.BoardRepository
}

func NewBoardUsecase(boardRepo repository.BoardRepository) BoardUsecase {
	return &boardUsecase{
		logger:    logging.NewLogger(),
		boardRepo: boardRepo,
	}
}

func (b *boardUsecase) GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error) {
	for _, it := range b.boardRepo.GetBoards(ctx) {
		if boardID == it.BoardId() {
			return it, nil
		}
	}
	return nil, fmt.Errorf("board record not found")
}

func (b *boardUsecase) GetBoards(ctx context.Context, userID string) []bbs.BoardRecord {
	boards := make([]bbs.BoardRecord, 0)
	for _, board := range b.boardRepo.GetBoards(ctx) {
		// TODO: Show Board by user level
		if board.IsClass() {
			continue
		}
		if !b.shouldShowOnUserLevel(board, userID) {
			continue
		}
		boards = append(boards, board)
	}
	return boards
}

func (b *boardUsecase) GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord {
	boards := make([]bbs.BoardRecord, 0)
	for _, board := range b.boardRepo.GetBoards(ctx) {
		// TODO: Show Board by user level
		if !b.shouldShowOnUserLevel(board, userID) {
			continue
		}
		if board.ClassId() != classID {
			continue
		}
		// m := marshalBoardHeader(board)
		// if board.IsClass() {
		// 	m["id"] = fmt.Sprintf("%v", bid+1)
		// }
		boards = append(boards, board)
	}
	return boards
}

func (b *boardUsecase) GetBoardArticles(ctx context.Context, boardID string) []interface{} {
	articleRecords, err := b.boardRepo.GetBoardArticleRecords(ctx, boardID)
	if err != nil {
		b.logger.Warningf("open directory file error: %v", err)
		// The board may not contain any article
	}

	items := []interface{}{}
	for _, f := range articleRecords {
		m := map[string]interface{}{
			"filename": f.Filename(),
			// Bug(pichu): f.Modified time will be 0 when file is vote
			"modified_time":   f.Modified(),
			"recommend_count": f.Recommend(),
			"post_date":       f.Date(),
			"title":           f.Title(),
			"money":           fmt.Sprintf("%v", f.Money()),
			"owner":           f.Owner(),
			// "aid": ""
			"url": getArticleURL(boardID, f.Filename()),
		}
		items = append(items, m)
	}
	return items
}

func (b *boardUsecase) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	buf, err := b.boardRepo.GetBoardArticle(ctx, boardID, filename)
	if err != nil {
		return nil, fmt.Errorf("read file %s/%s error: %w", boardID, filename, err)
	}
	return buf, nil
}

func (b *boardUsecase) GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} {
	fileHeaders, err := b.boardRepo.GetBoardTreasureRecords(ctx, boardID, treasuresID)
	if err != nil {
		b.logger.Warningf("open directory file error: %v", err)
		// The board may not contain any article
	}

	items := []interface{}{}
	for _, f := range fileHeaders {
		m := map[string]interface{}{
			"filename": f.Filename(),
			// Bug(pichu): f.Modified time will be 0 when file is vote
			"modified_time":   f.Modified(),
			"recommend_count": f.Recommend(),
			"post_date":       f.Date(),
			"title":           f.Title(),
			"money":           fmt.Sprintf("%v", f.Money()),
			"owner":           f.Owner(),
			// "aid": ""
			"url": getArticleURL(boardID, f.Filename()),
		}
		items = append(items, m)
	}
	return items
}

func (b *boardUsecase) shouldShowOnUserLevel(board bbs.BoardRecord, userID string) bool {
	// TODO: Get user Level
	return true
}

func getArticleURL(boardId string, filename string) string {
	return fmt.Sprintf("https://ptt-app-dev-codingman.pichuchen.tw/bbs/%s/%s.html", boardId, filename)
}
