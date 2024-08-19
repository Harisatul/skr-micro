package tryout

type detailTryoutRequest struct {
	ID string `json:"id,omitempty"`
}

type TryoutDetailResponse struct {
	ID          string              `json:"id"`
	CreatedAt   *int64              `json:"created_at"`
	UpdatedAt   *int64              `json:"updated_at"`
	TryOutName  string              `json:"tryout_name"`
	TestType    string              `json:"test_type"`
	TryOutQuota *int                `json:"tryout_quota"`
	TryOutPrice *int                `json:"tryout_price,omitempty"`
	DueDate     *int64              `json:"due_date,omitempty"`
	Question    []QuestionsResponse `json:"question"`
}

type QuestionsResponse struct {
	ID      string           `json:"id"`
	Content string           `json:"content"`
	Weight  *int             `json:"weight,omitempty"`
	Choice  []ChoiceResponse `json:"choice"`
}

type ChoiceResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Weight  *int   `json:"weight,omitempty"`
}
type createOrderTicketPayload struct {
	ID         uint32 `json:"id"`
	CategoryID uint8  `json:"category_id"`
	Version    uint32 `json:"version"`
}

type completeOrderMessage struct {
	ID uint64 `json:"id"`
}

type notifyCancelOrderMessage struct {
	ID uint64 `json:"id"`
}

type cancelOrderMessage struct {
	ID     uint64                   `json:"id"`
	Ticket cancelOrderTicketPayload `json:"ticket"`
}

type cancelOrderTicketPayload struct {
	Version uint32 `json:"version"`
}
