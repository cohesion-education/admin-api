package billing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/unrolled/render"
)

type BillingResponse struct {
	*cohesioned.APIResponse
	*cohesioned.PaymentDetails
}

func NewBillingResponse(p *cohesioned.PaymentDetails) *BillingResponse {
	return &BillingResponse{
		PaymentDetails: p,
		APIResponse:    &cohesioned.APIResponse{},
	}
}

func SavePaymentDetailsHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		p := &cohesioned.PaymentDetails{}
		resp := NewBillingResponse(p)

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&p); err != nil {
			resp.SetErrMsg("Unable to process payment details payload. Error: %v\n", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusBadRequest, resp)
			return
		}

		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		p.CreatedBy = currentUser.ID
		p.Created = time.Now()

		id, err := repo.Save(p)
		p.ID = id
		if err != nil {
			resp.SetErrMsg("Failed to save payment details %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		r.JSON(w, http.StatusOK, resp)
	}
}
