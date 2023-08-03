package main

import (
	"strings"
	"testing"
)

func TestS3ParseWithKey(t *testing.T) {
	s3ParseHelper(t, "s3://testbucketname/prefix1/prefix2/testData.csv", &ParsedS3Url{
		Key:      "testData.csv",
		Prefixes: []string{"prefix1", "prefix2"},
		Bucket:   "testbucketname",
	})
}

func TestS3ParseWithoutKey(t *testing.T) {
	s3ParseHelper(t, "s3://testbucketname/prefix1/prefix2/", &ParsedS3Url{
		Key:      "",
		Prefixes: []string{"prefix1", "prefix2"},
		Bucket:   "testbucketname",
	})
}

func s3ParseHelper(t *testing.T, url string, exp *ParsedS3Url) {
	res, err := ParseS3Url(url)
	if err != nil {
		t.Error(err)
	}

	// bucket
	if exp.Bucket != res.Bucket {
		t.Errorf("Bucket %s does not match expected value %s", res.Bucket, exp.Bucket)
	}

	// prefixes
	if len(exp.Prefixes) == len(res.Prefixes) {
		for i := range exp.Prefixes {
			if exp.Prefixes[i] != res.Prefixes[i] {
				t.Errorf("Prefixes %s do not match expected value %s", strings.Join(res.Prefixes, ","), strings.Join(exp.Prefixes, ","))
			}
		}
	} else {
		t.Errorf("Prefixes %s do not match expected value %s", strings.Join(res.Prefixes, ","), strings.Join(exp.Prefixes, ","))
	}

	// key
	if exp.Key != res.Key {
		t.Errorf("Key %s does not match expected value %s", res.Key, exp.Key)
	}
}
