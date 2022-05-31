package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

/*

A Go web server comprises one function - a `net/http.HandlerFunc(ResponseWriter, *Request)` - for each of
your API's endpoints. Our API has two endpoints: Produce for writing to the log and Consume for reading
from the log. When building a JSON/HTTP Go server, each handler consists of three steps:

1. Unmarshal the request's JSON body into a struct.
2. Run that endpoint's logic with the request to obtain a result.
3. Marshal and write that result to the response.

If your handlers become much more complicated than this, then you should move the code out, move request
and response handling into HTTP middleware, and move the business logic further down the stack.

*/

// NewHTTPServer takes in an address for the server to run on and returns an *http.Server.
func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()

	// Using gorilla/mux library to write nice, RESTful routes that match incoming requests
	// to their respective handlers.
	r := mux.NewRouter()
	r.HandleFunc("/", httpsrv.handleProduce).Methods(http.MethodPost)
	r.HandleFunc("/", httpsrv.handleConsume).Methods(http.MethodGet)

	// We wrap our server with a *http.Server so the user just needs to call ListenAndServe()
	// to listen for and handle incoming requests.
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

// httpServer references a log for the server to defer to in its handlers.
type httpServer struct {
	Log *Log
}

// newHTTPServer creates a new server containing a reference to a log.
func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

// ProduceRequest contains the record that the caller of our API wants appended to the log.
type ProduceRequest struct {
	Record Record `json:"record"`
}

// ProduceResponse tells the caller what offset the log stored the records under.
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

// ConsumeRequest specifies which records the caller of our API wants to read.
type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

// ConsumeResponse specifies the records to send back to the caller.
type ConsumeResponse struct {
	Record Record `json:"record"`
}

// handleProduce implements the three steps that a handler should consist of: unmarshling the
// request into a struct, using that struct to produce to the log and getting the offset that
// the log stored the record under, and marshaling and writing the result to the response.
func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ProduceResponse{Offset: off}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ConsumeResponse{Record: record}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
