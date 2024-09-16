# Change log
## [vx.x.xx- 2024-07-18]
### Added
Added error page with error string Error.err
### Changed
Moved part code from crbcreate to parse, now it parse and create cluster role binding

## [vx.x.xx- 2024-07-24]
### Added
Button cluster role description
New part of code for getcrdesc.go with parse code and description 
### Changed
Fixed return to crbcman after error message 
now when you try to create crb which one already exist press on button "return on previous page" redirect to crbcmain 

## [vx.x.xx- 2024-07-25]
### Added

### Changed
Pretty output for service account and cluster role bindings
Added new code part, now service account collect via loop and append to slice, previous output was unreadable
used string builder
output with yaml style

## [vx.x.xx- 2024-07-26]
### Added
Added new part of code with global var clientset "globalvar.Clientset"
dir globalvar package globalvar

### Changed
Buttons size changed, field select size changed
Vars below, deleted from all files
config, _    = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
Clientset, _ = kubernetes.NewForConfig(config)
Fixed incorrect output for cluster role, before fix apiGroup repeat every iteration because map not cleared
Resolution below:
// add data to map
m0["apiGroups"] = sl0
// add data to main slice
outSlice = append(outSlice, m0)
// clear map
m0 = make(map[string][]string)
// clear slice
sl0 = nil

## [vx.x.xx- 2024-07-29]
### Added
Validating webhook 
Check health path /health to check application 
### Changed

## [vx.x.xx- 2024-07-30]
### Added
netstat added to alpine image
### Changed
Changed _gen_certs.sh 
For generate certs use command ./_gen_certs.sh  <namespacename> <servicename>
manager-ns-webhook-tls add namespace

## [v0.1.5 2024-09-13]
### Added
Created ApplicationSet AppProject for Argo-CD
Deleted secrets and added annotation cert-manager.io/cluster-issuer: cluster-ca for generate secrets
### Changed
-----







