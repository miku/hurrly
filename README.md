README
======

Mass DOI resolver and checker.

Usage
-----

    $ hurrly < urls.txt > urls.tsv

Input is a list of DOI API URLs to check, one URL per line.
Only URL of the form `http://doi.org/api/handles/<DOI>` are accepted.

Output is TSV with the redirect location and some more information: status, request time, epoch, url, redirect.

    $ hurrly < fixtures/10.txt
    200 OK  0.8956  1430238589  http://.../10.1590/s0100-415820060...  http://www.scie...
    200 OK  0.1826  1430238589  http://.../10.1590/s0102-474420060...  http://www.scie...
    ...
    404 Not Found   0.3241  1430238589  http://.../10.sss1590/s0102-4744201...

Hurrly will try hard to get a result for each URL. It will write the status of the result
into the first column, either as HTTP status message, like `200 OK` or as internal error message,
like `E_REQ`, `E_READ`, `E_JSON`, etc.

To run things in parallel, adjust the `-w` parameter.
