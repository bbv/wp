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
	Url   string `yaml:"url" json:"url"`
	Delay int    `yaml:"delay" json:"-"`
	Id    int    `json:"id"`
}

func (t Task) Poll() {
	for {
		resp, _ := http.Get(t.Url)
		log.Println(t.Url, ": ", resp.StatusCode)
		SaveLogEntry(t, resp.StatusCode)
		time.Sleep(time.Duration(t.Delay) * time.Second)
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
	res, err := db.DB.Exec("INSERT INTO tasks SET url=?", t.Url)
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

func LoadTasks() []Task {
	tasks := make([]Task, 0)
	rows, err := db.DB.Query("SELECT id, url FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.Url)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, t)
	}
	return tasks
}
