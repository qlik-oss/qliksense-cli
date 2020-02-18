module github.com/qlik-oss/sense-installer

go 1.13

replace (
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a

	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191016120415-2ed914427d51

	sigs.k8s.io/kustomize/api => github.com/qlik-oss/kustomize/api v0.3.3-0.20200206224201-2e697eccbad9
)

require (
	cloud.google.com/go v0.52.0 // indirect
	cloud.google.com/go/storage v1.5.0 // indirect
	github.com/Shopify/ejson v1.2.1
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/aws/aws-sdk-go v1.28.9 // indirect
	github.com/bitly/go-hostpool v0.0.0-20171023180738-a3a6125de932 // indirect
	github.com/bugsnag/bugsnag-go v1.5.3 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cloudflare/cfssl v1.4.1 // indirect
	github.com/containerd/containerd v1.3.2 // indirect
	github.com/containerd/continuity v0.0.0-20191214063359-1097c8bae83b // indirect
	github.com/containers/image/v5 v5.1.0
	github.com/docker/cli v0.0.0-20191212191748-ebca1413117a
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.4.2-0.20190924003213-a8608b5b67c7
	github.com/docker/go v1.5.1-1 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobuffalo/logger v1.0.3 // indirect
	github.com/gobuffalo/packd v1.0.0 // indirect
	github.com/gobuffalo/packr/v2 v2.7.1
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect

	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/jinzhu/gorm v1.9.11 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/mattn/go-sqlite3 v2.0.1+incompatible // indirect
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/qlik-oss/k-apis v0.0.8
	github.com/rogpeppe/go-internal v1.5.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.1
	github.com/theupdateframework/notary v0.6.1 // indirect
	github.com/xlab/handysort v0.0.0-20150421192137-fb3537ed64a1 // indirect
	golang.org/x/crypto v0.0.0-20200214034016-1d94cc7ab1c6
	golang.org/x/exp v0.0.0-20200119233911-0405dc783f0a // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/sys v0.0.0-20200217220822-9197077df867 // indirect
	golang.org/x/tools v0.0.0-20200216192241-b320d3a0f5a2 // indirect
	google.golang.org/genproto v0.0.0-20200128133413-58ce757ed39b // indirect
	google.golang.org/grpc v1.27.0 // indirect
	gopkg.in/dancannon/gorethink.v3 v3.0.5 // indirect
	gopkg.in/fatih/pool.v2 v2.0.0 // indirect
	gopkg.in/gorethink/gorethink.v3 v3.0.5 // indirect
	gopkg.in/square/go-jose.v1 v1.1.2
	gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v3 v3.0.0-20190924164351-c8b7dadae555
	sigs.k8s.io/kustomize/api v0.3.2
	vbom.ml/util v0.0.0-20180919145318-efcd4e0f9787 // indirect
)

exclude github.com/Azure/go-autorest v12.0.0+incompatible
