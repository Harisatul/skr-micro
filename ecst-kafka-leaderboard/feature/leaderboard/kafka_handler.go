package leaderboard

import (
	"context"
	"ecst-kafka-leaderboard/feature/shared"
	"ecst-kafka-leaderboard/pkg"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log/slog"
)

type CreateGradeHandler struct {
}

func (*CreateGradeHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_message_status"

		lvState2       = shared.LogEventStateInsertDB
		lfState2Status = "state_2_create_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("replicate-grade"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var payload gradeMessage
	err := json.Unmarshal(msg.Value, &payload)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(payload),
	)

	/*------------------------------------
	| Step 2 : Replicate submission to db
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	parse, err := uuid.Parse(payload.Id)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	tid, err := uuid.Parse(payload.TryoutID)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	grade := Grade{
		ID:       parse,
		Score:    payload.Score,
		TryoutID: tid,
		Token:    payload.Token,
		Version:  payload.Version,
	}

	err = insertGrade(ctx, grade)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))
	pkg.LogInfoWithContext(ctx, "success replicate grade in leaderboard service", lf)
}

type UpdateGradeHandler struct {
}

func (*UpdateGradeHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_message_status"

		lvState2       = shared.LogEventStateUpdateDB
		lfState2Status = "state_2_update_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("update-grade"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var payload gradeMessage
	err := json.Unmarshal(msg.Value, &payload)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(payload),
	)

	/*------------------------------------
	| Step 2 : Update submission to db
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	parse, err := uuid.Parse(payload.Id)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	grade := Grade{
		ID:      parse,
		Score:   payload.Score,
		Token:   payload.Token,
		Version: payload.Version,
	}

	err = updateGrade(ctx, grade)

	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))
	pkg.LogInfoWithContext(ctx, "success replicate grade in leaderboard service", lf)
}
