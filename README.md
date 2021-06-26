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
05:01:51 PM ❯ ndt7-client
starting download
Avg. speed  :   196.8 Mbit/s
Avg. speed  :   209.0 Mbit/s
Avg. speed  :   214.1 Mbit/s
Avg. speed  :   211.4 Mbit/s

...

download: complete
starting upload
Avg. speed  :   216.6 Mbit/s
Avg. speed  :   239.9 Mbit/s
Avg. speed  :   233.0 Mbit/s
Avg. speed  :   229.4 Mbit/s

...

upload: complete
+------------------------+---------------------------+
|      MEASUREMENT       |           VALUE           |
+------------------------+---------------------------+
| Average Download Speed | 213.94406466171014 Mbit/s |
| Retrans Percent        |        1.1285983480443076 |
| MinRTT                 |                     3.447 |
| Average Upload Speed   | 205.73967321476167 Mbit/s |
+------------------------+---------------------------+
```
