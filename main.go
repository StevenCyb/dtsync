package main

import (
	"dtsync/pkg/args"
	"dtsync/pkg/fs"
	"dtsync/pkg/screen"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	Run(args.Parse(os.Args))
}

func Run(arguments args.Arguments) {
	var (
		srcErrChan, dstErrChan chan error
		err                    error
	)

	view := screen.NewView(os.Stdout)
	if err = view.Start(); err != nil {
		log.Println(err.Error())
	}

	defer view.Stop()

	operation := fs.NewOperation()
	scanner := fs.NewShadowScan()

	defer scanner.Stop()

	srcErrChan = scanner.Start(arguments.SrcRootPath, arguments.DstRootPath,
		func(srcPath, dstPath string) error {
			if !operation.Exists(dstPath) {
				view.AddStatus(screen.Status{SrcTotalFiles: 1, Copied: 1})

				return operation.Copy(srcPath, dstPath)
			} else if arguments.ReplaceNotMatchingFiles && !operation.Equal(srcPath, dstPath) {
				view.AddStatus(screen.Status{SrcTotalFiles: 1, Replaced: 1})

				return operation.Copy(srcPath, dstPath)
			}

			view.AddStatus(screen.Status{SrcTotalFiles: 1, Skipped: 1})

			return nil
		},
		func(srcPath, dstPath string) error {
			if !operation.Exists(dstPath) {
				view.AddStatus(screen.Status{SrcTotalDirectories: 1, Copied: 1})

				return operation.Copy(srcPath, dstPath)
			}

			view.AddStatus(screen.Status{SrcTotalDirectories: 1, Skipped: 1})

			return nil
		},
	)

	if arguments.RemoveDstLeftover {
		dstErrChan = scanner.Start(arguments.DstRootPath, arguments.SrcRootPath,
			func(srcPath, dstPath string) error {
				if !operation.Exists(dstPath) {
					view.AddStatus(screen.Status{DstTotalFiles: 1, Removed: 1})

					return operation.Delete(srcPath)
				}

				return nil
			},
			func(srcPath, dstPath string) error {
				if !operation.Exists(dstPath) {
					view.AddStatus(screen.Status{DstTotalDirectories: 1, Removed: 1})

					return operation.Delete(srcPath)
				}

				return nil
			},
		)
	}

	// waiting for any for ending conditions
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	for i := 0; i < 2; i++ {
		select {
		case <-signalChan:
			goto end
		case err = <-srcErrChan:
			if err != nil && !errors.Is(err, fs.ErrScannerAtEnd) {
				goto end
			}

			if !arguments.RemoveDstLeftover {
				goto end
			}
		case err = <-dstErrChan:
			if err != nil && !errors.Is(err, fs.ErrScannerAtEnd) {
				goto end
			}
		}
	}

end:
	view.Render()

	if err != nil && !errors.Is(err, fs.ErrScannerAtEnd) {
		log.Println(err.Error())
	}
}
