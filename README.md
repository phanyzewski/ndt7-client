# ndt7-client
Example Client Implementation of the ndt7 spec.

run locally

```sh
> go get -v github.com/phanyzewski/ndt7-client
> which dt7-clien
$GOPATH/bin/ndt7-client

> ndt7-client
```

example output

```sh
starting download
  214.9 Mbit/s
download: complete

starting upload
  206.6 Mbit/s
upload: complete

+------------------------+---------------+
|      MEASUREMENT       |     VALUE     |
+------------------------+---------------+
| Average Download Speed | 214.93 Mbit/s |
| Retrans Percent        | 0.7884 %      |
| MinRTT                 | 3.5130 ms     |
| Average Upload Speed   | 206.58 Mbit/s |
+------------------------+---------------+
```
