package ui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/a-h/templ"
	"github.com/orlandokj/just/application"
	"github.com/orlandokj/just/config"
	"github.com/orlandokj/just/ui/templates"
)

func isHtmxRequest(r *http.Request) bool {
    return r.Header.Get("hx-request") == "true"
}

func isSSERequest(r *http.Request) bool {
    // text/event-stream
    return r.Header.Get("Accept") == "text/event-stream"
}

func RunUI() error {
    mux := http.NewServeMux()
    mux.HandleFunc("/", createComponentHandleFunc(func(r *http.Request) (templ.Component, error) {
        var component templ.Component
        if (isHtmxRequest(r)) {
            component = templates.Dashboard()
        } else {
            component = templates.Index(templates.Dashboard())
        }
        return component, nil
    }))
    
    mux.HandleFunc("/dashboard", createComponentHandleFunc(func(r *http.Request) (templ.Component, error) {
        return templates.Dashboard(), nil
    }))

    mux.HandleFunc("/new-application", createComponentHandleFunc(func(r *http.Request) (templ.Component, error) {
        if r.Method == "GET" {
            return templates.ApplicationForm(application.Application{}), nil
        }

        if r.Method == "POST" {
            appConfig := config.Config{} 
            err := json.NewDecoder(r.Body).Decode(&appConfig)
            log.Println(appConfig)
            if err != nil {
                log.Printf("Error decoding new application: %v", err)
                return nil, err
            }
            err = application.NewApplication(appConfig)
            if err != nil {
                return nil, err
            }
            return templates.Dashboard(), nil
        }
        return nil, nil
    }))

    mux.HandleFunc("/edit-application/{name}", createComponentHandleFunc(func(r *http.Request) (templ.Component, error) {
        if r.Method == "GET" {
            app := application.GetApplication(r.PathValue("name"))
            if app == nil {
                return nil, errors.New("Application not found")
            }
            return templates.ApplicationForm(*app), nil
        }

        if r.Method == "POST" {
            appConfig := config.Config{} 
            err := json.NewDecoder(r.Body).Decode(&appConfig)
            log.Println(appConfig)
            if err != nil {
                log.Printf("Error decoding new application: %v", err)
                return nil, err
            }
            err = application.ModifyApplication(appConfig)
            if err != nil {
                return nil, err
            }
            return templates.Dashboard(), nil
        }
        return nil, nil
    }))

    mux.HandleFunc("/application/{name}/start", createComponentHandleFunc(func(r *http.Request) (templ.Component, error) {
        if r.Method != "POST" {
            return nil, nil
        }

        name := r.PathValue("name")
        err := application.RunApplication(name)
        if err != nil {
            log.Printf("Error starting application: %v", err)
            return nil, err
        }

        return templates.Dashboard(), nil
    }))

    mux.HandleFunc("/application/{name}/stop", createComponentHandleFunc(func(r *http.Request) (templ.Component, error) {
        if r.Method != "POST" {
            return nil, nil
        }

        name := r.PathValue("name")
        err := application.StopApplication(name)
        if err != nil {
            log.Printf("Error starting application: %v", err)
            return nil, err
        }

        return templates.Dashboard(), nil
    }))

    mux.HandleFunc("/application/{name}/delete", createComponentHandleFunc(func(r *http.Request) (templ.Component, error) {
        if r.Method != "DELETE" {
            return nil, nil
        }

        name := r.PathValue("name")
        err := application.DeleteApplication(name)
        if err != nil {
            return nil, err
        }

        return templates.Dashboard(), nil
    }))

    mux.HandleFunc("/application/{name}/logs", applicationLogHandler())
    return http.ListenAndServe(":7000", mux)
}

func applicationLogHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !isSSERequest(r) {
            createComponentHandleFunc(func(r *http.Request) (templ.Component, error) {
                return templates.ApplicationLogs(r.PathValue("name")), nil
            })(w, r)
            return
        }
        w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")
        w.Header().Set("Access-Control-Allow-Origin", "*")

        name := r.PathValue("name")
        outputChan := make(chan string)
        defer func() {
            application.StopWatching(name)
            close(outputChan)
        }()
        err := application.WatchLogs(name, outputChan)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(err.Error()))
            return
        }

        flusher, _ := w.(http.Flusher)
        for {
            outputStr := <- outputChan
            formattedOutput := convertAnsiToHtml(outputStr)
            _, err := fmt.Fprintf(w, "data: <div>%s</div>\n\n", formattedOutput)
            if err != nil {
                break
            }
            flusher.Flush()
        }
    }
}

type CreateComponentFunc func(r *http.Request) (templ.Component, error)

func createComponentHandleFunc(createComponent CreateComponentFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        component, err := createComponent(r)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(err.Error()))
        }
        if component == nil {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        err = component.Render(context.Background(), w)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(err.Error()))
        }
    }
}

var ansiToHtml = map[string]string{
	"0":  "reset",
	"":  "reset",
	"1":  "bold",
	"3":  "italic",
	"4":  "underline",
	"7":  "inverse",
	"30": "color:white",// Is black on the console
	"31": "color:red",
	"32": "color:green",
	"33": "color:yellow",
	"34": "color:blue",
	"35": "color:magenta",
	"36": "color:cyan",
	"37": "color:white",
	"40": "background-color:black",
	"41": "background-color:red",
	"42": "background-color:green",
	"43": "background-color:yellow",
	"44": "background-color:blue",
	"45": "background-color:magenta",
	"46": "background-color:cyan",
	"47": "background-color:white",
}

var htmlTags = map[string]string{
	"bold":      `<span style="font-weight:bold;">`,
	"italic":    `<span style="font-style:italic;">`,
	"underline": `<span style="text-decoration:underline;">`,
	"inverse":   `<span style="filter:invert(100%);">`,
    "reset":     `</span>`,
}

var ansiEscape = regexp.MustCompile(`\x1B\[([0-9;]*)([A-Za-z])`)

func convertAnsiToHtml(text string) string {
	// Replace ANSI escape sequences with HTML tags
	text = ansiEscape.ReplaceAllStringFunc(text, func(match string) string {
		matches := ansiEscape.FindStringSubmatch(match)
        // FIXME If there is more than one code, it will not work
        code := matches[1]


        codeOutput := ansiToHtml[code]

        if tag, ok := htmlTags[codeOutput]; ok {
            return tag;
        } else {
            return fmt.Sprintf(`<span style="%s;">`, codeOutput);
        }
	})

	return text
}
