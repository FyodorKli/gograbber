package libgograbber

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/swarley7/phantomjs"
)

// Screenshots a url derived from a Host{} object
func ScreenshotAURL(wg *sync.WaitGroup, s *State, cnt int, host Host, results chan Host, threads chan struct{}) (err error) {
	// Ideally this function would not use phantomjs - I've looked at WebKit-go and that looks promising
	defer func() {
		<-threads
		wg.Done()
	}()
	page, err := s.PhantomProcesses[cnt%len(s.PhantomProcesses)].CreateWebPage()
	url := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, host.Path)

	if err != nil {
		Error.Printf("Unable to Create webpage: %v (%v)\n", url, err)
		return err
	}
	defer page.Close()

	page.SetSettings(phantomjs.WebPageSettings{ResourceTimeout: s.Timeout + (time.Second * 2)}) // Time out the page if it takes too long to load. Sometimes JS is fucky and takes wicked long to do nothing forever :(

	if strings.HasPrefix(host.Path, "/") {
		host.Path = host.Path[1:] // strip preceding '/' char
	}
	if s.Debug {
		Debug.Printf("Trying to screenshot URL: %v\n", url)
	}
	ApplyJitter(s.Jitter)
	if err := page.Open(url); err != nil {
		Error.Printf("Unable to open page: %v (%v)\n", url, err)
		return err
	}
	// Setup the viewport and render the results view.
	if err := page.SetViewportSize(s.ImgX, s.ImgY); err != nil {
		Error.Printf("Unable to set Viewport size: %v (%v)\n", url, err)
		// <-target
		return err
	}
	currTime := GetTimeString()
	var screenshotFilename string
	if s.ProjectName != "" {
		screenshotFilename = fmt.Sprintf("%v/%v_%v-%v_%v.%v", s.ScreenshotDirectory, strings.ToLower(SanitiseFilename(s.ProjectName)), SanitiseFilename(url), currTime, rand.Int63(), s.ScreenshotFileType)
	} else {
		screenshotFilename = fmt.Sprintf("%v/%v-%v_%v.%v", s.ScreenshotDirectory, SanitiseFilename(url), currTime, rand.Int63(), s.ScreenshotFileType)
	}
	if err := page.Render(screenshotFilename, s.ScreenshotFileType, s.ScreenshotQuality); err != nil {
		Error.Printf("Unable to save Screenshot: %v (%v)\n", url, err)
		return err
	}
	Good.Printf("Screenshot for [%v] saved to: [%v]\n", g.Sprintf("%s", url), g.Sprintf("%s", screenshotFilename))
	host.ScreenshotFilename = screenshotFilename
	results <- host
	return err
}
