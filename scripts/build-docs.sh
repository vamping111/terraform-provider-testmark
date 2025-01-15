#!/bin/bash

# Path of the directory with the documentation
SOURCE_DIR="website/site/"

if [[ -z $S3_DOCS_BUCKET_NAME ]]; then
  echo "Define S3_DOCS_BUCKET_NAME environment variable."
  exit 1
fi

# Copying of configuration for the documentatiion to the website folder
cp -r docs/c2/mkdocs/assets docs/c2/mkdocs/images website/docs/
cp docs/c2/mkdocs/mkdocs.yml website/

# Generation of the documentation
cd website
echo "Generation of the documentation"
mkdocs build --clean

# Funtion to upload the documentation files to the bucket
upload_other_files () {
    s3cmd sync "$SOURCE_DIR" "$S3_DOCS_BUCKET_NAME" --acl-public
}

# Funtion to upload .css files to the bucket with specified Content-Type
upload_css_files () {
    find "$SOURCE_DIR" -type f -name "*.css" | while read -r file; do
    s3cmd modify "s3://$S3_DOCS_BUCKET_NAME${file#$SOURCE_DIR}" --add-header='Content-Type:text/css'
done
}

# Funtion to upload .js files to the bucket with specified Content-Type
upload_js_files () {
    find "$SOURCE_DIR" -type f -name "*.js" | while read -r file; do
    s3cmd modify "s3://$S3_DOCS_BUCKET_NAME${file#$SOURCE_DIR}" --add-header='Content-Type:application/javascript'
done
}

# Uploading of the documentation to the S3 bucket
cd ..
echo "Uploading of the documentation to the S3 bucket"
upload_other_files
upload_css_files
upload_js_files

# Removing of the temp files
echo "Removing of the temp files"
rm -fr website/docs/images website/docs/assets website/mkdocs.yml website/site
echo Complete!
