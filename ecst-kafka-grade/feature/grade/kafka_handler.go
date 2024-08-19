package grade

import (
	"context"
	"ecst-kafka-grade/feature/shared"
	"ecst-kafka-grade/pkg"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log/slog"
)

type CreateSubmissionHandler struct {
}

func (*CreateSubmissionHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_message_status"

		lvState2       = shared.LogEventStateUpdateDB
		lfState2Status = "state_2_create_db_status"

		lvState3       = shared.LogEventStateKafkaPublish
		lfState3Status = "state_3_publish_message_status"

		lf = []slog.Attr{
			pkg.LogEventName("replicate-submission"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var payload submissionMessage
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
	| Step 3 : Replicate submission to db
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	parse, err := uuid.Parse(payload.Id)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	submission := UserTestSubmission{
		ID:      parse,
		Token:   payload.Token,
		Version: payload.Version,
	}

	id, version, token, score, err := insertSubmission(ctx, submission)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	/*------------------------------------
	| Step 4 : Publish message
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	v := gradeMessage{
		Id:       id.String(),
		Token:    token,
		TryoutID: payload.TryOutId,
		Score:    score,
		Version:  version,
	}
	message, err := json.Marshal(v)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	err = pkg.PublishMessage(kp, CreateGradeTopic, string(message))

	lf = append(lf, pkg.LogStatusSuccess(lfState3Status))

	pkg.LogInfoWithContext(ctx, "success replicate submission in grade service", lf)
}

type CreateAnswerHandler struct {
}

func (*CreateAnswerHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_message_status"

		//lvState2       = shared.LogEventStateUpdateDB
		//lfState2Status = "state_2_create_db_status"

		lvState3       = shared.LogEventStateFetchDB
		lfState3Status = "state_3_fetch_db_status"

		lvState4       = shared.LogEventStateUpdateDB
		lfState4Status = "state_3_update_db_status"

		lvState5       = shared.LogEventStateKafkaPublish
		lfState5Status = "state_3_publish_message_status"

		lf = []slog.Attr{
			pkg.LogEventName("replicate-answer"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var payload answerMessage
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

	///*------------------------------------
	//| Step 2 : Replicate answer to db
	//* ----------------------------------*/
	//lf = append(lf, pkg.LogEventState(lvState2))
	//
	//parse, err := uuid.Parse(payload.Id)
	//if err != nil {
	//	lf = append(lf, pkg.LogStatusFailed(lfState2Status))
	//	pkg.LogErrorWithContext(ctx, err, lf)
	//	return
	//}
	//sparse, err := uuid.Parse(payload.SubmissionID)
	//if err != nil {
	//	lf = append(lf, pkg.LogStatusFailed(lfState2Status))
	//	pkg.LogErrorWithContext(ctx, err, lf)
	//	return
	//}
	//cparse, err := uuid.Parse(payload.ChoiceId)
	//if err != nil {
	//	lf = append(lf, pkg.LogStatusFailed(lfState2Status))
	//	pkg.LogErrorWithContext(ctx, err, lf)
	//	return
	//}
	//
	//answer := UserSubmittedAnswer{
	//	ID:                   parse,
	//	UserTestSubmissionID: sparse,
	//	ChoiceID:             cparse,
	//	Version:              payload.Version,
	//}

	//_, _, err = insertAnswer(ctx, answer)
	//if err != nil {
	//	lf = append(lf, pkg.LogStatusFailed(lfState2Status))
	//	pkg.LogErrorWithContext(ctx, err, lf)
	//	return
	//}
	//
	//lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	/*------------------------------------
	| Step 3 : Fetch Choice Data
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	idr, err := fetchChoices(ctx, payload.ChoiceId)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState3Status))

	/*------------------------------------
	| Step 4 : Update Submission Data
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState4))

	var (
		submission uuid.UUID
		v          int64
		token      string
		score      int
	)
	if idr == true {
		submission, v, token, score, err = updateSubmission(ctx, payload.SubmissionID)
		if err != nil {
			lf = append(lf, pkg.LogStatusFailed(lfState4Status))
			pkg.LogErrorWithContext(ctx, err, lf)
			return
		}
	} else {
		submission, v, token, score, err = updateSubmissionW(ctx, payload.SubmissionID)
		if err != nil {
			lf = append(lf, pkg.LogStatusFailed(lfState4Status))
			pkg.LogErrorWithContext(ctx, err, lf)
			return
		}
	}
	lf = append(lf, pkg.LogStatusSuccess(lfState3Status))

	/*------------------------------------
	| Step 5 : Publish message
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState5))

	message, err := json.Marshal(gradeMessage{
		Id:      submission.String(),
		Token:   token,
		Score:   score,
		Version: v,
	})
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState5Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	err = pkg.PublishMessage(kp, UpdateGradeTopic, string(message))

	lf = append(lf, pkg.LogStatusSuccess(lfState5Status))

	pkg.LogInfoWithContext(ctx, "success replicate submission in grade service", lf)
}
