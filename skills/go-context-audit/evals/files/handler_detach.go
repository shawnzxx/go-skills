package auditfixture

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Mailer struct{}

func (m *Mailer) Send(ctx context.Context, to, body string) error {
	time.Sleep(2 * time.Second)
	fmt.Printf("sent email to %s\n", to)
	return nil
}

type OrderHandler struct {
	mailer *Mailer
}

func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	userEmail := r.FormValue("email")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := processOrder(ctx); err != nil {
		http.Error(w, "order failed", http.StatusInternalServerError)
		return
	}

	go h.mailer.Send(r.Context(), userEmail, "Your order has been placed!")

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "order placed")
}

func processOrder(ctx context.Context) error {
	return nil
}
