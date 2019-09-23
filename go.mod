module github.com/adrienaury/owl

go 1.13

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/MakeNowJust/heredoc v0.0.0-20170808103936-bb23615498cd
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	github.com/docker/docker-credential-helpers v0.6.3
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/jsimonetti/pwscheme v0.0.0-20160922125227-76804708ecad
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.1
	github.com/pkg/errors v0.8.1 // indirect
	github.com/russross/blackfriday v1.5.2
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0 // indirect
	golang.org/x/sys v0.0.0-20190616124812-15dcb6c0061f // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/ldap.v3 v3.0.3
	gopkg.in/yaml.v3 v3.0.0-20190905181640-827449938966
	gotest.tools v2.2.0+incompatible // indirect
)

replace k8s.io/kubernetes/pkg/kubectl/util/term => k8s.io/kubectl/pkg/util/term v0.0.0-20190918164019-21692a0861df
