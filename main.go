package main

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"path"
	"syscall"
)

func main() {
	app := &cli.App{
		Name:    "dist",
		Usage:   "tool to distribute data with filplus",
		Version: UserVersion(),
		Commands: []*cli.Command{
			userManager,
			dataSetManager,
			pieceManager,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "repo",
				EnvVars: []string{"DIST_PATH"},
				Value:   "~/.dist",
			},
		},
		Before: func(ctx *cli.Context) error {
			// 打开或者创建 lockFile 文件
			file, err := os.OpenFile(os.TempDir()+"/myapp.lock", os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				panic(err)
			}
			//defer file.Close()

			// 尝试获取文件锁
			err = unix.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
			if err != nil {
				if err == syscall.EWOULDBLOCK {
					fmt.Println("Another instance is already running...")
				} else {
					panic(err)
				}
				os.Exit(1)
			}

			homeDir, err := homedir.Expand(ctx.String("repo"))

			if err != nil {
				return err
			}
			_, err = os.Stat(homeDir)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					err = os.MkdirAll(homeDir, 0755)
					if err != nil {
						return err
					}
				}
			}
			orgsJson = path.Join(homeDir, "users.json")
			dataSetsJson = path.Join(homeDir, "datasets.json")

			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}
}
