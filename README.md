# Github Release Puller

Periodic automated downloads of github releases.

## Build

1. Install [Golang](https://golang.org/)

2. Clone source
   ```bash
   git clone https://github.com/ZASXCDFVA/github-release-puller
   ```
   
3. Build
   ```bash
   cd github-release-puller
   go build
   ```

## Usage

#### Start daemon

```bash
./github-release-puller <path/to/configuration>
```

#### Example configuration

```yaml
asset-pullers:
  - name: "Clash Premium Core"
    owner: "Dreamacro"
    repository: "clash"
    tag: "premium"
    destination: "/srv/clash-updater/"
    interval: 3600
    filters:
      - match: "(.*)linux-amd64(.*)"
  - name: "Clash Open Source Core"
    owner: "Dreamacro"
    repository: "clash"
    destination: "/srv/clash-updater/"
    interval: 3600
    filters:
      - match: "(.*)linux-amd64(.*)"
```