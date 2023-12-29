<p align="center">
    <img width="256" src="assets/logo.png">
</p>

# Troll-A

[![License: Apache](https://img.shields.io/github/license/crissyfield/troll-a?color=orange)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/crissyfield/troll-a)](https://goreportcard.com/report/github.com/crissyfield/troll-a)
[![Go Reference](https://pkg.go.dev/badge/github.com/crissyfield/troll-a.svg)](https://pkg.go.dev/github.com/crissyfield/troll-a)

`Troll-A` is a command line tool for extracting secrets such as passwords, API keys, and tokens from WARC (Web
ARChive) files. `Troll-A` is an **easy-to-use**, **comprehensive**, and **fast** solution for finding
secrets in web archive.


## Features

- **Protocols:** Supports retrieving web archives directly from a network server via HTTP/HTTPS, from the
  [Amazon S3](https://aws.amazon.com/pm/serv-s3/) object storage service, or from the local file system.
- **Compression:** Supports web archives compressed with [GZip](https://www.gzip.org),
  [BZip2](https://sourceware.org/bzip2/), [XZ](https://github.com/tukaani-project/xz), or 
  [ZStd](https://github.com/facebook/zstd). For ZStd, it also supports custom dictionaries prepended to the
  compressed data stream (as used by `*.megawarc.warc.zst` files).
- **Comprehensive:** Uses the battle-tested ruleset from the [Gitleaks](https://gitleaks.io) project to
  detect up to 166 different types of secrets, tokens, keys, or other sensitive information.
- **Performance:** Works concurrently and optionally uses optimized regular expressions (via
  [go-re2](https://github.com/wasilibs/go-re2)) to process a typical [Common Crawl](https://commoncrawl.org)
  web archive (~34.000 pages) in less than 30 seconds on AWS `c7g.12xlarge`.


## Installation

### Download

`Troll-A` is available in binary form for macOS and Linux on the
[releases page](https://github.com/crissyfield/troll-a/releases).

### Building

For better performance, it is recommended to build `Troll-A` from source, as this allows to use optimized
regular expression engine provided by [go-re2](https://github.com/wasilibs/go-re2). For this to work, the
[RE2](https://github.com/google/re2) dependency must be installed first.

#### macOS

```
# Install dependencies
brew install re2

# Install with RE2 activated
go install -tags re2_cgo github.com/crissyfield/troll-a@v1.0.0
```

#### Debian / Ubuntu

```
# Install dependencies
sudo apt install -u build-essential libre2-dev

# Install with RE2 activated
go install -tags re2_cgo github.com/crissyfield/troll-a@v1.0.0
```


## Usage

```
Usage:
  troll-a [flags] url

This tool allows to extract (potential) secrets such as passwords, API keys, and tokens
from WARC (Web ARChive) files. Extracted information is output as structured text org
JSON, which simplifies further processing of the data.

"url" can be either a regular HTTP or HTTPS reference ("https://domain/path"), an Amazon
S3 reference ("s3://bucket/path"), or a file path (either "file:///path" or simply
"path"). If the data is compressed with either GZip, BZip2, or ZStd it is automatically
decompressed. ZStd with a prepended custom dictionary (as used by "*.megawarc.warc.zstd")
is also handled transparently.

This tool uses rules from the Gitleaks project (https://gitleaks.io) to detect secrets.

Flags:
  -e, --enclosed               only report secrets that are enclosed within their context
  -h, --help                   help for troll-a
  -j, --jobs uint              detect secrets with this many concurrent jobs (default 8)
  -s, --json                   output detected secrets as JSON
  -p, --preset rules-preset    rules preset to use. This could be one of the following:
                               all:         All known rules will be applied, which can
                                            result in a significant amount of noise for
                                            large data sets.
                               most:        Most of the rules are applied, skipping the
                                            biggest culprits for false positives.
                               secret:      Only rules are applied that are most likely
                                            to result in an actual leak of a secret.
                               No other values are allowed. (default secret)
  -q, --quiet                  suppress success message(s)
  -r, --retry retry-strategy   retry strategy to use. This could be one of the following:
                               never:       This strategy will fail after the first fetch
                                            failure and will not attempt to retry.
                               constant:    This strategy will attempt to retry up to 5
                                            times, with a 5s delay after each attempt.
                               exponential: This strategy will attempt to retry for 15
                                            minutes, with an exponentially increasing
                                            delay after each attempt.
                               always:      This strategy will attempt to retry forever,
                                            with no delay at all after each attempt.
                               No other values are allowed. (default never)
  -t, --timeout duration       fetching timeout (does not apply to files) (default 30m0s)
  -v, --version                version for troll-a
```


## Examples

### Common Crawl

[Common Crawl](https://commoncrawl.org) maintains a free, open repository of web crawl data that can be used by
anyone. The Common Crawl corpus contains petabytes of data collected regularly since 2008.

For example, to extract secrets from all of the 3.35 billion pages of the November/December 2023 crawl
(called `CC-MAIN-2023-50`), you can do this:

```bash
# Download the list of all 90.000 WARC paths
curl -sSL -O https://data.commoncrawl.org/crawl-data/CC-MAIN-2023-50/warc.paths.gz

# Iterate through all paths using 64 scanning jobs, output matches as JSON
gzcat warc.paths.gz | \
xargs -I{} -- troll-a -e -s -j64 https://data.commoncrawl.org/{} > secrets.json
```

> [!WARNING]
> This will take a long time! Depending on your hardware and Internet connection, this can take anywhere from
> a week to several months. You may want to run this example only for the first few lines of `warc.paths.gz`.

### Internet Archive

The [Archive Team](http://archiveteam.org/index.php) is a group dedicated to digital preservation and web
archiving founded in 2009. Web archives are stored as WARC files (more specifically, in MegaWARC format) and
made available through the [Internet Archive](https://archive.org/details/archiveteam).

For example, to extract secrets from the 113.372 pages the Archive Team crawled from
[pastebin.com](https://pastebin.com) in April of 2023 (here's the corresponding
[publication](https://archive.org/details/archiveteam_pastebin_20230421003309_a3b951b4) on the Internet
Archive), you can do this:

```bash
# Call troll-a directly with the MegaWARC URL
troll-a -e https://archive.org/download/archiveteam_pastebin_20230421003309_a3b951b4/pastebin_20230421003309_a3b951b4.1603050931.megawarc.warc.zst
```

...which results in...

```
Detected: secret="acf30fb56amsh654fa8104418601p1e420cjsn3152a0032f0b" rule="rapidapi-access-token" uri="https://pastebin.com/raw/bKMJXkQE" line=36 column=15
Detected: secret="acf30fb56amsh654fa8104418601p1e420cjsn3152a0032f0b" rule="rapidapi-access-token" uri="https://pastebin.com/raw/bKMJXkQE" line=36 column=15
Detected: secret="acf30fb56amsh654fa8104418601p1e420cjsn3152a0032f0b" rule="rapidapi-access-token" uri="https://pastebin.com/raw/nferefe2" line=37 column=6
Detected: secret="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" rule="github-pat" uri="https://pastebin.com/print/cQEA2GCS" line=39 column=123
Detected: secret="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" rule="github-pat" uri="https://pastebin.com/embed_js/cQEA2GCS" line=11 column=2688
Detected: secret="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" rule="github-pat" uri="https://pastebin.com/embed_iframe/cQEA2GCS?theme=dark" line=49 column=123
Detected: secret="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" rule="github-pat" uri="https://pastebin.com/cQEA2GCS" line=222 column=123
Detected: secret="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" rule="github-pat" uri="https://pastebin.com/raw/cQEA2GCS" line=22 column=22
Detected: secret="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" rule="github-pat" uri="https://pastebin.com/embed_iframe/cQEA2GCS" line=48 column=123
Detected: secret="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" rule="github-pat" uri="https://pastebin.com/embed_js/cQEA2GCS?theme=dark" line=11 column=2796
Detected: secret="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" rule="github-pat" uri="https://pastebin.com/clone/cQEA2GCS" line=152 column=27
Success: Processed https://archive.org/download/archiveteam_pastebin_20230421003309_a3b951b4/pastebin_20230421003309_a3b951b4.1603050931.megawarc.warc.zst (113372 records)
```


## Credits

The set of rules used to detect the actual secrets is part of the [Gitleaks](https://gitleaks.io) project. We
are very grateful for the tremendous work they have done in compiling all this information!


## What's up with the name?

The [Troll A platform](https://en.wikipedia.org/wiki/Troll_A_platform) is a natural gas platform in the Troll
gas field off the west coast of Norway. As of 2014, it was the tallest structure that has ever been moved to
another position, relative to the surface of the Earth, and is among the largest and most complex engineering
projects in history. In 1996, the platform set the Guinness World Record for the largest offshore gas platform.

> [!NOTE]
> While we deeply dislike the exploitation of natural resources, we admire the engineering feat!
