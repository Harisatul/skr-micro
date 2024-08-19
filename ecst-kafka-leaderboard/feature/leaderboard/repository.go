package leaderboard

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

func fetchGrade(ctx context.Context, tid string, limit, offset int) ([]ParticipantScore, error) {

	var participants []ParticipantScore

	// SQL query with placeholders for LIMIT and OFFSET
	sql := `
        SELECT
            p.id AS participant_id,
            p.participants_name AS participant_name,
            g.score
        FROM
            grades g
        JOIN
            participants p ON g.token = p.token
       WHERE
            g.tryout_id = $1
        ORDER BY
            g.score DESC
        LIMIT $2
        OFFSET $3
    `

	rows, err := db.Query(ctx, sql, tid, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var participantID string
		var participantName string
		var score int

		err := rows.Scan(&participantID, &participantName, &score)
		if err != nil {
			return nil, err
		}

		participant := ParticipantScore{
			ParticipantID:   participantID,
			ParticipantName: participantName,
			Score:           score,
		}
		participants = append(participants, participant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return participants, nil
}

func insertGrade(ctx context.Context, grade Grade) error {
	sql := `
        INSERT INTO grades (id, created_at, updated_at,tryout_id, score, version, token)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := db.Exec(ctx, sql, grade.ID, time.Now().UnixMilli(), time.Now().UnixMilli(), grade.TryoutID, grade.Score, grade.Version, grade.Token)
	if err != nil {
		return err
	}

	return nil
}

func updateGrade(ctx context.Context, grade Grade) error {

	sql := `
        UPDATE grades
        SET updated_at = $1, score = $2, version = $3
        WHERE id = $4 AND version = $5
    `

	tag, err := db.Exec(ctx, sql, time.Now().UnixMilli(), grade.Score, grade.Version, grade.ID, grade.Version-1)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return err
}
