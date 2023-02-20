package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	spClient "github.com/bnb-chain/greenfield-go-sdk/client/sp"
	"github.com/urfave/cli/v2"
)

// cmdCalHash create a new Bucket
func cmdCalHash() *cli.Command {
	return &cli.Command{
		Name:      "get-hash",
		Action:    computeHashRoot,
		Usage:     "compute hash roots of object ",
		ArgsUsage: "filePath",
		Description: `

Examples:
# Compute file path
$ gnfd get-hash --segSize 16  --shards 6 /home/test.text `,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "segSize",
				Value:    16,
				Usage:    "the segment size (MB)",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "dataShards",
				Value:    4,
				Usage:    "the ec encode data shard number",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "parityShards",
				Value:    2,
				Usage:    "the ec encode parity shard number",
				Required: true,
			},
		},
	}
}

func cmdChallenge() *cli.Command {
	return &cli.Command{
		Name:      "challenge",
		Action:    getChallengeInfo,
		Usage:     "make a challenge to sp",
		ArgsUsage: "",
		Description: `

Examples:
# Make challenge
$ gnfd challenge --objectId "test" --pieceIndex 2  --spIndex -1`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "objectId",
				Value:    "",
				Usage:    "the objectId to be challenge",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "pieceIndex",
				Value:    0,
				Usage:    "show which piece to be challenge",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "spIndex",
				Value:    -1,
				Usage:    "show which sp of the s",
				Required: true,
			},
		},
	}
}

func computeHashRoot(ctx *cli.Context) error {
	// read the local file payload to be uploaded
	filePath := ctx.Args().Get(0)

	exists, objectSize, err := pathExists(filePath)
	if !exists {
		return errors.New("file not exists")
	} else if objectSize > int64(500*1024*1024) {
		return errors.New("upload file can not be larger than 500M")
	}

	// Open the referenced file.
	fReader, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fReader.Close()

	segmentSize := ctx.Int("segSize")
	if segmentSize <= 0 {
		return errors.New("segment size should be more than 0 ")
	}

	dataShards := ctx.Int("dataShards")
	if dataShards <= 0 {
		return errors.New("data shards number should be more than 0 ")
	}

	parityShards := ctx.Int("parityShards")
	if parityShards <= 0 {
		return errors.New("parity shards number should be more than 0 ")
	}

	s3Client, err := NewClient(ctx)
	if err != nil {
		return err
	}

	priHash, secondHash, _, err := s3Client.GetPieceHashRoots(fReader, int64(segmentSize*1024*1024), dataShards, parityShards)
	if err != nil {
		return err
	}

	fmt.Printf("get primary sp hash root: \n%s\n", priHash)
	fmt.Println("get secondary sp hash list:")
	for _, hash := range secondHash {
		fmt.Println(hash)
	}

	return nil
}

func getChallengeInfo(ctx *cli.Context) error {
	objectId := ctx.String("objectId")
	if objectId == "" {
		return errors.New("object id empty ")
	}

	pieceIndex := ctx.Int("pieceIndex")
	if pieceIndex <= 0 {
		return errors.New("pieceIndex should not be less than 0 ")
	}

	spIndex := ctx.Int("spIndex")
	if spIndex < -1 {
		return errors.New("redundancyIndex should not be less than -1")
	}

	s3Client, err := NewClient(ctx)
	if err != nil {
		return err
	}

	filePath := ctx.Args().Get(0)
	log.Printf("download challenge payload into file:%s \n", filePath)

	st, err := os.Stat(filePath)
	if err == nil {
		// If the destination exists and is a directory.
		if st.IsDir() {
			return errors.New("fileName is a directory.")
		}
	}

	// If file exist, open it in append mode
	fd, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}

	info := spClient.ChallengeInfo{
		ObjectId:        objectId,
		PieceIndex:      pieceIndex,
		RedundancyIndex: spIndex,
	}

	c, cancelCreateBucket := context.WithCancel(globalContext)
	defer cancelCreateBucket()

	res, err := s3Client.ChallengeSP(c, info, spClient.NewAuthInfo(false, ""))
	if err != nil {
		return err
	}

	if res.PiecesHash != nil {
		fmt.Println("get hash result", res.PiecesHash)
	} else {
		return errors.New("fail to fetch piece hashes")
	}

	if res.PieceData != nil {
		defer res.PieceData.Close()
		_, err = io.Copy(fd, res.PieceData)
		fd.Close()
		if err != nil {
			return err
		}

		fmt.Println("get challenge piece data:", res.PieceData)
	} else {
		return errors.New("fail to fetch challenge data")
	}

	return nil
}
