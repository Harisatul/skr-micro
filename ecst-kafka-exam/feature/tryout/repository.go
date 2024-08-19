package tryout

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

func insertSubmission(ctx context.Context, submission UserTestSubmission) (id uuid.UUID, version int64, token string, tryoutid uuid.UUID, err error) {

	// Prepare the SQL insert statement
	sql := `INSERT INTO submission (id, created_at, updated_at, submission_start_time, version, submission_end_time, token, tryout_id) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, version, token, tryout_id`

	// Execute the insert statement
	err = db.QueryRow(ctx, sql, submission.ID, time.Now().UnixMilli(), time.Now().UnixMilli(), time.Now(), 1, nil, submission.Token, submission.TryoutID).Scan(&id, &version, &token, &tryoutid)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}
	}
	return
}

func insertAnswer(ctx context.Context, answer UserSubmittedAnswer) (id uuid.UUID, version int64, choice_id uuid.UUID, user_test_submission_id uuid.UUID, err error) {

	sql := `INSERT INTO answer_submitted (id, created_at, updated_at, question_id, choice_id, user_test_submission_id, version) 
            VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, version, choice_id, user_test_submission_id`

	// Execute the insert statement
	err = db.QueryRow(ctx, sql, answer.ID, time.Now().UnixMilli(), time.Now().UnixMilli(), answer.QuestionID, answer.ChoiceID, answer.UserTestSubmissionID, 1).Scan(&id, &version, &choice_id, &user_test_submission_id)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}
	}
	return
}
