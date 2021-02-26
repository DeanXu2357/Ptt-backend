package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetPopularBoardList(t *testing.T) {

	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)
	req, err := http.NewRequest("GET", "/v1/popular-boards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/popular-boards", delivery.routePopularBoards)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responseMap := map[string]interface{}{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
	if err != nil {
		t.Errorf("get unexpect json: %w", err)
	}
	t.Logf("got response %v", rr.Body.String())
	responseData := responseMap["data"]
	popularBoards := responseData.(map[string]interface{})["items"].([]interface{})

	var prevNum int
	for i := range popularBoards {
		curr := popularBoards[i].(map[string]interface{})["number_of_user"]
		currNum, err := strconv.Atoi(curr.(string))
		if err != nil {
			t.Fatalf("handler returned unexpected body, invalid number_of_user: got %v",
				currNum)
		}

		if i > 0 && prevNum < currNum {
			t.Fatalf("handler returned unexpected body, invalid order: got %v before %v",
				prevNum, currNum)
		}
		prevNum = currNum
	}
}

func TestGetBoardList(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	req, err := http.NewRequest("GET", "/v1/boards/", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	w := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards", delivery.routeBoards)
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetBoardInformation(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	req, err := http.NewRequest("GET", "/v1/boards/class/information", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	w := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/class/information", delivery.routeBoards)
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
