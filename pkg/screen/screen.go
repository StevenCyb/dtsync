package screen

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/fatih/color"
)

// ErrNoOutputBuffer is returned when the output buffer is not set.
var ErrNoOutputBuffer = errors.New("no output buffer")

// Status holds the current progress status.
type Status struct {
	SrcTotalFiles       int
	SrcTotalDirectories int
	DstTotalFiles       int
	DstTotalDirectories int
	Copied              int
	Replaced            int
	Removed             int
	Skipped             int
}

// View provides a CLI view that shows a fixed text with the given number sets.
type View struct {
	stdout         io.Writer
	status         Status
	numberSetsLock sync.Mutex
	ticker         *time.Ticker
	stopChan       chan struct{}
	startTime      time.Time
	firstPrint     bool
}

func NewView(stdout io.Writer) View {
	return View{
		stdout:     stdout,
		firstPrint: true,
		startTime:  time.Now(),
		stopChan:   make(chan struct{}),
	}
}

// AddStatus sets number sets to display.
func (v *View) AddStatus(status Status) {
	v.numberSetsLock.Lock()
	defer v.numberSetsLock.Unlock()

	v.status.SrcTotalFiles += status.SrcTotalFiles
	v.status.SrcTotalDirectories += status.SrcTotalDirectories
	v.status.DstTotalFiles += status.DstTotalFiles
	v.status.DstTotalDirectories += status.DstTotalDirectories
	v.status.Copied += status.Copied
	v.status.Removed += status.Removed
	v.status.Replaced += status.Replaced
	v.status.Skipped += status.Skipped
}

// Render renders the view.
func (v *View) Render() {
	v.numberSetsLock.Lock()
	defer v.numberSetsLock.Unlock()

	if v.firstPrint {
		v.firstPrint = false
	} else {
		fmt.Print("\033[11F")
	}

	fmt.Printf("Elapsed       : %s\n\n", color.BlueString("%v", time.Since(v.startTime)))

	fmt.Printf("TotalSrcFiles : %s\n", color.CyanString("%d", v.status.SrcTotalFiles))
	fmt.Printf("TotalSrcDirs  : %s\n", color.CyanString("%d", v.status.SrcTotalDirectories))
	fmt.Printf("TotalDstFiles : %s\n", color.CyanString("%d", v.status.DstTotalFiles))
	fmt.Printf("TotalDstDirs  : %s\n\n", color.CyanString("%d", v.status.DstTotalDirectories))

	fmt.Printf("Copied        : %s\n", color.HiCyanString("%d", v.status.Copied))
	fmt.Printf("Replaced      : %s\n", color.HiCyanString("%d", v.status.Replaced))
	fmt.Printf("Removed       : %s\n", color.HiCyanString("%d", v.status.Removed))
	fmt.Printf("Skipped       : %s\n", color.HiCyanString("%d", v.status.Skipped))
}

// Start the view rendering.
func (v *View) Start() error {
	if v.stdout == nil {
		return ErrNoOutputBuffer
	}

	v.stopChan = make(chan struct{})

	go func() {
		v.ticker = time.NewTicker(time.Second)

		v.Render()

		for {
			select {
			case <-v.ticker.C:
				v.Render()
			case <-v.stopChan:
				v.ticker.Stop()
				v.Render()

				return
			}
		}
	}()

	return nil
}

// Stop the view rendering.
func (v *View) Stop() {
	if v.stopChan != nil {
		close(v.stopChan)
	}
}
