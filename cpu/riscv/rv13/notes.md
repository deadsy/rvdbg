# Notes on Expected Performance

## Absolute Maximums

TCK clock rate = T

If all clocks gave us 1 bit of useful data:

rate = T / 8

E.g. T = 4e6 Hz

rate = 4e6 / (8 * 1024) = 488 KiB/sec

## DMI Operation Overhead

For a dmi operation:

op = preamble + dmi + postamble

dmi = abits + 32 + 2

preamble = 3 (idle to dr shift)

postamble = 3 + idle (dr shift to idle)

Eg. idle = 5 bits

Eg. abits = 5 bits

dmi = 3 + 5 + 32 (data) + 2 + 5 + 3

So data = 32 of 50 cyles

rate = 488 * (32/50) = 312 KiB/sec

## Program Buffer Setup Overhead






