package leaderboard

import (
	"ecst-kafka-leaderboard/feature/shared"
	"ecst-kafka-leaderboard/pkg"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

func HttpRoute(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/http/get", getGrade)
}

func getGrade(w http.ResponseWriter, r *http.Request) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateFetchDB
		lfState2Status = "state_2_fetch_grade"

		ctx = r.Context()

		lf = []slog.Attr{
			pkg.LogEventName("get-leaderboard"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var req leaderboardReq
	// Retrieve the query parameter
	tryoutID := r.URL.Query().Get("tid")
	sizeStr := r.URL.Query().Get("size")
	pageStr := r.URL.Query().Get("page")

	if tryoutID == "" {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx, "missing id query parameter", nil, lf)
		shared.WriteErrorResponse(w, http.StatusBadRequest, errors.New("missing id query parameter"))
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx, "failed to cast size query param", nil, lf)
		shared.WriteErrorResponse(w, http.StatusBadRequest, errors.New("failed to cast size query param"))
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx, "failed to cast page query param", nil, lf)
		shared.WriteErrorResponse(w, http.StatusBadRequest, errors.New("failed to cast page query param"))
		return
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(req),
	)
	/*------------------------------------
	| Step 2 : Fetch Leaderboard
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	offset := (page - 1) * size

	grade, err := fetchGrade(ctx, tryoutID, size, offset)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}
	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	shared.WriteSuccessResponse(w, http.StatusOK,
		grade,
	)
	pkg.LogInfoWithContext(ctx, "success get leaderboard", lf)
}
