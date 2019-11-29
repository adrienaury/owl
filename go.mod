module github.com/adrienaury/owl

go 1.13

require (
	github.com/docker/docker-credential-helpers v0.6.3
	github.com/kr/pretty v0.1.0 // indirect
	github.com/spf13/cobra v0.0.5
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/ldap.v3 v3.1.0
	gopkg.in/yaml.v3 v3.0.0-20190905181640-827449938966
)

replace k8s.io/kubernetes/pkg/kubectl/util/term => k8s.io/kubectl/pkg/util/term v0.0.0-20190918164019-21692a0861df
