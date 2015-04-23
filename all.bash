#!/bin/bash

SNAPPY="snappy-1.1.1"
LEVELDB="leveldb-1.18"

rm -rf deps/include deps/libs deps/${SNAPPY} deps/${LEVELDB} || exit 1

cd deps/ && mkdir libs 

tar -zxf ${SNAPPY}.tar.gz && cd ${SNAPPY} && ./configure --disable-shared --with-pic && make || exit 1
SNAPPY_PATH=`pwd`

cp ${SNAPPY_PATH}/.libs/libsnappy.a ../libs

cd ../libs

export LIBRARY_PATH=`pwd`
export C_INCLUDE_PATH=${SNAPPY_PATH}
export CPLUS_INCLUDE_PATH=${SNAPPY_PATH}

cd ../

tar -zxf ${LEVELDB}.tar.gz && cd ${LEVELDB} && make || exit 1
cp libleveldb.a ../libs
mv include ../

cd ../../

make && make redisprotocol && make restore
