package permissions

import (
	"net/http"

	"log"

	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/retrieval"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type Middleware struct {
	*handler.State
	T          Permission
	TargetType uint32
	M          func(state *handler.State, p Permission, TargetType uint32, next http.Handler) http.Handler
}

func (m Middleware) Handler(next http.Handler) http.Handler {
	return m.M(m.State, m.T, m.TargetType, next)
}

func PermissionMiddle(state *handler.State, p Permission, TargetType uint32, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var usr interface{}
		usr = model.Ref{}
		id, _ := mux.Vars(r)["ID"]
		if p != CanView {
			usr, ok := context.GetOk(r, "auth")
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Auth params not set")
				return
			}
			if usr == nil {
				log.Println("User is nil")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		var tarRef model.Ref
		var err error

		switch TargetType {
		case model.Images:
			tarRef, err = retrieval.GetImageRef(state.DB, id)
		case model.Users:
			tarRef, err = retrieval.GetUserRef(state.DB, id)
		}

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		user := usr.(model.Ref)
		valid, err := Valid(state.DB, user.Id, p, tarRef.Id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !valid && p != CanView {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}