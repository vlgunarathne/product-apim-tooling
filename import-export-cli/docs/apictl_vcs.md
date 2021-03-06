## apictl vcs

Checks status and deploys projects

### Synopsis

Checks status and deploys projects to the specified environment. In order to 
use this command, 'git' must be installed in the system.'

```
apictl vcs [flags]
```

### Examples

```
apictl status  -e dev
apictl deploy -e dev
```

### Options

```
  -h, --help   help for vcs
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications
* [apictl vcs deploy](apictl_vcs_deploy.md)	 - Deploys projects to the specified environment
* [apictl vcs status](apictl_vcs_status.md)	 - Shows the list of projects that are ready to deploy

