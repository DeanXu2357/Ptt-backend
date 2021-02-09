package repository

import (
	"context"
	"testing"

	"github.com/PichuChen/go-bbs"
)

func TestGetUsers(t *testing.T) {

	repo := repository{
		userRecords:  []bbs.UserRecord{},
		boardRecords: []bbs.BoardRecord{},
	}

	actual := repo.GetUsers(context.TODO())
	if actual == nil {
		t.Errorf("GetUsers got %v, expected not equal nil", actual)
	}

}
