package main

import (
	"github.com/PichuChen/go-bbs"

	// "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	// "os"
)

func getBoardTreasures(w http.ResponseWriter, r *http.Request, boardId string) {
	logger.Debugf("getBoardTreasures: %v", r)
	token := getTokenFromRequest(r)
	_, treasuresId, filename, err := parseBoardTreasurePath(r.URL.Path)
	if err != nil {
		logger.Warningf("parseBoardTreasurePath error: %v", err)
		// TODO return 400?
	}
	if filename != "" {
		// get file
		getBoardTreasuresFile(w, r, boardId, treasuresId, filename)
		return
	}

	// Check permission for board
	err = checkTokenPermission(token,
		[]permission{PermissionReadTreasureInformation},
		map[string]string{
			"board_id":    boardId,
			"treasure_id": strings.Join(treasuresId, ","),
		})
	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	filepath, err := bbs.GetBoardTreasuresDirectoryPath(globalConfig.BBSHome, boardId, treasuresId)
	logger.Debugf("open DIR file: %v", filepath)

	var fileHeaders []*bbs.FileHeader
	fileHeaders, err = bbs.OpenFileHeaderFile(filepath)
	if err != nil {
		logger.Warningf("open directory file error: %v", err)
		// The board may not contain any article
	}

	items := []interface{}{}
	for _, f := range fileHeaders {
		m := map[string]interface{}{
			"filename": f.Filename,
			// Bug(pichu): f.Modified time will be 0 when file is vote
			"modified_time":   f.Modified,
			"recommend_count": f.Recommend,
			"post_date":       f.Date,
			"title":           f.Title,
			"money":           fmt.Sprintf("%v", f.Money),
			"owner":           f.Owner,
			// "aid": ""
			"url": getArticleURL(boardId, f.Filename),
		}
		items = append(items, m)
	}
	logger.Debugf("fh: %v", fileHeaders)

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": items,
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)

}

func getBoardTreasuresFile(w http.ResponseWriter, r *http.Request, boardId string, treasuresId []string, filename string) {
	logger.Debugf("getBoardTreasuresFile %v board: %v, treasuresId: %v, filename: %v", r, boardId, treasuresId, filename)

	w.WriteHeader(http.StatusNotImplemented)
}

// parseBoardTreasurePath parse covert url path from /v1/boards/SYSOP/article to
// {SYSOP, article) or /v1/boards to {,}
func parseBoardTreasurePath(path string) (boardId string, treasuresId []string, filename string, err error) {
	pathSegment := strings.Split(path, "/")

	if len(pathSegment) == 6 {
		// /{{version}}/boards/{{board_id}}/treasures/articles
		boardId = pathSegment[3]
		treasuresId = []string{}
		filename = ""
		return
	} else if len(pathSegment) >= 7 {
		// /{{version}}/boards/{{board_id}}/treasures/{{treasures_id ... }}/articles
		// or
		// /{{version}}/boards/{{board_id}}/treasures/{{treasures_id ... }}/articles/{{filename}}
		boardId = pathSegment[3]
		if pathSegment[len(pathSegment)-1] == "articles" {
			treasuresId = pathSegment[5 : len(pathSegment)-1]
			filename = ""
		} else {
			treasuresId = pathSegment[5 : len(pathSegment)-2]
			filename = pathSegment[len(pathSegment)-1]
		}
		return
	}
	// should not be reached
	logger.Warningf("parseBoardTreasurePath got malform path: %v", path)
	return

}
