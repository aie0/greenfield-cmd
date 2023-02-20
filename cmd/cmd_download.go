package main

import (
	"context"
	"fmt"
	"log"

	spClient "github.com/bnb-chain/greenfield-go-sdk/client/sp"
	"github.com/urfave/cli/v2"
)

// cmdGetObj return the command to finish downloading object payload
func cmdGetObj() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Action:    getObject,
		Usage:     "Download object",
		ArgsUsage: "[filePath] OBJECT-URL",
		Description: `
Download a specific object from storage provider

Examples:
# download a file
$ gnfd get --start 1  --end 4000 gnfd://bucketname/file.txt file.txt `,
		Flags: []cli.Flag{
			&cli.Uint64Flag{
				Name:  "start",
				Value: 0,
				Usage: "start offset of download range",
			},
			&cli.Uint64Flag{
				Name:  "end",
				Value: 0,
				Usage: "end offset of download range",
			},
		},
	}
}

// getObject download the object payload from sp
func getObject(ctx *cli.Context) error {
	if ctx.NArg() != 2 {
		return fmt.Errorf("the args number should be two")
	}

	urlInfo := ctx.Args().Get(0)
	bucketName, objectName := ParseBucketAndObject(urlInfo)

	s3Client, err := NewClient(ctx)
	if err != nil {
		log.Println("create client fail", err.Error())
		return err
	}

	c, cancelCreateBucket := context.WithCancel(globalContext)
	defer cancelCreateBucket()

	filePath := ctx.Args().Get(1)
	log.Printf("download object %s into file:%s \n", objectName, filePath)

	startIndex := ctx.Uint64("start")
	endIndex := ctx.Uint64("end")
	option := spClient.DownloadOption{}
	option.SetRange(int64(startIndex), int64(endIndex))

	err = s3Client.FGetObject(c, bucketName, objectName, filePath, option, spClient.NewAuthInfo(false, ""))
	if err != nil {
		return err
	}
	log.Println("downaload object finish")
	return nil
}
