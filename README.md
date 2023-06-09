kLeZ/dir2cm
===============

[![Build Status](https://travis-ci.org/kLeZ/dir2cm.svg?branch=master)](https://travis-ci.org/kLeZ/dir2cm)

Creates a Kubernetes `ConfigMap` from a directory. Due to limitations with `ConfigMap`,
only files (not directories) in the first-level directory are processed.

```
$ mkdir -p configfiles
$ for i in $(seq 1 5); do \
    echo "This is file $i" >> "configfiles/file_${i}.txt"; \
  done
$ dir2cm -dir configfiles -name my-configs
```

Generates a ConfigMap like so:
```
---
apiversion: v1
kind: ConfigMap
metadata:
  name: my-configs
  labels: {}
data:
  file_1.txt: |
    This is file 1
  file_2.txt: |
    This is file 2
  file_3.txt: |
    This is file 3
  file_4.txt: |
    This is file 4
  file_5.txt: |
    This is file 5
```

Docker
------

Given the setup above, the command

    docker run --rm -v "$(pwd)/configfiles:/data" kLeZ/dir2cm -name foobar

Should dump the following to stdout:

    apiversion: v1
    kind: ConfigMap
    metadata:
      name: foobar
      labels: {}
    data:
      file_1.txt: |
        This is file 1
      file_2.txt: |
        This is file 2
      file_3.txt: |
        This is file 3
      file_4.txt: |
        This is file 4
      file_5.txt: |
        This is file 5

