#!/bin/sh

cd `dirname $0`

dir=vendor/github.com/hashicorp/go-plugin/examples/bidirectional

rm -f $dir/multer
rm -f $dir/protob

ln -s $PWD/multer $dir/multer
ln -s $PWD/protob $dir/protob

