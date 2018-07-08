package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/ansel1/merry"
	"github.com/eloff/scheduler/server/data"
)

type ServiceResource interface {
	Register() error
}

func APIServer(ds data.DataStore, address string) (*http.Server, error) {
	services := []ServiceResource{
		&ScheduleResource{ds: ds},
	}

	server := &http.Server{Addr: address}
	for _, service := range services {
		err := service.Register()
		if err != nil {
			return nil, err
		}
	}

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("../client/src/static"))))
	http.Handle("/", http.FileServer(http.Dir("../client/build")))
	/*target, err := url.Parse("http://127.0.0.1:3000")
	if err != nil {
		return nil, merry.Wrap(err)
	}
	http.Handle("/", httputil.NewSingleHostReverseProxy(target))*/

	return server, nil
}

type JSONHandler func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func handlerWrapper(handler JSONHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := handler(w, r)
		if err == nil {
			if result != nil {
				if url, ok := result.(*url.URL); ok {
					// If result is a URL, then this is a redirect
					http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
					return
				}

				var jsonText []byte
				jsonText, err = json.Marshal(result)
				if err == nil {
					w.Header().Set("Content-Type", "application/json")
					_, err = w.Write(jsonText)
				}
				err = merry.Wrap(err)
			}
		}
		if err != nil {
			code := merry.HTTPCode(err)
			w.WriteHeader(code)
			log.Printf("error in request: %s", merry.Details(err))
		}
	})
}
