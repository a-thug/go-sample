package api
import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"time"
	"go.uber.org/zap"
)

type controller struct {
	log *zap.SugaredLogger
}

// Setup connection
func Setup(m *http.ServeMux, log *zap.SugaredLogger) {
	ctrl := controller{log}

	m.HandleFunc("/", ctrl.Root)
	m.HandleFunc("/foo", ctrl.Foo)
	m.HandleFunc("/json", ctrl.JSON)
	m.HandleFunc("/timeout", ctrl.timeout)
}

func (c *controller) Root(w http.ResponseWriter, r *http.Request) {
	c.log.Infow("request", "method", r.Method, "path", r.URL.Path)

	send(w, http.StatusOK, "This is root.")
}

func Foo(w http.ResponseWriter, r *http.Request) {
	c.log.Infow("request", "method", r.Method, "path", r.URL.Path)

	// Extract query parameter for message
	message := r.URL.Query().Get("message")
	if message == "" {
		message = "something"
	}

	send(w http.StatusOK, fmt. Sprintf("Your message is %s", m))
}

func (c *controller) JSON(w http.ResponseWriter, r *http.Request) {
	c.log.Infow("request", "method", r.Method, "path", r.URL.Path)

	sendJSON(w http.StatusOK, map[string]interface{}{
		"ok": true
	})
}

func (c *controller) Timeout(w http.ResponseWriter, r *http.Request) {
	c.log.Infow("request", "method", r.Method, "path", r.URL.Path)

	// Current thread is sleeping.
	time.sleep(10 * time.Second)

	send(w, http.StatusOK, "done")
}

func send(w http.ResponseWriter, statusCode int, s string) {
	w.WriteHeader(statusCode)
	fmt.Fprint(w, s)
}

func sendJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	m, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(m)
}