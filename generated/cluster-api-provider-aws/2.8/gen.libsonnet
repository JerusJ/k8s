{
  local d = (import 'doc-util/main.libsonnet'),
  '#':: d.pkg(name='cluster-api-provider-aws', url='github.com/JerusJ/k8s/cluster-api-provider-aws-libsonnet/2.8/main.libsonnet', help=''),
  bootstrap:: (import '_gen/bootstrap/main.libsonnet'),
  controlplane:: (import '_gen/controlplane/main.libsonnet'),
  infrastructure:: (import '_gen/infrastructure/main.libsonnet'),
}
