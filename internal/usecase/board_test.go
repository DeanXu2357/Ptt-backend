package usecase

import (
	"context"
	"testing"

	"github.com/PichuChen/go-bbs"
)

func TestSearchArticles(t *testing.T) {
	repository := &MockRepository{}
	articleRecords, _ := repository.GetBoardArticleRecords(context.Background(), "")

	type TestCase struct {
		input         *ArticleSearchCond
		expectedItems []bbs.ArticleRecord
	}

	mockArticle1 := MockArticleRecord{
		recommendCount: 10,
		owner:          "SYSOP",
		title:          "[討論] 偶爾要發個廢文",
	}

	mockArticle2 := MockArticleRecord{
		recommendCount: -20,
		owner:          "XDXD",
		title:          "[問題] UNICODE",
	}

	cases := []TestCase{
		{
			input: &ArticleSearchCond{
				Title:                 "廢文",
				Author:                "s",
				RecommendCountGe:      0,
				RecommendCountLe:      50,
				RecommendCountGeIsSet: true,
				RecommendCountLeIsSet: true,
			},
			expectedItems: []bbs.ArticleRecord{&mockArticle1},
		},
		{
			input: &ArticleSearchCond{
				Title:                 "",
				Author:                "X",
				RecommendCountLe:      10,
				RecommendCountLeIsSet: true,
			},
			expectedItems: []bbs.ArticleRecord{&mockArticle2},
		},
	}

	for index, c := range cases {
		cond := c.input
		expectedItems := c.expectedItems

		actualItems := searchArticles(articleRecords, cond)
		for i, v := range actualItems {
			if v.Title() != expectedItems[i].Title() ||
				v.Owner() != expectedItems[i].Owner() ||
				v.Recommend() != expectedItems[i].Recommend() {
				t.Errorf("article not match on index %d, expected: %v, got: %v", index, expectedItems[i], v)
			}
		}
	}
}
