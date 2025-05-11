{
  local d = (import 'doc-util/main.libsonnet'),
  '#':: d.pkg(name='cluster-api', url='github.com/JerusJ/k8s/cluster-api-libsonnet/v1.0.2/main.libsonnet', help=''),
  addons:: (import '_gen/addons/main.libsonnet'),
  cluster:: (import '_gen/cluster/main.libsonnet'),
}
