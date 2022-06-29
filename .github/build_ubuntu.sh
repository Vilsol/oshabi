#!/usr/bin/env bash

set -ex

sudo apt install -y unzip wget build-essential cmake curl git libgtk2.0-dev pkg-config libavcodec-dev libavformat-dev libswscale-dev libtbb2 libtbb-dev libjpeg-dev libpng-dev libtiff-dev libdc1394-22-dev

mkdir -p /tmp/opencv
pushd /tmp/opencv
curl -Lo opencv.zip https://github.com/opencv/opencv/archive/4.6.0.zip
unzip -q -u opencv.zip
curl -Lo opencv_contrib.zip https://github.com/opencv/opencv_contrib/archive/4.6.0.zip
unzip -q -u opencv_contrib.zip
rm opencv.zip opencv_contrib.zip
popd

pushd /tmp/opencv/opencv-4.6.0
mkdir -p build
pushd build

cmake -D CMAKE_BUILD_TYPE=RELEASE -D CMAKE_INSTALL_PREFIX=/usr/local -D BUILD_SHARED_LIBS=ON -D OPENCV_EXTRA_MODULES_PATH=/tmp/opencv/opencv_contrib-4.6.0/modules -D BUILD_DOCS=OFF -D BUILD_EXAMPLES=OFF -D BUILD_TESTS=OFF -D BUILD_PERF_TESTS=OFF -D BUILD_opencv_java=NO -D BUILD_opencv_python=NO -D BUILD_opencv_python2=NO -D BUILD_opencv_python3=NO -D WITH_JASPER=OFF -D WITH_TBB=ON -DOPENCV_GENERATE_PKGCONFIG=ON ..
make -j $(nproc --all)
make preinstall

sudo make install
sudo ldconfig

popd
popd
