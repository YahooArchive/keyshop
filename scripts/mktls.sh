#!/usr/bin/env sh
set -x
domain="localhost"
domains="${domain} $(id -u -n).local"
tls=data/tls
if test -d ${tls} ; then
      rm ${tls}/*     
else
      mkdir ${tls}
fi
go run ks/cmd/localcert/main.go ${domains} &&
      cat ${domain}.server.pem ${domain}.ca.pem > ${domain}.chain.pem &&
      mv ${domain}* ${tls}
