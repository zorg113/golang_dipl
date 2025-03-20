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

func NewRouter(auth *handlers.Authorization, balckList *handlers.BlackList, whiteList *handlers.WhiteList, bucket *handlers.Bucket, log *zerolog.Logger) *HttpApiRouter {
	router := mux.NewRouter()
	return &HttpApiRouter{
		router:    router,
		auth:      auth,
		blackLits: balckList,
		whiteList: whiteList,
		bucket:    bucket,
		log:       log,
	}
}

func (r *HttpApiRouter) InitRouters() {
	r.router.HandleFunc("/auth/check", r.auth.AuthorizationHanler).Methods("POST")
	r.router.HandleFunc("/auth/reset", r.bucket.ResetBucket).Methods("DELETE")
	r.router.HandleFunc("/auth/blacklist", r.blackLits.AddIP).Methods("POST")
	r.router.HandleFunc("/auth/blacklist", r.blackLits.GetIPs).Methods("GET")
	r.router.HandleFunc("/auth/blacklist", r.blackLits.DeleteIP).Methods("DELETE")
	r.router.HandleFunc("/auth/whitelist", r.whiteList.AddIP).Methods("POST")
	r.router.HandleFunc("/auth/whitelist", r.whiteList.GetIPs).Methods("GET")
	r.router.HandleFunc("/auth/whitelist", r.whiteList.DeleteIP).Methods("DELETE")
}

func (r *HttpApiRouter) GetRouter() *mux.Router {
	return r.router
}
