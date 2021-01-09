package main

import (
	"github.com/PichuChen/go-bbs"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	// "strings"
)

// getBoardArticles handles request with `/v1/boards/SYSOP/articles` and will return
// article list to client
func getBoardArticles(w http.ResponseWriter, r *http.Request, boardId string) {
	logger.Debugf("getBoardArticles: %v", r)
	token := getTokenFromRequest(r)
	// Check permission for board
	err := checkTokenPermission(token,
		[]permission{PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	filepath, err := bbs.GetBoardArticleDirectoryath(globalConfig.BBSHome, boardId)
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

func getBoardArticlesFile(w http.ResponseWriter, r *http.Request, boardId string, filename string) {
	logger.Debugf("getBoardArticlesFile: %v", r)

	token := getTokenFromRequest(r)
	err := checkTokenPermission(token,
		[]permission{PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	path, err := bbs.GetBoardArticlePath(globalConfig.BBSHome, boardId, filename)
	file, err := os.Open(path)
	if err != nil {
		logger.Errorf("open file %v error: %v", path, err)
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Errorf("read file %v error: %v", path, err)
	}

	bufStr := base64.StdEncoding.EncodeToString(buf)

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"raw": bufStr,
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)
}

func getArticleURL(boardId string, filename string) string {
	return fmt.Sprintf("https://ptt-app-dev-codingman.pichuchen.tw/bbs/%s/%s.html", boardId, filename)
}
