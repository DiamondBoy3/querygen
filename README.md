## querygen

A powerful Go-based tool for generating search engine queries from Nuclei templates across multiple search engines and severity levels.

## Features
- **Multiple Search Engines**: Support for Shodan, Google, Censys, Fofa, Hunter, and ZoomEye
- **Flexible Severity Levels**: Choose from `all`, `critical`, `high`, `medium`, `low`, and `info` levels
- **Comma-Separated Severities**: Combine multiple severity levels in one command
- **Clean Output**: Converts queries to your preferred search engine format

## Supported Search Engines
| Engine | Output Format | Notes |
|--------|---------------|-------|
| `shodan` | `http.title:` | Default engine |
| `google` | `intitle:` | Requires grep-like handling for binary data |
| `censys` | `services.http.response.html_title=` | |
| `fofa` | `title=` | |
| `hunter` | `web.title=` | |
| `zoomeye` | `title=` | |

## Supported Severity Levels
- `all` - All available templates
- `critical` - Critical severity templates
- `high` - High severity templates  
- `medium` - Medium severity templates
- `low` - Low severity templates
- `info` - Informational templates

## How It Works
1. **Fetches Templates**: Downloads query templates from the [nucleihubquery](https://github.com/rix4uni/nucleihubquery) GitHub repository
2. **Processes by Severity**: Filters templates based on selected severity levels
3. **Converts Format**: Transforms queries from source format to target search engine format
4. **Outputs Results**: Streams converted queries to stdout for immediate use

## Installation
```
go install github.com/rix4uni/querygen@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/querygen/releases/download/v0.0.1/querygen-linux-amd64-0.0.1.tgz
tar -xvzf querygen-linux-amd64-0.0.1.tgz
rm -rf querygen-linux-amd64-0.0.1.tgz
mv querygen ~/go/bin/querygen
```
Or download [binary release](https://github.com/rix4uni/querygen/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/querygen.git
cd querygen; go install
```

## Usage
```yaml
Usage of querygen:
  -e, --engine string     Search engine: shodan, google, censys, fofa, hunter, zoomeye (default "shodan")
  -s, --severity string   Severity level: all, critical, high, medium, low, info (comma-separated) (default "all")
      --silent            Silent mode.
      --version           Print the version of the tool and exit.
```

## Usage Examples
```yaml
# Generate Shodan queries for all severity levels (default)
querygen --engine shodan --severity all

# Generate Google queries for critical severity
querygen --engine google --severity critical

# Generate Fofa queries with short flags
querygen -e fofa -s critical
```

### Multiple Severities
```yaml
# Combine critical and high severity queries
querygen --engine shodan --severity critical,high

# Mix multiple severity levels
querygen --engine google --severity critical,high,medium

# All severity levels (default)
querygen --engine shodan --severity all
```

## Output Example
```yaml
▶ querygen --silent --engine shodan --severity critical | unew
http.title:="综合安防管理平台
http.title:="综合安防管理平台"
http.title:="金蝶云星空 管理中心"
http.title:="微信管理后台"
http.title:="SANGFOR上网优化管理"
http.title:"Dashboard [Jenkins]"
http.title:"Script Console [Jenkins]"
http.title:M1-Server
http.title:"AnalyticsCloud 分析云"
http.title:"omnipcx for enterprise"
http.title:"zeroshell"
http.title:"struts2 showcase"
http.title:"manageengine desktop central 10"
http.title:"manageengine desktop central"
http.title:"zabbix-server"
http.title:"oracle peoplesoft sign-in"
http.title:"apache solr"
http.title:"solr admin"
http.title:"kentico database setup"
http.title:"active management technology"
http.title:"cobbler web interface"
http.title:"roteador wireless"
http.title:"coldfusion administrator login"
http.title:"fuel cms"
http.title:"kibana"
```