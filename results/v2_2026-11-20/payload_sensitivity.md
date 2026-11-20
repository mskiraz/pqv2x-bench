# COER payload-size sensitivity (limit 2481 B)

Sizes are COER encoder output for the implemented IEEE 1609.2-2022 subset with hypothetical PQ extension additions; only the application payload varies; the issuing authority uses the same scheme as the pseudonym certificate; region, assurance level and other optional certificate fields are not present.

## Largest application payload that fits

| Scheme | Format | Certificate (B) | SPDU at 39 B payload | Max payload (B) | SPDU at max | First exceeding payload (B) | SPDU there |
|---|---|---:|---:|---:|---:|---:|---:|
| ECDSA P-256 | digest | 144 | 132 | 2386 | 2481 | 2387 | 2482 |
| ECDSA P-256 | certificate | 144 | 270 | 2248 | 2481 | 2249 | 2482 |
| Falcon-512 | digest | 1615 | 736 | 1782 | 2481 | 1783 | 2482 |
| Falcon-512 | certificate | 1615 | 2345 | 174 | 2481 | 175 | 2482 |
| ML-DSA-44 | digest | 3784 | 2490 | 30 | 2481 | 31 | 2482 |
| ML-DSA-44 | certificate | 3784 | 6268 | none | – | 0 | 6229 |
| SLH-DSA-SHA2-128s | digest | 7938 | 7926 | none | – | 0 | 7887 |
| SLH-DSA-SHA2-128s | certificate | 7938 | 15858 | none | – | 0 | 15819 |

## Encoded SPDU size (bytes) by payload length; remaining bytes to the limit in parentheses, negative when over

| Payload (B) | ECDSA P-256 digest | ECDSA P-256 certificate | Falcon-512 digest | Falcon-512 certificate | ML-DSA-44 digest | ML-DSA-44 certificate | SLH-DSA-SHA2-128s digest | SLH-DSA-SHA2-128s certificate |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 39 | 132 (+2349) | 270 (+2211) | 736 (+1745) | 2345 (+136) | 2490 (-9) | 6268 (-3787) | 7926 (-5445) | 15858 (-13377) |
| 64 | 157 (+2324) | 295 (+2186) | 761 (+1720) | 2370 (+111) | 2515 (-34) | 6293 (-3812) | 7951 (-5470) | 15883 (-13402) |
| 100 | 193 (+2288) | 331 (+2150) | 797 (+1684) | 2406 (+75) | 2551 (-70) | 6329 (-3848) | 7987 (-5506) | 15919 (-13438) |
| 128 | 222 (+2259) | 360 (+2121) | 826 (+1655) | 2435 (+46) | 2580 (-99) | 6358 (-3877) | 8016 (-5535) | 15948 (-13467) |
| 160 | 254 (+2227) | 392 (+2089) | 858 (+1623) | 2467 (+14) | 2612 (-131) | 6390 (-3909) | 8048 (-5567) | 15980 (-13499) |
| 170 | 264 (+2217) | 402 (+2079) | 868 (+1613) | 2477 (+4) | 2622 (-141) | 6400 (-3919) | 8058 (-5577) | 15990 (-13509) |
| 174 | 268 (+2213) | 406 (+2075) | 872 (+1609) | 2481 (+0) | 2626 (-145) | 6404 (-3923) | 8062 (-5581) | 15994 (-13513) |
| 175 | 269 (+2212) | 407 (+2074) | 873 (+1608) | 2482 (-1) | 2627 (-146) | 6405 (-3924) | 8063 (-5582) | 15995 (-13514) |
| 176 | 270 (+2211) | 408 (+2073) | 874 (+1607) | 2483 (-2) | 2628 (-147) | 6406 (-3925) | 8064 (-5583) | 15996 (-13515) |
| 200 | 294 (+2187) | 432 (+2049) | 898 (+1583) | 2507 (-26) | 2652 (-171) | 6430 (-3949) | 8088 (-5607) | 16020 (-13539) |
| 256 | 351 (+2130) | 489 (+1992) | 955 (+1526) | 2564 (-83) | 2709 (-228) | 6487 (-4006) | 8145 (-5664) | 16077 (-13596) |
| 400 | 495 (+1986) | 633 (+1848) | 1099 (+1382) | 2708 (-227) | 2853 (-372) | 6631 (-4150) | 8289 (-5808) | 16221 (-13740) |

Every listed size was produced by signing a real message with that payload and checking that the receiver accepts it.
