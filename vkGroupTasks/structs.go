package vkGroupSchedule


type responseItem struct {
	ID       int    `json:"id"`
	Owner    int    `json:"owner_id"`
	Text     string `json:"text"`
	IsPinned int    `json:"is_pinned"`
	Date     int    `json:"date"`
}

type vkResponse struct {
	Response response `json:"response"`
}

type response struct {
	Items []responseItem `json:"items"`
}
