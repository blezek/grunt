#!/bin/sh

# This script takes several arguments:
# dso.zip  -  a zip file containing DSO (DICOM segmentation objects) files
# dicom.zip - a zip file contating a DICOM series
# output.zip - the name of an output zip file

dso=$1
dicom=$2
output=$3

mkdir -p work/dso work/dicom work/output
unzip -d work/dso $dso
unzip -d dicom/dicom $dicom
/riipl/runFeaturePipeline work/dso work/dicom work/output

wd=$(pwd)
cd work/output
zip -r $wd/$output *

rm -rf work/
