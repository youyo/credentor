# credentor

[![Go project version](https://badge.fury.io/go/github.com%2Fyouyo%2Fcredentor.svg)](https://badge.fury.io/go/github.com%2Fyouyo%2Fcredentor)
[![Go Report Card](https://goreportcard.com/badge/github.com/youyo/credentor)](https://goreportcard.com/report/github.com/youyo/credentor)
[![License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](./LICENSE)
[![CircleCI](https://circleci.com/gh/youyo/credentor/tree/master.svg?style=svg)](https://circleci.com/gh/youyo/credentor/tree/master)

AWS assume role credential wrapper.

## Description

credentor is useful for some commands which couldn't resolve an assume role credentials in ~/.aws/credentials.

For example,

- [Terraform](https://www.terraform.io/)
- [Packer](https://www.packer.io/)
- etc.

## Install

Place a `credentor` command to your PATH and set an executable flag.  
Download the latest release from github. https://github.com/youyo/credentor/releases/latest

```console
# darwin/amd64
$ curl -s https://api.github.com/repos/youyo/credentor/releases/latest \
	| grep "browser_download_url.*darwin" \
	| cut -d : -f 2,3 \
	| tr -d \" \
	| wget -qi -

# linux/amd64
$ curl -s https://api.github.com/repos/youyo/credentor/releases/latest \
	| grep "browser_download_url.*linux" \
	| cut -d : -f 2,3 \
	| tr -d \" \
	| wget -qi -
```

## Usage

```ini
# ~/.aws/credentials

[my-profile]
aws_access_key_id=XXX
aws_secret_access_key=YYY
```

```ini
# ~/.aws/config

[profile foo]
role_arn=arn:aws:iam::999999999999:role/MyRole
source_profile=my-profile
```

### As command wrapper

```console
$ AWS_PROFILE=foo credentor -- some_command [arg1 arg2...]
```

`credentor` works as below.

1. Find `AWS_PROFILE` section in ~/.aws/credentials and ~/.aws/config .
2. Call `aws sts assume-role` to a get temporary credentials.
3. Set the credentilas to environment variables.
   - `AWS_ACCESS_KEY_ID`
   - `AWS_SECRET_ACCESS_KEY`
   - `AWS_SESSION_TOKEN`
4. Execute `some_command` with args.

### As env exporter

When credentor is executed with no arguments, credentor outputs shell script to export AWS credentials environment variables.

```console
$ export AWS_PROFILE=foo
$ credentor
export AWS_ACCESS_KEY_ID=XXXXXXXXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=zWarBXUtMKJYnC8y4dNAf9e5HQqFTp....
export AWS_SESSION_TOKEN=Wj3YGuSMwn8aJx4AN6TFsbtB5URKHEpVgdDkPvy7....
```

You can set the credentials in current shell by `eval`.

```console
$ eval "$(credentor)"
```

Temporary credentials has expiration time (about 1 hour).

## References

credentor is inspired by aswrap.  
Original software is [aswrap](https://github.com/fujiwara/aswrap). https://github.com/fujiwara/aswrap
