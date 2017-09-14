package poller

import (
	"../db"
	"database/sql"
	"log"
	"net/http"
	"time"
)

type Poller struct {
	Tasks []Task
}

type Task struct {
	Url   string `yaml:url`
	Delay int    `yaml:delay`
	Id    int
}

type Response struct {
	TaskId     int
	StatusCode int
	Time       time.Time
}

func (t Task) Poll() {
	for {
		resp, _ := http.Get(t.Url)
		log.Println(t.Url, ": ", resp.StatusCode)
		//time.Sleep(time.Duration(t.Delay) * time.Second)
		t.SaveResponse(resp.StatusCode)
		time.Sleep(2 * time.Second)
	}
}

func NewPoller(tasks []Task) Poller {
	for _, task := range tasks {
		task.Load()
		go task.Poll()
	}
	return Poller{Tasks: tasks}
}

func (t *Task) Load() error {
	var id int
	err := db.DB.QueryRow("SELECT id FROM tasks WHERE url=?", t.Url).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		t.Save()
	case err != nil:
		log.Fatal(err)
	default:
		t.Id = id
	}
	return nil
}

func (t *Task) Save() error {
	res, err := db.DB.Exec("INSERT INtO tasks SET url=?", t.Url)
	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	t.Id = int(id)
	return nil
}

func (t Task) SaveResponse(statusCode int) error {
	_, err := db.DB.Exec("REPLACE INTO logs SET task_id=?, time=NOW(), status_code=?", t.Id, statusCode)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
