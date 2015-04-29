README
======

Mass DOI resolver and checker.

Installation
------------

Releases: https://github.com/miku/hurrly/releases

Or the go tool:

    $ go get github.com/miku/hurrly/cmd/hurrly

Usage
-----

Input is a list of DOI API URLs to check, one URL per line.
Only URL of the form `http://doi.org/api/handles/<DOI>` are accepted.

    $ cat fixtures/urls.txt
    http://doi.org/api/handles/10.1590/s0100-41582006000200010
    http://doi.org/api/handles/10.1590/s0034-71402005000400003
    http://doi.org/api/handles/10.1590/s0034-71402005000400005
    http://doi.org/api/handles/10.1590/s0100-512x2004000200004
    http://doi.org/api/handles/10.1590/s0100-512x2004000200005
    http://doi.org/api/handles/10.1590/s0102-47442006000100001
    http://doi.org/api/handles/10.1590/s0102-47442006000100002
    http://doi.org/api/handles/10.1590/s0102-47442006000100003
    http://doi.org/api/handles/10.1590/s0102-47442006000100005
    http://doi.org/api/handles/10.1590/s0102-47442006000100006
    http://doi.org/api/handles/10.xxxx1590xxxx/nononoxxxx

Output are TSV with the redirect location and additional information (status, request time, epoch, url, redirect).

    $ hurrly < fixtures/urls.txt
    200 OK  0.4543  1430316353  http://doi.org/...0004  http://www.scie...g=pt
    200 OK  0.4542  1430316353  http://doi.org/...0005  http://www.scie...g=pt
    200 OK  0.4546  1430316353  http://doi.org/...0003  http://www.scie...g=pt
    ...
    404 Not Found   0.3166  1430316354  http://doi.org/api/...nononoxxxx   NOT_AVAILABLE

Hurrly will try hard to get a result for each URL. It will write the status of the result
into the first column, either as HTTP status message, like `200 OK`, `404 Not Found` or as internal error message,
like `E_REQ`, `E_READ`, `E_JSON`, etc.

To run things in parallel, adjust the `-w` parameter.

You can use `hurrly` to look up single DOI links as well:

    $ echo "http://doi.org/api/handles/10.1021/la025770y" | hurrly | cut -f5
    http://pubs.acs.org/doi/abs/10.1021/la025770y
