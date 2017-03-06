# Payroll-index

## Build the docker image

```
docker build -t payroll-index .
```

## Launch the docker image

```
docker run --rm -e AWS_ACCESS_KEY_ID=[ID] -e AWS_SECRET_ACCESS_KEY=[KEY] payroll-index
```
