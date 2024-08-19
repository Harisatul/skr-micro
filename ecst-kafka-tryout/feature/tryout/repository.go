package tryout

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
)

func getTryoutWithQuestions(ctx context.Context, tryoutID string) (tryoutDetail1 TryoutDetailResponse, err error) {

	// SQL query
	query := `
        SELECT
            t.id AS tryout_id,
            t.created_at AS tryout_created_at,
            t.updated_at AS tryout_updated_at,
            t.tryout_name,
            t.tryout_type,
            t.tryout_quota,
            t.tryout_price,
            t.due_date,
            q.id AS question_id,
            q.content AS question_content,
            q.weight AS question_weight,
            c.id AS choice_id,
            c.content AS choice_content,
            c.weight AS choice_weight
        FROM
            try_outs t
            LEFT JOIN questions q ON t.id = q.tryout_id
            LEFT JOIN choices c ON q.id = c.question_id
        WHERE
            t.id = $1
    `
	rows, err := db.Query(ctx, query, tryoutID)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()
	// Check if any rows are returned
	if !rows.Next() {
		return TryoutDetailResponse{}, errors.New("no tryout detail found")
	}

	tryoutDetail := TryoutDetailResponse{
		Question: []QuestionsResponse{},
	}

	questionMap := make(map[string]*QuestionsResponse)

	for rows.Next() {
		var tryoutID uuid.UUID
		var tryoutName, tryoutType string
		var tryoutCreatedAt, tryoutUpdatedAt, dueDate *int64
		var tryoutQuota, tryoutPrice *int
		var questionID, questionContent *string
		var questionWeight *int
		var choiceID, choiceContent *string
		var choiceWeight *int

		err := rows.Scan(
			&tryoutID,
			&tryoutCreatedAt,
			&tryoutUpdatedAt,
			&tryoutName,
			&tryoutType,
			&tryoutQuota,
			&tryoutPrice,
			&dueDate,
			&questionID,
			&questionContent,
			&questionWeight,
			&choiceID,
			&choiceContent,
			&choiceWeight,
		)
		if err != nil {
			return TryoutDetailResponse{}, err
		}

		// Populate tryout details (only once)
		if tryoutDetail.ID == "" {
			tryoutDetail.ID = tryoutID.String()
			tryoutDetail.CreatedAt = tryoutCreatedAt
			tryoutDetail.UpdatedAt = tryoutUpdatedAt
			tryoutDetail.TryOutName = tryoutName
			tryoutDetail.TestType = tryoutType
			tryoutDetail.TryOutQuota = tryoutQuota
			tryoutDetail.TryOutPrice = tryoutPrice
			tryoutDetail.DueDate = dueDate
		}

		if questionID != nil {
			question, exists := questionMap[*questionID]
			if !exists {
				question = &QuestionsResponse{
					ID:      *questionID,
					Content: *questionContent,
					Weight:  questionWeight,
					Choice:  []ChoiceResponse{},
				}
				questionMap[*questionID] = question
				tryoutDetail.Question = append(tryoutDetail.Question, *question)
			}

			if choiceID != nil {
				choice := ChoiceResponse{
					ID:      *choiceID,
					Content: *choiceContent,
					Weight:  choiceWeight,
				}
				question.Choice = append(question.Choice, choice)
				// Update the QuestionsResponse in tryoutDetail
				for i := range tryoutDetail.Question {
					if tryoutDetail.Question[i].ID == *questionID {
						tryoutDetail.Question[i].Choice = questionMap[*questionID].Choice
					}
				}
			}
		}

	}
	if err != nil {
		return
	}
	return tryoutDetail, err

}
