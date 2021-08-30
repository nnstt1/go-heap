# go-heap

Kubernetes の Scheduler で実装されている heap を元に、Go の理解を深めるためのリポジトリです。

## モジュール

`k8s.io/kubernetes/pkg/scheduler/framework` が参照する `k8s.io/api` のバージョンが `v0.0.0` となっておりビルドで失敗しました。

https://github.com/kubernetes/kubernetes/issues/79384 を参照に `go-mod-kube.sh` スクリプトを実行することで、指定バージョンのモジュールを使用することができました。

```bash
./go-mod-kube.sh 1.22.0
go mod tidy
go build
```
