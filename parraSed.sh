#!/bin/sh

## Parrallel sed adaptation from : 
## https://stackoverflow.com/questions/1828236/how-to-make-this-sed-script-faster
## from https://stackoverflow.com/users/874188/tripleee

SEDSCRIPT=${1:-script.sed}
   INFILE=${2:-input.txt}
  OUTFILE=${3:-output.txt}


cpus=`nproc`

SPLITLIMIT=`wc -l $INFILE | awk -v cpus=$cpus '{printf("%i", $1 / cpus ) }'`

split -d -l $SPLITLIMIT $INFILE x_

for chunk in x_??
do
  sed -f $SEDSCRIPT $chunk > $chunk.out &
done

wait 

cat x_??.out >> $OUTFILE

rm -f x_??
rm -f x_??.out
