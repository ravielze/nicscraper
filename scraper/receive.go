package scraper

import (
	"sort"

	"github.com/mkamadeus/nicscraper/models"
	"github.com/sirupsen/logrus"
)

func (s Scraper) Receive() ([]models.Student, []string) {
	size := len(s.Args.Prefixes) * len(s.Args.Years) * s.Args.Limit
	result := make([]models.Student, 0)
	failed := make([]string, 0)

	logrus.Debug(size)

	for i := 0; i < size; i++ {
		select {
		case student := <-s.Students:
			logrus.Debugf("Received student: %s", student)
			result = append(result, student)
		case err := <-s.Failed:
			logrus.Debugf("Received failed: %s", err)
			failed = append(failed, err)
		}
	}

	// Sort output
	sort.Slice(result[:], func(i, j int) bool {
		if result[i].MajorID == "" && result[j].MajorID == "" {
			return result[i].FacultyID < result[j].FacultyID
		}
		return result[i].MajorID < result[j].MajorID
	})
	sort.Slice(failed[:], func(i, j int) bool {
		return failed[i] < failed[j]
	})

	return result, failed
}
