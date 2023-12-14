<p align="center">
    <img width="256" src="assets/logo.png">
</p>

# Troll-A

[![Go Reference](https://pkg.go.dev/badge/github.com/crissyfield/troll-a.svg)](https://pkg.go.dev/github.com/crissyfield/troll-a)
[![Go Report Card](https://goreportcard.com/badge/github.com/crissyfield/troll-a)](https://goreportcard.com/report/github.com/crissyfield/troll-a)

## Performance

```bash
time ./troll-a \
    --jobs=64 \
    --preset=secret \
    https://data.commoncrawl.org/crawl-data/CC-MAIN-2023-40/segments/1695233510326.82/warc/CC-MAIN-20230927203115-20230927233115-00771.warc.gz
```

|  Instance Type   |    Arch     |  Jobs  |  Duration  |  Instance Cost  |    Total Cost    |    Total Time    |
|:-----------------|:-----------:|:------:|:----------:|----------------:|-----------------:|-----------------:|
|  `c7g.16xlarge`  |  `aarch64`  |  `16`  |  1:10.078  |  2.32000 USD/h  |    406.45 USD    |     7d 7h 12m    |
|  `c7g.12xlarge`  |  `aarch64`  |  `16`  |  1:10.217  |  1.74000 USD/h  |    305.44 USD    |     7d 7h 33m    |
|  `r7g.8xlarge`   |  `aarch64`  |  `16`  |  1:10.502  |  1.71360 USD/h  |    302.03 USD    |     7d 8h 15m    |
|  `c6a.8xlarge`   |   `x86_64`  |  `16`  |  1:32.323  |  1.22400 USD/h  |    282.23 USD    |    9d 14h 48m    |
|  `m7g.8xlarge`   |  `aarch64`  |  `16`  |  1:10.220  |  1.30560 USD/h  |    229.20 USD    |     7d 7h 33m    |
|  `c6a.8xlarge`   |   `x86_64`  |  `32`  |  1:13:498  |  1.22400 USD/h  |    224.90 USD    |    7d 15h 45m    |
|  `c6a.8xlarge`   |   `x86_64`  |  `64`  |  1:13.377  |  1.22400 USD/h  |    224.53 USD    |    7d 15h 27m    |
|  `c6a.4xlarge`   |   `x86_64`  |  `16`  |  2:24.422  |  0.61200 USD/h  |    220.97 USD    |     15d 1h 3m    |
|  `c6a.4xlarge`   |   `x86_64`  |  `32`  |  2:23.686  |  0.61200 USD/h  |    219.84 USD    |   14d 23h 13m    |
|  `c6a.4xlarge`   |   `x86_64`  |  `64`  |  2:23.563  |  0.61200 USD/h  |    219.65 USD    |   14d 22h 54m    |
|  `c7g.16xlarge`  |  `aarch64`  |  `32`  |  0:36.644  |  2.32000 USD/h  |    212.53 USD    |  **3d 19h 37m**  |
|  `c7g.8xlarge`   |  `aarch64`  |  `16`  |  1:10.427  |  1.16000 USD/h  |    204.24 USD    |      7d 8h 4m    |
|  `r7g.8xlarge`   |  `aarch64`  |  `32`  |  0:38.103  |  1.71360 USD/h  |    163.23 USD    |    3d 23h 15m    |
|  `c7g.12xlarge`  |  `aarch64`  |  `32`  |  0:36.759  |  1.74000 USD/h  |    159.90 USD    |  **3d 19h 54m**  |
|  `r7g.8xlarge`   |  `aarch64`  |  `64`  |  0:37.063  |  1.71360 USD/h  |    158.78 USD    |  **3d 20h 39m**  |
|  `c7g.16xlarge`  |  `aarch64`  |  `64`  |  0:25.209  |  2.32000 USD/h  |    146.21 USD    |   **2d 15h 1m**  |
|  `c7g.16xlarge`  |  `aarch64`  |  `96`  |  0:25.196  |  2.32000 USD/h  |    146,14 USD    |  **2d 14h 59m**  |
|  `m7g.8xlarge`   |  `aarch64`  |  `32`  |  0:38.078  |  1.30560 USD/h  |    124.29 USD    |    3d 23h 12m    |
|  `c7g.12xlarge`  |  `aarch64`  |  `64`  |  0:27.754  |  1.74000 USD/h  |    120.72 USD    |  **2d 21h 23m**  |
|  `m7g.8xlarge`   |  `aarch64`  |  `64`  |  0:36.801  |  1.30560 USD/h  |  **120.12 USD**  |   **3d 20h 0m**  |
|  `c7g.12xlarge`  |  `aarch64`  |  `96`  |  0:27.428  |  1.74000 USD/h  |  **119.31 USD**  |  **2d 20h 34m**  |
|  `c7g.8xlarge`   |  `aarch64`  |  `32`  |  0:37.938  |  1.16000 USD/h  |  **110.02 USD**  |  **3d 22h 51m**  |
|  `c7g.8xlarge`   |  `aarch64`  |  `64`  |  0:36.954  |  1.16000 USD/h  |  **107.17 USD**  |  **3d 20h 23m**  |
|  `c7g.4xlarge`   |  `aarch64`  |  `16`  |  1:13.004  |  0.58000 USD/h  |  **105.86 USD**  |    7d 14h 31m    |
|  `c7g.2xlarge`   |  `aarch64`  |  `16`  |  2:22.966  |  0.29000 USD/h  |  **103.65 USD**  |   14d 21h 25m    |
|  `c7g.2xlarge`   |  `aarch64`  |  `32`  |  2:22.924  |  0.29000 USD/h  |  **103.62 USD**  |   14d 21h 19m    |
|  `c7g.4xlarge`   |  `aarch64`  |  `32`  |  1:11.450  |  0.58000 USD/h  |  **103.60 USD**  |    7d 10h 38m    |
|  `c7g.2xlarge`   |  `aarch64`  |  `64`  |  2:22.824  |  0.29000 USD/h  |  **103.55 USD**  |    14d 21h 4m    |
|  `c7g.4xlarge`   |  `aarch64`  |  `64`  |  1:11.135  |  0.58000 USD/h  |  **103.15 USD**  |     7d 9h 50m    |
