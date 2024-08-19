package tryout

import (
	"ecst-kafka-tryout/feature/shared"
	"ecst-kafka-tryout/pkg"
	"errors"
	"log/slog"
	"net/http"
)

func HttpRoute(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/tryout/detail", detailTryoutHandler)
}

func detailTryoutHandler(w http.ResponseWriter, r *http.Request) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateFetchDB
		lfState2Status = "state_2_get_detail_tryout"

		ctx = r.Context()

		lf = []slog.Attr{
			pkg.LogEventName("DetailTryout"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var req detailTryoutRequest
	// Retrieve the query parameter
	tryoutID := r.URL.Query().Get("id")
	if tryoutID == "" {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx, "missing id query parameter", nil, lf)
		shared.WriteErrorResponse(w, http.StatusBadRequest, errors.New("missing id query parameter"))
		return
	}
	req.ID = tryoutID

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(req),
	)

	/*------------------------------------
	| Step 2 : GET TRYOUT
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	tryout, err := getTryoutWithQuestions(ctx, req.ID)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))
	shared.WriteSuccessResponse(w, http.StatusOK,
		tryout,
	)
	pkg.LogInfoWithContext(ctx, "success get tryout", lf)
}
