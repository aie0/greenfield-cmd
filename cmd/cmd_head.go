package main

import (
	"context"
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

// cmdHeadObj return the command to finish uploading payload of the object
func cmdHeadObj() *cli.Command {
	return &cli.Command{
		Name:      "head-obj",
		Action:    headObject,
		Usage:     "query object info",
		ArgsUsage: "OBJECT-URL",
		Description: `
send headObject txn to chain and fetch objectInfo on greenfield chain
Examples:
$ gnfd-cmd head-bucket gnfd://bucket-name/object-name`,
	}
}

// cmdHeadBucket return the command to finish uploading payload of the object
func cmdHeadBucket() *cli.Command {
	return &cli.Command{
		Name:      "head-bucket",
		Action:    headBucket,
		Usage:     "query bucket info",
		ArgsUsage: "BUCKET-URL",
		Description: `
send headBucket txn to chain and fetch bucketInfo on greenfield chain
Examples:
$ gnfd-cmd head-bucket gnfd://bucket-name`,
	}
}

func headObject(ctx *cli.Context) error {
	urlInfo := ctx.Args().Get(0)
	bucketName, objectName, err := getObjAndBucketNames(urlInfo)
	if err != nil {
		return nil
	}

	client, err := NewClient(ctx)
	if err != nil {
		log.Println("failed to create client", err.Error())
		return err
	}

	c, cancelCreateBucket := context.WithCancel(globalContext)
	defer cancelCreateBucket()

	objectInfo, err := client.HeadObject(c, bucketName, objectName)
	if err != nil {
		fmt.Println("headObject fail:", err.Error())
		return err
	}
	fmt.Println("object id:", objectInfo.Id.Uint64())
	fmt.Println("object status", objectInfo.ObjectStatus)
	fmt.Println("object size:", objectInfo.Size)

	return nil
}

// headBucket send the create bucket request to storage provider
func headBucket(ctx *cli.Context) error {
	bucketName, err := getBucketName(ctx)
	if err != nil {
		return err
	}

	client, err := NewClient(ctx)
	if err != nil {
		log.Println("failed to create client", err.Error())
		return err
	}

	c, cancelCreateBucket := context.WithCancel(globalContext)
	defer cancelCreateBucket()

	bucketInfo, err := client.HeadBucket(c, bucketName)
	if err != nil {
		fmt.Println("headBucket fail:", err.Error())
		return err
	}
	fmt.Println("bucket id:", bucketInfo.Id.Uint64())
	fmt.Println("bucket owner:", bucketInfo.Owner)

	return nil
}
