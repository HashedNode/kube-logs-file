# Kube-Logs-File

Util for download and write logs from Kubernetes Pods to file.
## Usage/Examples

From the cloned source code, move inside the project and in the shell type:
```bash
go run main.go --context=example --namespace=example-namespace --pods=podName1,podName2
```

Or you can run the executable after download it:

```bash
./kube-logs-file --context=example --namespace=example-namespace --pods=podName1,podName2
```

For a guide aboute the usage run
```bash
./kube-logs-file --help
Or
go run main.go --help 
```
## Arguments References

#### Path to Kube conf file

```bash
  --kubeconfig
```

| Optional | Description                |
| :------- | :------------------------- |
| `true`   | Path to the kubeconfig file, it will use default path if not provided|



#### Context

```bash
  --context
```

| Optional | Description                       |
| :------- | :-------------------------------- |
| `true`   | The context of kubernetes to use, if not passed will be used the current active context|


#### Namespace

```bash
  --namespace
```

| Optional | Description                       |
| :------- | :-------------------------------- |
| `false`  | Namespace where to get the pods|


#### Pods Name

```bash
  --pods
```

| Optional | Description                       |
| :------- | :-------------------------------- |
| `false`  | Comma-separated list of pods names with no spaces, if you want to put spaces between pods names just escape with \"\" example \"podname1, podname2\"|

## Authors

- [@HashedNode](https://github.com/HashedNode)


## License

[MIT](https://choosealicense.com/licenses/mit/)

