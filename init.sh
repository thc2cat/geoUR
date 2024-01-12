#!/bin/sh

# Source your free maxmind licence key first 
. ./license_key

WGET="wget -N --no-if-modified-since"

# Only download if modified since previous
# touch .last-modified 

curl --head "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-ASN&license_key=$license_key&suffix=tar.gz" 2>/dev/null |
grep last-modified > tmpmod

curl --head "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=$license_key&suffix=tar.gz" 2>/dev/null |
grep last-modified >> tmpmod

if cmp --silent -- "tmpmod" ".last-modified"; then
 rm tmpmod
 echo "No fresh updates !"
 exit
fi

mkdir tmp data 2>/dev/null

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

mv tmpmod .last-modified
rm -fr tmp
