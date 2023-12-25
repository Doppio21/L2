package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Address string
}

type Dependencies struct {
	Events *Events
	Log    *logrus.Logger
}

type Server struct {
	cfg  Config
	deps Dependencies
}

func New(cfg Config, deps Dependencies) *Server {
	return &Server{
		cfg:  cfg,
		deps: deps,
	}
}

func (s *Server) Run(ctx context.Context) {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    s.cfg.Address,
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	logMiddlewares := LoggerMiddleware(s.deps.Log.WithContext(ctx))

	mux.Handle("/create_event", logMiddlewares(wrapHandler(http.MethodPost, s.handlePost("create"))))
	mux.Handle("/update_event", logMiddlewares(wrapHandler(http.MethodPost, s.handlePost("update"))))
	mux.Handle("/delete_event", logMiddlewares(wrapHandler(http.MethodPost, s.handlePost("delete"))))
	mux.Handle("/events_for_day", logMiddlewares(wrapHandler(http.MethodGet, s.handleGet("day"))))
	mux.Handle("/events_for_week", logMiddlewares(wrapHandler(http.MethodGet, s.handleGet("week"))))
	mux.Handle("/events_for_month", logMiddlewares(wrapHandler(http.MethodGet, s.handleGet("month"))))

	serverClosed := make(chan struct{})
	go func() {
		defer close(serverClosed)
		s.deps.Log.Info("server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.deps.Log.Fatalf("listen and serve: %v", err)
		}
	}()

	select {
	case <-serverClosed:
	case <-ctx.Done():
		s.deps.Log.Info("server stoped")
		_ = srv.Close()
	}
}

func getQueries(r *http.Request) (int, string, error) {
	q := r.URL.Query()
	userID := q.Get("userid")
	if userID == "" {
		return 0, "", errors.New("invalid userID")
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, "", err
	}

	date := q.Get("date")
	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		return 0, "", errors.New("invalid date")
	}

	return id, date, nil
}

func writeError(err error, w http.ResponseWriter) {
	e, _ := json.Marshal(fmt.Sprintf("Error: " + err.Error()))
	w.Write(e)
}

func wrapHandler(method string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (s *Server) handlePost(pattern string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userid, date, err := getQueries(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeError(err, w)
			s.deps.Log.Errorf("error get queries: %v", err)
			return
		}

		buff := make([]byte, 1024)
		n, err := r.Body.Read(buff)
		if err != nil && err != io.EOF {
			w.WriteHeader(http.StatusBadRequest)
			writeError(err, w)
			s.deps.Log.Errorf("error reading request: %v", err)
			return
		}
		buff = buff[:n]
		req := &GetRequest{}
		err = json.Unmarshal(buff, req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeError(err, w)
			s.deps.Log.Errorf("error json unmarshal: %v", err)
			return
		}
		switch pattern {
		case "create":
			err := s.deps.Events.Create(userid, date, req.Name, req.Desc)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				writeError(err, w)
				s.deps.Log.Errorf("event creation error: %v", err)
				return
			}
		case "update":
			err := s.deps.Events.Update(userid, date, req.Name, req.Desc)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				writeError(err, w)
				s.deps.Log.Errorf("event update error: %v", err)
				return
			}
		case "delete":
			err := s.deps.Events.Delete(userid, date, req.Name)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				writeError(err, w)
				s.deps.Log.Errorf("event deletion error: %v", err)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (s *Server) handleGet(pattern string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userid, date, err := getQueries(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeError(err, w)
			s.deps.Log.Errorf("error get queries: %v", err)
			return
		}

		events := []Event{}
		switch pattern {
		case "day":
			events, err = s.deps.Events.Get(userid, date, "day")
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				writeError(err, w)
				s.deps.Log.Errorf("error get events for day: %v", err)
				return
			}
		case "week":
			events, err = s.deps.Events.Get(userid, date, "week")
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				writeError(err, w)
				s.deps.Log.Errorf("error get events for week: %v", err)
				return
			}
		case "month":
			events, err = s.deps.Events.Get(userid, date, "month")
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				writeError(err, w)
				s.deps.Log.Errorf("error get events for month: %v", err)
				return
			}
		}

		if events == nil {
			w.WriteHeader(http.StatusNotFound)
			writeError(errors.New("event not found"), w)
			s.deps.Log.Error("event not found")
			return
		}

		resp := Response{Events: events}
		data, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeError(err, w)
			s.deps.Log.Errorf("error json marshal: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(data); err != nil {
			writeError(err, w)
			s.deps.Log.Errorf("error writing response: %v", err)
			return
		}
	})
}
