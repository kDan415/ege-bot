package ege

import (
	"fmt"
	"time"
)

type Credentials struct {
	Time   time.Time `json:"-"`
	Result struct {
		Exams []struct {
			ExamId     int    `json:"ExamId"`
			Subject    string `json:"Subject"`
			TestMark   int    `json:"TestMark"`
			Mark5      int    `json:"Mark5"`
			MinMark    int    `json:"MinMark"`
			Status     int    `json:"Status"`
			IsHidden   bool   `json:"IsHidden"`
			StatusName string `json:"StatusName"`
			HasResult  bool   `json:"HasResult"`
		} `json:"Exams"`
	} `json:"Result"`
}

func (c *Credentials) String() string {
	if c == nil {
		return ""
	}
	result := "Date: " + c.Time.Format(time.DateTime) + "\n\n"
	for _, exam := range c.Result.Exams {
		result += fmt.Sprintf("Exam: %s\nHasResult?: %t\nIsHidden?: %t\nStatus: %d %s\nTestMark: %d\nMark5: %d\n\n",
			exam.Subject,
			exam.HasResult,
			exam.IsHidden,
			exam.Status,
			exam.StatusName,
			exam.TestMark,
			exam.Mark5)
	}
	return result
}

// returns true if 2 result sets are exactly the same
func compare(first *Credentials, second *Credentials) bool {
	if first == nil || second == nil {
		return false
	}
	if len(first.Result.Exams) != len(second.Result.Exams) {
		return false
	}
	for i, exam1 := range first.Result.Exams {
		exam2 := second.Result.Exams[i]
		if exam1 != exam2 {
			return false
		}
	}
	return true
}
