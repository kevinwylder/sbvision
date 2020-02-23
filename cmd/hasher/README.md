# hasher program

This program exists to migrate away from the centralized images table

## USAGE

1. Run the db migration `database/migrations/03-image-duplicates-1.sql` 

2. execute the build output from this directory (`hasher`) in the environment you will be migrating. It will look at recreate the environment from the following env variables used by `server`

* `DB_CREDS` to identify which database to scan
* `ASSET_DIR` to look for stored images
* `S3_BUCKET` if specified, uses the given S3 bucket to store in asset dir

3. Check the result. You should see the `frames.image_hash` and `videos.thumbnail_hash` columns filled out. You should also see that the image assets for frames and video are moved to be based off their ID

4. Run the second db migration `database/migrations/04-image-duplicates-2.sql`