#!/bin/sh

. ./license_key

WGET="wget -N --no-if-modified-since"

mkdir tmp
$WGET -q "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=$license_key&suffix=tar.gz" -O  tmp/GeoLite2-Country.tar.gz
$WGET -q "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-ASN&license_key=$license_key&suffix=tar.gz" -O  tmp/GeoLite2-ASN.tar.gz

cd tmp
tar zxf GeoLite2-Country.tar.gz
tar zxf GeoLite2-ASN.tar.gz

ls *.gz | xargs -n 1  tar zxf
rm *.gz
find . -name \*.mmdb | xargs -I[] cp []  ../data/


COUNTRYFILE=`ls -d GeoLite2-Country_*`
ASNFILE=`ls -d GeoLite2-ASN*`

echo COUNTRYFILE : $COUNTRYFILE
echo ASNFILE : $ASNFILE

cd ..

cat init.tpl |\
sed "s+%COUNTRYFILE%+$COUNTRYFILE+g" |\
sed "s+%ASNFILE%+$ASNFILE+g" > init.go

rm -fr tmp
