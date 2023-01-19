gots
====
`gots` is a crude UNIX timestamp manipulation CLI utility.

In addition to converting between human readable time/date format and UNIX timestamps, it also allows you to increment and decrement time by seconds, minutes, hours, days, months and years. Supports unix nanos.

### Usage
```
❯ gots
RFC3339: 2023-01-18T21:36:37-08:00
Unix: 1674106597
UnixNano: 1674106597905104000
---------------------
RFC3339 (UTC): 2023-01-19T05:36:37Z
Unix (UTC): 1674106597
UnixNano (UTC): 1674106597905104000

❯ gots 1674106597
RFC3339: 2023-01-18T21:36:37-08:00
Unix: 1674106597
UnixNano: 1674106597000000000
---------------------
RFC3339 (UTC): 2023-01-19T05:36:37Z
Unix (UTC): 1674106597
UnixNano (UTC): 1674106597000000000

❯ gots 1674106597905104000 +1y
handling time shift
RFC3339: 2024-01-18T21:36:37-08:00
Unix: 1705642597
UnixNano: 1705642597905104000
---------------------
RFC3339 (UTC): 2024-01-19T05:36:37Z
Unix (UTC): 1705642597
UnixNano (UTC): 1705642597905104000
```

Note that in order to convert dates from "string" format (ie. "01/02/2014 01:02:03"), the provided time format must match it exactly.

### Accuracy
The generated timestamps can be a bit inaccurate - namely because I decided to not meddle with timezones, DST or any other time-related pains. This works for me, but may be a problem for you.


NOTE: I wrote this in 2015 when I just picked up Go :)
