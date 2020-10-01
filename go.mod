module github.com/qlik-oss/sense-installer

go 1.14

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20200319173657-742aab907b54
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a

	k8s.io/apimachinery => k8s.io/apimachinery v0.17.0
	k8s.io/client-go => k8s.io/client-go v0.17.0
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191219154910-1528d4eea6dd
	sigs.k8s.io/kustomize/api => github.com/qlik-oss/kustomize/api v0.6.3-0.20201001044635-f23c10b755f6
)

require (
	cloud.google.com/go v0.52.0 // indirect
	cloud.google.com/go/storage v1.5.0 // indirect
	github.com/Masterminds/semver/v3 v3.1.0
	github.com/Shopify/ejson v1.2.1
	github.com/aws/aws-sdk-go v1.28.9 // indirect
	github.com/bugsnag/bugsnag-go v1.5.3 // indirect
	github.com/containers/image/v5 v5.1.0
	github.com/go-git/go-git/v5 v5.1.0
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobuffalo/logger v1.0.3 // indirect
	github.com/gobuffalo/packd v1.0.0 // indirect
	github.com/gobuffalo/packr/v2 v2.7.1
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/logrusorgru/aurora v0.0.0-20200102142835-e9ef32dff381
	github.com/mattn/go-colorable v0.1.4
	github.com/mattn/go-tty v0.0.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/otiai10/copy v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/qlik-oss/k-apis v0.1.16
	github.com/rancher/k3d/v3 v3.0.2
	github.com/robfig/cron/v3 v3.0.1
	github.com/rogpeppe/go-internal v1.5.2 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.0.1-0.20200629195214-2c5a0d300f8b
	github.com/spf13/viper v1.6.1
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4 // indirect
	golang.org/x/exp v0.0.0-20200119233911-0405dc783f0a // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
	golang.org/x/tools v0.0.0-20200312194400-c312e98713c2 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.17.2
	k8s.io/apiextensions-apiserver v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kubectl v0.17.2
	sigs.k8s.io/kustomize/api v0.3.2
	sigs.k8s.io/yaml v1.2.0
)

exclude github.com/Azure/go-autorest v12.0.0+incompatible
