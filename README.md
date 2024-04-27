# FTP client for golang

## Upload

```golang
package main

import (
	"fmt"
	"github.com/civet148/log"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"time"
)

const (
	Version     = "v0.1.0"
	ProgramName = "ftp-client"
)

var (
	BuildTime = "2024-04-27"
	GitCommit = ""
)

const (
	CMD_NAME_UPLOAD = "upload"
)

const (
	CMD_FLAG_NAME_DEBUG    = "debug"
	CMD_FLAG_NAME_FTP_URL  = "ftp-url"
)

func main() {
    local := []*cli.Command{
        uploadCmd,
    }
    app := &cli.App{
        Name:     ProgramName,
        Version:  fmt.Sprintf("%s %s commit %s", Version, BuildTime, GitCommit),
        Flags:    []cli.Flag{},
        Commands: local,
        Action:   nil,
    }
    if err := app.Run(os.Args); err != nil {
    log.Errorf("exit in error %s", err)
    os.Exit(1)
    return
    }
}

var uploadCmd = &cli.Command{
	Name:      CMD_NAME_UPLOAD,
	Usage:     "upload file to FTP server",
	ArgsUsage: "<src> <dest>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    CMD_FLAG_NAME_DEBUG,
			Usage:   "open debug log mode",
			Aliases: []string{"d"},
		},
		&cli.StringFlag{
			Name:    CMD_FLAG_NAME_FTP_URL,
			Usage:   "ftp url, e.g ftp://user:password@127.0.0.1:21",
			Aliases: []string{"f"},
			Value:   "ftp://127.0.0.1:21",
		},
	},
	Action: func(cctx *cli.Context) (err error) {
		var strFtpUrl string
		var strSrcPath = cctx.Args().Get(0)
		var strDstPath = cctx.Args().Get(1)
		if strSrcPath == "" || strDstPath == "" {
			return fmt.Errorf("no src or dest file path")
		}
		if cctx.Bool(CMD_FLAG_NAME_DEBUG) {
			log.SetLevel("debug")
		}
		strFtpUrl = cctx.String(CMD_FLAG_NAME_FTP_URL)
		client := ftp.NewFtpClient(strFtpUrl)
		err = client.Upload(strSrcPath, strDstPath)
		if err != nil {
			return log.Errorf(err.Error())
		}
		return nil
	},
}
```

## Build and run

```shell

go build && ./ftp-client upload /local/file/path/test.txt /remote/file/path/test.txt


```