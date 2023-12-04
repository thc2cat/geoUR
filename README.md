# geoUR

Using maxmind geoip databases, embed Country and ASN data in a portable binary, and reading inputs of ipv4 addresses on stdin, output Country and ASN numbers

Maxmind license_key file should be present in current directory, containing 'licence_key=xxxxx' in order to download Maxmind mmdb databases via init.sh script.

Most of the job is done using [oschwald/geoip2-golang lib]("https://github.com/oschwald/geoip2-golang").

## Build (need golang installed)

```shell
> make build
```

On *nix, init.sh should download maxmind dbs in data/ dir and produce init.go , wich mainly contains names of the assets.
Your can also manually download these db from maxmind and adapt init.tpl manually.

## Usage

```shell
> geoUR.exe -h             
=== geoUR.exe  embedded geoip2 databases
        - GeoLite2-Country_xxxx
        - GeoLite2-ASN_xxxx

Usage: geoUR.exe ip or stdin input

> geoUR.exe < testdata.ips
...
> geoUR 217.69.139.150
217.69.139.150 (Russia / AS47764 LLC VK)
```

ip addresses must be extracted before ( see my extractip repo), as exemple :

```shell
> rg "pattern" logs/* | extractip | geoUR | rg "Russia|COLOCROSSING"
```
