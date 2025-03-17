package httpapi

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/httpapi/handlers"
)

type HttpApiRouter struct {
	router    *mux.Router
	auth      *handlers.Authorization
	blackLits *handlers.BlackList
	whiteList *handlers.WhiteList
	bucket    *handlers.Bucket
	log       *zerolog.Logger
}

func NewRouter(auth *handlers.Authorization, balckList *handlers.BlackList, whiteList *handlers.WhiteList, log *zerolog.Logger) *HttpApiRouter {
	router := mux.NewRouter()
	return &HttpApiRouter{
		router:    router,
		auth:      auth,
		blackLits: balckList,
		whiteList: whiteList,
		log:       log,
	}
}

func (r *HttpApiRouter) InitRoutes() {
	r.router.HandleFunc("/auth/check", r.auth.AuthorizationHanler).Methods("POST")
	r.router.HandleFunc("/bucket/reset", r.bucket.ResetBucket).Methods("DELETE")
	r.router.HandleFunc("/blacklist/add", r.blackLits.AddIP).Methods("POST")
	r.router.HandleFunc("/blacklist/get", r.blackLits.GetIPs).Methods("GET")
	r.router.HandleFunc("/blacklist/remove", r.blackLits.DeleteIP).Methods("DELETE")
	r.router.HandleFunc("/whitelist/add", r.whiteList.AddIP).Methods("POST")
	r.router.HandleFunc("/whitelist/get", r.whiteList.GetIPs).Methods("GET")
	r.router.HandleFunc("/whitelist/remove", r.whiteList.DeleteIP).Methods("DELETE")
}

func (r *HttpApiRouter) GetRouter() *mux.Router{
	return r.router
}