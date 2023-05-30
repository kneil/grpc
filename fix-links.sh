#!/bin/sh

cd /tmp/repo/vendor/github.com/hashicorp/go-plugin/examples/bidirectional/

rm -f multer
rm -f protob

ln -s /tmp/repo/multer
ln -s /tmp/repo/protob

