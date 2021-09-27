package exam

import "time"

type Variant struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Answer struct {
	Id   int    `json:"id"`
	Text string `json:"string"`
}

type Task struct {
	Id        int      `json:"id"`
	VariantID int      `json:"variant_id"`
	Question  string   `json:"question"`
	Correct   int      `json:"correct"`
	Answers   []Answer `json:"answers"`
}

type Test struct {
	Id        int       `json:"id"`
	UserID    int       `json:"user_id"`
	VariantID int       `json:"variant_id"`
	StartAt   time.Time `json:"start_at"`
}

type UserAnswer struct {
	Id     int `json:"id"`
	TestID int `json:"test_id"`
	TaskID int `json:"task_id"`
	Answer int `json:"answer"`
}

type Result struct {
	Id      int `json:"id"`
	TestID  int `json:"test_id"`
	Percent int `json:"percent"`
}
