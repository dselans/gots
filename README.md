gots
====
`gots` is a (mostly crude) UNIX timestamp manipulation CLI utility.

In addition to converting between human readable time/date format and UNIX timestamps, it also allows you to increment and decrement time by seconds, minutes, hours, days, months and years.

### Installation
```
go install github.com/dselans/gots@latest
```

### Usage
```
$ gots -h
Usage: ./gots [-h] [date_string|unix_timestamp|[+|-123s|m|h|d|M|y]]
$ gots
1424657615 <=> Sun Feb 22 18:13:35 PST 2015
$ gots +10d
1425521628 <=> Wed Mar  4 18:13:48 PST 2015
$ gots 1234567890
1234567890 <=> Fri Feb 13 15:31:30 PST 2009
$ gots "01/02/2014 23:59:02"
1388707142 <=> Thu Jan  2 23:59:02 UTC 2014
```

Note that in order to convert dates from "string" format (ie. "01/02/2014 01:02:03"), the provided time format must match it exactly.

### Accuracy
The generated timestamps can be a bit inaccurate - namely because I decided to not meddle with timezones, DST or any other time-related pains. This works for me, but may be a problem for you; if so, please send a PR!
