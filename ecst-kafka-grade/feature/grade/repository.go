package grade

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

func insertSubmission(ctx context.Context, submission UserTestSubmission) (id uuid.UUID, version int64, token string, score int, err error) {

	// Prepare the SQL insert statement
	sql := `INSERT INTO submission (id, created_at, updated_at, version, token, score) 
            VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, version, token, score`

	// Execute the insert statement
	err = db.QueryRow(ctx, sql, submission.ID, time.Now().UnixMilli(), time.Now().UnixMilli(), submission.Version, submission.Token, 0).Scan(&id, &version, &token, &score)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}
	}
	return
}

func insertAnswer(ctx context.Context, answer UserSubmittedAnswer) (id uuid.UUID, version int64, err error) {

	sql := `INSERT INTO answer_submitted (id, created_at, updated_at, choice_id, user_test_submission_id, version) 
            VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, version`

	// Execute the insert statement
	err = db.QueryRow(ctx, sql, answer.ID, time.Now().UnixMilli(), time.Now().UnixMilli(), answer.ChoiceID, answer.UserTestSubmissionID, answer.Version).Scan(&id, &version)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}
	}
	return
}

func fetchChoices(ctx context.Context, id string) (idr bool, err error) {

	query := "SELECT is_correct FROM choices WHERE id = $1"

	// Execute the insert statement
	err = db.QueryRow(ctx, query, id).Scan(&idr)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}
	}
	return
}

func updateSubmission(ctx context.Context, submission string) (id uuid.UUID, version int64, token string, score int, err error) {

	// Prepare the SQL update statement
	sql := `
		UPDATE submission 
		SET updated_at = $1, version = version + 1,  score = score + 5
		WHERE id = $2
		RETURNING id, version, token, score
	`

	// Execute the update statement
	err = db.QueryRow(ctx, sql, time.Now().UnixMilli(), submission).Scan(&id, &version, &token, &score)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}
	}
	return
}

func updateSubmissionW(ctx context.Context, submission string) (id uuid.UUID, version int64, token string, score int, err error) {

	// Prepare the SQL update statement
	sql := `
		UPDATE submission 
		SET updated_at = $1, version = version + 1,  score = score + 0
		WHERE id = $2
		RETURNING id, version, token, score
	`

	// Execute the update statement
	err = db.QueryRow(ctx, sql, time.Now().UnixMilli(), submission).Scan(&id, &version, &token, &score)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}
	}
	return
}
