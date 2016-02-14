#!/bin/sh
UCL_VERSION=0.7.3
set -ex
mkdir /tmp/libucl
cd /tmp/libucl
wget https://github.com/vstakhov/libucl/archive/$UCL_VERSION.tar.gz
tar xzf $UCL_VERSION.tar.gz
cd libucl-$UCL_VERSION && ./autogen.sh && ./configure --prefix=/usr --enable-urls --enable-signatures && make && sudo make install
rm -fr /tmp/libucl
