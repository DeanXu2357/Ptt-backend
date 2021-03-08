package repository

import (
	"context"
	"fmt"

	"github.com/PichuChen/go-bbs"
)

// TODO: go-bbs lacks following interfaces, remove when go-bbs will have implemented functions
type PostsLimitedBoardRecord interface {
	PostLimitPosts() uint8
}

type LoginsLimitedBoardRecord interface {
	PostLimitLogins() uint8
}

type BadPostLimitedBoardRecord interface {
	PostLimitBadPost() uint8
}

type postsLimitedBoardRecord struct{}
type loginsLimitedBoardRecord struct{}
type badPostLimitedBoardRecord struct{}

func (r *postsLimitedBoardRecord) PostLimitPosts() uint8 {
	return 0
}

func (r *loginsLimitedBoardRecord) PostLimitLogins() uint8 {
	return 0
}

func (r *badPostLimitedBoardRecord) PostLimitBadPost() uint8 {
	return 0
}

func (repo *repository) GetBoards(_ context.Context) []bbs.BoardRecord {
	return repo.boardRecords
}

func (repo *repository) GetBoardArticle(_ context.Context, boardID, filename string) ([]byte, error) {
	return repo.db.ReadBoardArticleFile(boardID, filename)
}

func (repo *repository) GetBoardArticleRecords(_ context.Context, boardID string) ([]bbs.ArticleRecord, error) {
	return repo.db.ReadBoardArticleRecordsFile(boardID)
}

func (repo *repository) GetBoardTreasureRecords(_ context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error) {
	return repo.db.ReadBoardTreasureRecordsFile(boardID, treasureIDs)
}

func (repo *repository) GetBoardPostsLimited(_ context.Context, boardID string) (PostsLimitedBoardRecord, error) {
	// TODO: replace postsLimitedBoardRecord to real bbs record
	return &postsLimitedBoardRecord{}, nil
}

func (repo *repository) GetBoardLoginsLimited(_ context.Context, boardID string) (LoginsLimitedBoardRecord, error) {
	// TODO: replace loginsLimitedBoardRecord to real bbs record
	return &loginsLimitedBoardRecord{}, nil
}

func (repo *repository) GetBoardBadPostLimited(_ context.Context, boardID string) (BadPostLimitedBoardRecord, error) {
	// TODO: replace badPostLimitedBoardRecord to real bbs record
	return &badPostLimitedBoardRecord{}, nil
}

func loadBoardFile(db *bbs.DB) ([]bbs.BoardRecord, error) {
	boardRecords, err := db.ReadBoardRecords()
	if err != nil {
		logger.Errorf("get board header error: %v", err)
		return nil, fmt.Errorf("failed to read board records: %w", err)
	}
	for index, board := range boardRecords {
		logger.Debugf("loaded %d %v", index, board.BoardId())
	}
	return boardRecords, nil
}
