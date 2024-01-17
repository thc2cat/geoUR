# geoUR

geoUR Meaning geoip based who are you ?

geoUR use maxmind geoip databases, embed Country and ASN data in a portable binary, and reads inputs of ipv4 addresses on stdin, then output Country and ASN numbers in differents format ( SED, JSON )

A free Maxmind license_key file must be present in current directory, containing 'licence_key=xxxxx' in order to download Maxmind mmdb databases via init.sh script, wich will be embeded in cli binary.

ip to geoip info is done via Maxmind geoip available at [oschwald/geoip2-golang lib]("https://github.com/oschwald/geoip2-golang").

## Build (need golang installed)

```shell
> sh init.sh
> go mod tidy
> go build
```

On \*nix, init.sh should download maxmind dbs in data/ dir and produce init.go , wich mainly contains names of the assets.
Your can also manually download these db from maxmind and adapt init.tpl manually.

## Usage

```shell
> go run . -h
geoUR (v1.4) build with MaxMind free edition db :
        - GeoLite2-Country_20240116
        - GeoLite2-ASN_20240116

Usage: geoUR ip or stdin input
use BULKFORMAT env for JSON, SED format

> geoUR 217.69.139.150
217.69.139.150 (Russia / AS47764 LLC VK)

ip addresses must be extracted ( checkout my extractip repo ) , as exemple :

```shell
> rg "pattern" logs/* | extractip | geoUR | rg "Russia|Ivory"
...
> geoUR.exe < testdata.ips
...

> head -2  testuser | extractip | geoUR
93.2.135.70 (France / AS15557 Societe Francaise Du Radiotelephone - SFR SA)
176.134.198.13 (France / AS5410 Bouygues Telecom SA)

> setenv BULKFORMAT JSON
> head -2  testuser | extractip | geoUR
{ "ip":"93.2.135.70", "geoip": { "country":"France","AS":"AS15557 Societe Francaise Du Radiotelephone - SFR SA"} }
{ "ip":"176.134.198.13", "geoip": { "country":"France","AS":"AS5410 Bouygues Telecom SA"} }
```

 geoUR remove caracters such as quotes, \& , and \\. in AS to avoid issues when slurping with `jq -s '.'`.

```shell
> extractip < userip |geoUR|jq -s '.' >  usergeo.s
```

## JSON merging, filtering, searching  with jq examples

Merging two JSON files using "ip" key with jq is a way faster than sed.

```shell
# Combining files
> jq -s '[ .[0] + .[1] | group_by(.ip)[] | select(length > 1) | add ]' userip.s usergeo.s  > combined

# filtering users with more than 3 different geoip.<FIELD> entries
> jq -s '.[]|reduce .[] as $i ({}; (.[$i.user]+=[$i.geoip.<FIELD>]))|to_entries[]|{user:.key,val:[.value|unique]}|select((.val[]|length)>3)|.user' combined

# Quels sont les utilisateurs de Nord VPN : ( "AS136787 TEFINCOM S.A." )
> jq --arg p "AS136787 TEFINCOM S.A." '.[]|select(.geoip.AS==$p) |{user:.user}' combined |jq -s '.|group_by(.user)|add|unique|.[].user'

# show user without ips
> jq --arg p "user.name@domain.com" '.[]|select(.user==$p) |{user:.user, geoip: .geoip }' combined |jq -s '.|group_by(.user)|add|unique'

```

## SED output

With sed, look at parraSED.sh script for a speed up parrallel sed
but be aware of memory involved with concurrents sed.

```shell
> setenv BULKFORMAT SED
> head -2  testuser | extractip | geoUR
s/"93.2.135.70"/"France AS15557 Societe Francaise Du Radiotelephone - SFR SA"/g
s/"176.134.198.13"/"France AS5410 Bouygues Telecom SA"/g

```
