package bridge

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var serverMap = make(map[string]*http.Server)

type ResponseData struct {
	Status  int
	Headers map[string]string
	Body    string
}

type RequestData struct {
	Method string
	URL    string
	Header map[string][]string
	Body   string
}

func (a *App) StartServer(address string, serverID string) FlagResult {
	log.Printf("StartServer: %s", address)

	server := &http.Server{
		Addr: address,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(500)
				log.Printf("reading body error: %v", err)
				return
			}

			requestID := uuid.New().String()
			println(body, requestID)
			respChan := make(chan ResponseData)
			respBody := []byte{}

			defer close(respChan)

			a.Ctx.Events.On(requestID, func(event *application.WailsEvent) {
				a.Ctx.Events.Off(requestID)
				resp := ResponseData{200, make(map[string]string), "A sample http server"}
				println(event)
				// if len(event.Data) >= 4 {
				// 	if status, ok := event.Data[0].(float64); ok {
				// 		resp.Status = int(status)
				// 	}
				// 	if headers, ok := event.Data[1].(string); ok {
				// 		json.Unmarshal([]byte(headers), &resp.Headers)
				// 	}
				// 	if body, ok := event.Data[2].(string); ok {
				// 		resp.Body = body
				// 		respBody = []byte(body)
				// 	}
				// 	if options, ok := event.Data[3].(string); ok {
				// 		ioOptions := IOOptions{Mode: "Text"}
				// 		json.Unmarshal([]byte(options), &ioOptions)
				// 		if ioOptions.Mode == Binary {
				// 			body, err = base64.StdEncoding.DecodeString(resp.Body)
				// 			if err != nil {
				// 				resp.Status = 500
				// 				respBody = []byte(err.Error())
				// 			} else {
				// 				respBody = body
				// 			}
				// 		}
				// 	}
				// }
				respChan <- resp
			})

			a.Ctx.Events.Emit(&application.WailsEvent{
				Name: serverID,
				Data: &RequestData{
					Method: r.Method,
					URL:    r.URL.RequestURI(),
					Header: r.Header,
					Body:   string(body),
				},
			})

			res := <-respChan
			for key, value := range res.Headers {
				w.Header().Set(key, value)
			}
			w.WriteHeader(res.Status)
			w.Write(respBody)
		}),
	}

	var result error

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			result = err
		}
	}()

	time.Sleep(1 * time.Second)

	if result != nil {
		return FlagResult{false, result.Error()}
	}

	serverMap[serverID] = server

	return FlagResult{true, "Success"}
}

func (a *App) StopServer(id string) FlagResult {
	log.Printf("StopServer: %s", id)

	server, ok := serverMap[id]
	if !ok {
		return FlagResult{false, "server not found"}
	}

	err := server.Close()
	if err != nil {
		return FlagResult{false, err.Error()}
	}

	delete(serverMap, id)

	return FlagResult{true, "Success"}
}

func (a *App) ListServer() FlagResult {
	log.Printf("ListServer: %v", serverMap)

	var servers []string

	for serverID := range serverMap {
		servers = append(servers, serverID)
	}

	return FlagResult{true, strings.Join(servers, "|")}
}
