package api

import (
	"avilego.me/recent_news/news"
	"encoding/json"
	"net/http"
)

type searchResponse struct {
	Count int
	Data  searchData
}

type searchData struct {
	Sources  []news.Source
	Previews []previewData
}

type previewData struct {
	Title       string
	Link        string
	Description string
	SourceLink  string
}

func newPreviewData(preview news.Preview) previewData {
	return previewData{
		Title:       preview.Title,
		Link:        preview.Link,
		Description: preview.Description,
		SourceLink:  preview.Source.Link,
	}
}

func newSearchResponse(previews []news.Preview) searchResponse {
	sourcesMap := make(map[string]news.Source)
	sources := make([]news.Source, 0)
	prvsData := make([]previewData, len(previews))

	for i, preview := range previews {
		sourcesMap[preview.Source.Link] = *preview.Source
		prvsData[i] = newPreviewData(preview)
	}
	for _, v := range sourcesMap {
		sources = append(sources, v)
	}

	return searchResponse{
		Count: len(prvsData),
		Data: searchData{
			Sources:  sources,
			Previews: prvsData,
		},
	}
}

type SearchHandler struct {
	Finder news.Finder
}

func (h SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	expr := r.URL.Query().Get("keywords")
	previews := h.Finder.FindRelated(expr)
	searchResponse := newSearchResponse(previews)
	jsonResponse, err := json.Marshal(searchResponse)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		panic(err)
	}
}
