package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/* тест для GET запроса без параметров */
func TestShowArticles(t *testing.T) {
	request, err := http.NewRequest("GET", "/articles", nil) // запрос к тестироваему API
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()        // рекордер
	handler := http.HandlerFunc(ShowArticles) // функция которую будем тестировать
	handler.ServeHTTP(recorder, request)      // записываем в рекордер ответ API

	/* check status code */
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	/* check answer (только если данные не динамичны, в ответе присутствует перенос строки) */
	expectedAnswer := `[{"Id":"1","Title":"First title","Author":"First author","Content":"First content"},{"Id":"2","Title":"Second title","Author":"Second author","Content":"Second content"}]
`
	if recorder.Body.String() != expectedAnswer {
		t.Error("Request body does not match expectedAnswers")
	}
}

/* тест для GET запроса c параметрами */
func TestShowArticleById(t *testing.T) {
	request, err := http.NewRequest("GET", "/article", nil)
	if err != nil {
		t.Fatal(err)
	}
	query := request.URL.Query()          // создаём параметры запроса
	query.Add("id", "1")                  // записываем значение параметра
	request.URL.RawQuery = query.Encode() // добавляем парметры к запросу

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(ShowArticleById)
	handler.ServeHTTP(recorder, request)

	/* check status code */
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
