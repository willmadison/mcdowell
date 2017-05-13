FROM centurylink/ca-certs
MAINTAINER Will Madison <will@willmadison.com>

ADD mcdowell /mcdowell

ENTRYPOINT ["/mcdowell"]