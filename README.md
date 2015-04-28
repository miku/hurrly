README
======

Find the destination location for DOIs.

Usage
-----

    $ hurrly < urls.txt > urls.ldj

Input is a list of DOI API URLs to check, one URL per line. Output is TSV with
the redirect location and some more information.

    $ hurrly < fixtures/10.txt
    200 OK  0.8956  1430238589  http://.../10.1590/s0100-41582006000200010  http://www.scielo.br/...
    200 OK  0.1826  1430238589  http://.../10.1590/s0102-47442006000100006  http://www.scielo.br/...
    ...
    404 Not Found   0.3241  1430238589  http://doi.org/api/handles/10.sss1590/s0102-47442006000100006
