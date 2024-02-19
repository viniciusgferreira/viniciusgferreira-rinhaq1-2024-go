package handlers

import (
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/data"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/services"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Client struct{}
type ClientsHandler struct {
	service *services.Service
}

func NewClientHandler(service *services.Service) *ClientsHandler {
	return &ClientsHandler{service: service}
}

func (c *ClientsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathID(r)
	if !ok {
		w.WriteHeader(400)
		return
	}
	if id > 5 || id < 1 {
		w.WriteHeader(404)
		return
	}

	if r.Method == http.MethodPost {
		c.createTransaction(w, r, id)
	}

	if r.Method == http.MethodGet {
		c.createStatement(w, r, id)
	}
}

func parsePathID(r *http.Request) (uint8, bool) {
	reg := regexp.MustCompile(`(\d+)`)
	g := reg.FindAllStringSubmatch(r.URL.Path, -1)
	if len(g) != 1 || len(g[0]) != 2 {
		return 0, false
	}
	id, err := strconv.Atoi(g[0][1])
	if err != nil {
		return 0, false
	}
	return uint8(id), true
}

func (c *ClientsHandler) createTransaction(w http.ResponseWriter, r *http.Request, id uint8) {
	t, err := data.NewTransaction(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if !t.Validate() {
		w.WriteHeader(400)
		return
	}

	newBalance, err := c.service.CreateTransaction(t, id)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(422)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_ = newBalance.ToJSON(w)
}

func (c *ClientsHandler) createStatement(w http.ResponseWriter, r *http.Request, id uint8) {
	client := c.service.CreateStatement(id)
	_ = client.ToJSON(w)
}
