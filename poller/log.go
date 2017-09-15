package poller

import (
	"../db"
	"log"
	"time"
)

type LogEntry struct {
	TaskId     int       `json:"taskId"`
	StatusCode int       `json:"statusCode"`
	Time       time.Time `json:"time"`
}

func SaveLogEntry(task Task, statusCode int) {
	_, err := db.DB.Exec("REPLACE INTO logs SET task_id=?, time=NOW(), status_code=?", task.Id, statusCode)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadLogEntries(taskId int, limit int) []LogEntry {
	query := "SELECT task_id, time, status_code FROM logs WHERE task_id=? ORDER BY time"
	args := []interface{}{taskId}
	if limit > 0 {
		query = "SELECT * FROM (SELECT task_id, time, status_code FROM logs WHERE task_id=? ORDER BY time DESC LIMIT ?) t ORDER BY time"
		args = append(args, limit)
	}
	logEntries := make([]LogEntry, 0)
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var e LogEntry
		err := rows.Scan(&e.TaskId, &e.Time, &e.StatusCode)
		if err != nil {
			log.Fatal(err)
		}
		logEntries = append(logEntries, e)
	}
	return logEntries

}
