# ssmenv

## Description

Manage environment variables with the AWS SSM Parameter Store.

## Usage

### Prepare

In order to run ssmenv, it is necessary to complete the setting of AWS authentication and define the path of Parameter Store to store environment variables.

```sh
AWS_PROFILE=sample-profile
AWS_REGION=ap-northeast-1

SSMENV_PATH=/dotenv/development
```

### Commands

#### ssmenv get

Get a environment variable.

```sh
$ ssmenv get SAMPLE_VALUE_1
SAMPLE_VALUE_1=sapmle_value

$ ssmenv get SAMPLE_VALUE_2
SAMPLE_VALUE_2="sapmle value"
```

#### ssmenv list

List environment variables.

```sh
$ ssmenv list
SAMPLE_VALUE_1=sapmle_value
SAMPLE_VALUE_2="sapmle value"
```

#### ssmenv put

Set environment variables.
Stored in the Parameter Store as `${SSMENV_PATH}/SAMPLE_VALUE`.

```sh
$ ssmenv put -e SAMPLE_VALUE_1=update_value
$ ssmenv get SAMPLE_VALUE_1
SAMPLE_VALUE_1=update_value

$ ssmenv put -e SAMPLE_VALUE_1=update_value -e SAMPLE_VALUE_2="update value"
$ ssmenv list
SAMPLE_VALUE_1=update_value
SAMPLE_VALUE_2="update value"

$ cat .env
SAMPLE_VALUE_1=update_value_by_file
SAMPLE_VALUE_2="update value by file"
$ ssmenv put -f .env
$ ssmenv list
SAMPLE_VALUE_1=update_value_by_file
SAMPLE_VALUE_2="update value by file"

$ ssmenv put -e SAMPLE_VALUE_3=new_value --dry-run
- key: SAMPLE_VALUE_3
  old_value:
  new_value: new_value

$ ssmenv put -e SAMPLE_VALUE_3=new_value --diff
- key: SAMPLE_VALUE_3
  old_value:
  new_value: new_value

$ ssmenv put -f secret.env --diff=key
- key: SAMPLE_CREDENTIAL
```

#### ssmenv run

Command execution with applying environment variables.

```sh
$ ssmenv run echo $SAMPLE_VALUE_1
sapmle_value
```

## Inspired by

- https://github.com/sachaos/s3env
