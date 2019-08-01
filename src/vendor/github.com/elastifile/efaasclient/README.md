# eFaaS client for Go

## Development process highlights

Prior to pushing your changes, make sure all unit tests pass (you may need to update eFaaS URL, project number, service account key etc.):
```bash
go build # Make sure there are no compilation errors
go test -v -timeout 60m
```

## Release process highlights

```bash
git tag vX.Y.X
git push --tags origin master
```

## Usage

### Highlights

In order to use private repo, you may need to
1. Make sure you have access to GutHub via SSH
2. Configure git to use SSH instead of HTTPS that is used with go modules by default
```bash
git config --global url.ssh://git@github.com/.insteadOf https://github.com/
```

### Example

```go
    opts := ClientCreateOpts{
        ProjectNumber: "<Your GCP project's numeric id>",
        BaseURL:       "https://cloud-file-service-gcp.elastifile.com",
    }
    client, _ := NewClient("<Your service account's key in JSON format>", opts)
    instances, _ := client.ListInstances()
```
Note: error handling was skipped for brevity's sake
