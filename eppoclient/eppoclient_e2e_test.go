package eppoclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const TEST_DATA_DIR = "test-data/assignment"
const BUCKET_NAME = "sdk-test-data"

func Test_e2e(t *testing.T) {
	downloadTestData()
}

func downloadTestData() {
	if _, err := os.Stat(TEST_DATA_DIR); os.IsNotExist(err) {
		if err := os.MkdirAll(TEST_DATA_DIR, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	} else {
		return
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		fmt.Println(err)
	}

	query := &storage.Query{Prefix: "assignment/test-case"}
	bkt := client.Bucket(BUCKET_NAME)
	it := bkt.Objects(ctx, query)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(attrs)
		obj := bkt.Object(attrs.Name)
		rdr, err := obj.NewReader(ctx)

		if err != nil {
			log.Fatal(err)
		}
		defer rdr.Close()

		out, err := os.Create("test-data/" + obj.ObjectName())
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		io.Copy(out, rdr)
	}
}
