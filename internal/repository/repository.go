package repository

import (
	"context"
	"fmt"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
)

var (
	logger = logging.NewLogger()
)

// Repository directly interacts with database via db handler.
type Repository interface {

	// board.go
	// GetBoards return all board record
	GetBoards(ctx context.Context) []bbs.BoardRecord
	// GetBoardArticle returns an article file in a specified board and filename
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	// GetBoardArticleRecords returns article records of a board
	GetBoardArticleRecords(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error)
	// GetBoardTreasureRecords returns treasure article records of a board
	GetBoardTreasureRecords(ctx context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error)
	// GetBoardPostsLimited returns posts limited record of a board
	// TODO: replace PostsLimitedBoardRecord with real bbs record
	GetBoardPostsLimited(ctx context.Context, boardID string) (PostsLimitedBoardRecord, error)
	// GetBoardLoginsLimited returns logins limited record of a board
	// TODO: replace LoginsLimitedBoardRecord with real bbs record
	GetBoardLoginsLimited(ctx context.Context, boardID string) (LoginsLimitedBoardRecord, error)
	// GetBoardBadPostLimited returns bad posts limited record of a board
	// TODO: replace BadPostLimitedBoardRecord with real bbs record
	GetBoardBadPostLimited(ctx context.Context, boardID string) (BadPostLimitedBoardRecord, error)

	// user.go
	// GetUsers returns all user reords
	GetUsers(ctx context.Context) ([]bbs.UserRecord, error)
	// GetUserFavoriteRecords returns favorite records of a user
	GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error)
	// GetUserArticles returns user's articles
	GetUserArticles(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error)

	// article.go
	// GetPopularArticles returns all popular articles
	GetPopularArticles(ctx context.Context) ([]PopularArticleRecord, error)
}

type repository struct {
	db           *bbs.DB
	userRecords  []bbs.UserRecord
	boardRecords []bbs.BoardRecord
}

func NewRepository(db *bbs.DB) (Repository, error) {
	userRecords, err := loadUserRecords(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load user records: %w", err)
	}

	boardRecords, err := loadBoardFile(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load board file: %w", err)
	}

	return &repository{
		db:           db,
		boardRecords: boardRecords,
		userRecords:  userRecords,
	}, nil
}
