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

func GetPaymentDetailsHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := NewBillingResponse(nil)
		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		p, err := repo.FindByCreatedByID(currentUser.ID)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get your payment details %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		if p == nil {
			r.JSON(w, http.StatusNotFound, resp)
			return
		}

		resp.PaymentDetails = p
		r.JSON(w, http.StatusOK, resp)
	}
}

func SaveOrUpdatePaymentDetailsHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := NewBillingResponse(nil)
		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.PaymentDetails, err = repo.FindByCreatedByID(currentUser.ID)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get your payment details %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		if resp.PaymentDetails == nil {
			resp.PaymentDetails = &cohesioned.PaymentDetails{}
		}

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&resp.PaymentDetails); err != nil {
			resp.SetErrMsg("Unable to process payment details payload. Error: %v\n", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusBadRequest, resp)
			return
		}

		if resp.PaymentDetails.ID == 0 {
			fmt.Println("payment details id was not set - creating new payment details")
			resp.PaymentDetails.Created = time.Now()
			resp.PaymentDetails.CreatedBy = currentUser.ID

			resp.PaymentDetails.ID, err = repo.Save(resp.PaymentDetails)
			if err != nil {
				resp.SetErrMsg("Failed to save payment details %v", err)
				fmt.Println(resp.ErrMsg)
				r.JSON(w, http.StatusInternalServerError, resp)
				return
			}
		} else {
			fmt.Println("payment details id was set - updating payment details")
			resp.PaymentDetails.UpdatedBy = currentUser.ID
			resp.PaymentDetails.Updated = time.Now()

			if err := repo.Update(resp.PaymentDetails); err != nil {
				resp.SetErrMsg("Failed to update payment details %v", err)
				fmt.Println(resp.ErrMsg)
				r.JSON(w, http.StatusInternalServerError, resp)
				return
			}
		}

		r.JSON(w, http.StatusOK, resp)
	}
}
