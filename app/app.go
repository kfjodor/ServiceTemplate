package app

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/spf13/viper"

	"service_template/infra"
	"service_template/logger"
	"service_template/middlewares"
	"service_template/storage"
)

type App struct {
	Router               *mux.Router
	DB                   *storage.Storage
	Infra                infra.Config
	keywalletRemoveAllow bool
}

func (a *App) Initialize(ctx context.Context) {
	log := logger.FromContext(ctx).WithField("m", "Initialize")
	log.Debugf("Initialize:: ")

	db := new(storage.Storage)
	err := db.InitPostgress(ctx, viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.name"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		nil)
	if err != nil {
		log.Errorf("Cannot connect to db connection: %v", err)

		return
	}

	a.DB = db
	a.Router = mux.NewRouter()
	a.setRouters()

	a.keywalletRemoveAllow = viper.GetBool("keywallet_remove_allow")
}

func (a *App) Run(infraCtx context.Context, host string) {
	log := logger.FromContext(infraCtx)
	log = log.WithField("m", "Run")
	log.Debugf("Run:: ")

	a.Router.Use(middlewares.LoggingMiddleware)
	a.Router.Use(middlewares.AuthMiddlewareGenerator(infraCtx, a.DB.DB))

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"localhost:3000"}),
		handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowCredentials(),
	)(a.Router)

	handler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tc := r.Context()
			r = r.WithContext(infra.GenerateNewTrace(tc))
			h.ServeHTTP(w, r)
		})
	}(cors)

	log.Errorf(infra.ServeHTTP(infraCtx, host, handler))
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) setRouters() {
}
