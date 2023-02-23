package main

import (
	"fmt"

	"github.com/bnb-chain/greenfield/sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/urfave/cli/v2"
)

// cmdDelBucket delete a exist Bucket,the bucket must be empty
func cmdDelBucket() *cli.Command {
	return &cli.Command{
		Name:      "del-bucket",
		Action:    deleteBucket,
		Usage:     "deletea a  bucket",
		ArgsUsage: "BUCKET-URL",
		Description: `
Send a deleteBucket txn to greenfield chain

Examples:
# Del a exist bucket
$ gnfd del gnfd://bucketname`,
	}
}

// cmdDelObject delete a exist object in bucket
func cmdDelObject() *cli.Command {
	return &cli.Command{
		Name:      "del-obj",
		Action:    deleteObject,
		Usage:     "create a new bucket",
		ArgsUsage: "BUCKET-URL",
		Description: `
Send a deleteObject txn to greenfield chain

Examples:
# Del a exist object
$ gnfd del gnfd://bucketname/objectname`,
	}
}

// deleteBucket send the deleteBucket msg to greenfield
func deleteBucket(ctx *cli.Context) error {
	bucketName, err := getBucketName(ctx)
	if err != nil {
		return err
	}

	client, err := NewClient(ctx)
	if err != nil {
		return err
	}

	broadcastMode := tx.BroadcastMode_BROADCAST_MODE_BLOCK
	gnfdResp := client.DelBucket(client.SPClient.GetAccount(), bucketName, types.TxOption{Mode: &broadcastMode})
	if gnfdResp.Err != nil {
		return err
	}

	fmt.Println("delete bucket finish, txn hash:", gnfdResp.TxnHash)
	return nil
}

// deleteObject send the deleteBucket msg to greenfield
func deleteObject(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf("the args number should be one")
	}

	urlInfo := ctx.Args().Get(0)
	bucketName, objectName, err := getObjAndBucketNames(urlInfo)
	if err != nil {
		return nil
	}

	client, err := NewClient(ctx)
	if err != nil {
		return err
	}

	broadcastMode := tx.BroadcastMode_BROADCAST_MODE_BLOCK
	gnfdResp := client.DelObject(client.SPClient.GetAccount(), bucketName, objectName, types.TxOption{Mode: &broadcastMode})
	if gnfdResp.Err != nil {
		return err
	}

	fmt.Println("delete object finish, txn hash:", gnfdResp.TxnHash)
	return nil
}
