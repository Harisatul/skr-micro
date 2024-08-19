package tryout

import (
	"ecst-kafka-exam/feature/shared"
	"ecst-kafka-exam/pkg"
	"encoding/json"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func HttpRoute(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/submission/create", insertSubmissionHandler)
	mux.HandleFunc("POST /api/answer/create", insertAnswerHandler)
}

func insertSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateInsertDB
		lfState2Status = "state_2_insert_submission`"

		lvState3       = shared.LogEventStateKafkaPublish
		lfState3Status = "state_3_publish_message_status"

		ctx = r.Context()

		lf = []slog.Attr{
			pkg.LogEventName("insert-submission"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var req submissionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx, "invalid request", err, lf)
		shared.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(req),
	)
	/*------------------------------------
	| Step 2 : Insert Submission
	* ----------------------------------*/

	id, err := pkg.GenerateId()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	lf = append(lf, pkg.LogEventState(lvState2))

	parse, err := uuid.Parse(req.TryoutID)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	submission := UserTestSubmission{
		ID:       id,
		Token:    req.Token,
		TryoutID: parse,
	}

	sid, version, token, tryout, err := insertSubmission(ctx, submission)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}
	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	/*------------------------------------
	| Step 3 : Publish message
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	message := submissionMessage{
		Id:       sid.String(),
		TryOutID: tryout.String(),
		Token:    token,
		Version:  version,
	}

	messageByte, err := json.Marshal(message)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	err = pkg.PublishMessage(kp, CreateSubmissionTopic, string(messageByte))
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState3Status))
	shared.WriteSuccessResponse(w, http.StatusOK,
		submissionResponse{
			Id: sid,
		},
	)
	pkg.LogInfoWithContext(ctx, "success get tryout", lf)
}

func insertAnswerHandler(w http.ResponseWriter, r *http.Request) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateInsertDB
		lfState2Status = "state_2_insert_answer`"

		lvState3       = shared.LogEventStateKafkaPublish
		lfState3Status = "state_3_publish_message_status"

		ctx = r.Context()

		lf = []slog.Attr{
			pkg.LogEventName("insert-answer"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var req answerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx, "invalid request", err, lf)
		shared.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(req),
	)
	/*------------------------------------
	| Step 2 : Insert Answer
	* ----------------------------------*/

	id, err := pkg.GenerateId()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	lf = append(lf, pkg.LogEventState(lvState2))

	sparse, err := uuid.Parse(req.SubmissionID)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	qparse, err := uuid.Parse(req.QuestionID)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	cparse, err := uuid.Parse(req.ChoiceID)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	answer := UserSubmittedAnswer{
		ID:                   id,
		QuestionID:           qparse,
		ChoiceID:             cparse,
		UserTestSubmissionID: sparse,
	}

	aid, version, choiceId, submissionId, err := insertAnswer(ctx, answer)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}
	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	/*------------------------------------
	| Step 3 : Publish message
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	message := answerMessage{
		Id:           aid.String(),
		SubmissionID: submissionId.String(),
		ChoiceId:     choiceId.String(),
		Version:      version,
	}

	messageByte, err := json.Marshal(message)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	err = pkg.PublishMessage(kp, CreateAnswerTopic, string(messageByte))
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState3Status))
	shared.WriteSuccessResponse(w, http.StatusOK,
		answerResponse{
			Id: aid,
		},
	)
	pkg.LogInfoWithContext(ctx, "success get tryout", lf)
}
