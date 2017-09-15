package main

import (
	"./db"
	"./poller"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	appConfig, err := poller.ReadAppConfig("app.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Init(appConfig.Db)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/tasks/", tasks)
	http.ListenAndServe(fmt.Sprintf(":%d", appConfig.Port), nil)
}

func tasks(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) == 2 {
		tasks := poller.LoadTasks()
		json, err := json.Marshal(tasks)
		if err != nil {
			log.Fatal(err)
		}
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprint(res, string(json))
		return
	}
	if len(parts) > 2 {
		taskId, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Println(err)
		}
		limit := 0
		if len(parts) > 3 {
			limit, err = strconv.Atoi(parts[3])
			if err != nil {
				log.Println(err)
			}
		}
		logEntries := poller.LoadLogEntries(taskId, limit)
		json, err := json.Marshal(logEntries)
		fmt.Fprint(res, string(json))
		return
	}
}
