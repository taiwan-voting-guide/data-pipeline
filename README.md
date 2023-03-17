# Taiwan Voting Guide

## Set up

1. sign up [prefect](https://www.prefect.io/) and create a project
2. `prefect cloud login`
3. `python data_pipeline/legislators.py`

## Golang

### Setup
```
# install golang 1.19
cp .env.example golang/.env
python scripts/vote_data.py
```

### Commands
```sh
# party data
go run .golang/cmd/parties/main.go

# politician data

```
