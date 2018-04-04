package lib

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	tm "github.com/buger/goterm"
)

func MarkdownReport(s *State) string {
	var report bytes.Buffer
	currTime := strings.Replace(time.Now().Format(time.RFC3339), ":", "_", -1)
	reportFile := path.Join(s.ReportDirectory, fmt.Sprintf("%v_Report.md", currTime))
	file, err := os.Create(reportFile)
	if err != nil {
		panic(err)
	}
	// Header
	report.WriteString(fmt.Sprintf("# Gograbber report - %v (%v)\n", s.ProjectName, currTime))
	for _, URLComponent := range s.URLComponents {
		var path string
		for a := range URLComponent.Paths.Set {
			path = a
		}
		url := fmt.Sprintf("%v://%v:%v/%v\n", URLComponent.Protocol, URLComponent.HostAddr, URLComponent.Port, path)
		report.WriteString(fmt.Sprintf("## %v\n", url))
		report.WriteString(fmt.Sprintf("![%v](../../%v)\n", URLComponent.ScreenshotFilename, URLComponent.ScreenshotFilename))

	}
	file.WriteString(report.String())
	return reportFile
}

func TextOutput(s *State) {
	tm.Clear()
	box := tm.NewBox(100|tm.PCT, 5, 1)
	fmt.Fprint(box, "Some box content")
	time.Sleep(5 * time.Second)
}

func JSONify(s *State) {

}

func SanitiseFilename(UnsanitisedFilename string) (filename string) {
	r := regexp.MustCompile("[0-9a-zA-Z-._]")
	return r.ReplaceAllString(UnsanitisedFilename, "[0-9a-zA-Z-._]")
}
