package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/assets/amazon"
	"github.com/kevinwylder/sbvision/assets/filesystem"
)

// GetLocalAssets looks at environment variables to determine where assets are stored
// it returns the key value store and an optional string if a temporary directory was used
// if ASSET_DIR exists, then local storage will be in the given dir, and no cleanup string will be used
// if S3_BUCKET exists, then the key value store will be in s3, backed by the file store
func GetLocalAssets() (sbvision.KeyValueStore, string) {
	var err error
	tmpAssets := ""
	assetDir, exists := os.LookupEnv("ASSET_DIR")
	if !exists {
		assetDir, err = ioutil.TempDir("", "")
		if err != nil {
			log.Fatal("Could not create tmp dir for image storage", err)
		}
		tmpAssets = assetDir
	}
	cache, err := filesystem.NewAssetDirectory(assetDir)
	if err != nil {
		log.Fatal("Could not use the asset dir - ", err)
	}
	if bucket, exists := os.LookupEnv("S3_BUCKET"); exists {
		assets, err := amazon.NewS3BucketManager(bucket, cache)
		if err != nil {
			log.Fatal("Error using S3_BUCKET ", err)
		}
		return assets, tmpAssets
	}
	return cache, tmpAssets
}
