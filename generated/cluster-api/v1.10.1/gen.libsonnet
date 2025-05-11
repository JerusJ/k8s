{
  local d = (import 'doc-util/main.libsonnet'),
  '#':: d.pkg(name='cluster-api', url='github.com/JerusJ/k8s/cluster-api-libsonnet/v1.10.1/main.libsonnet', help=''),
  addons:: (import '_gen/addons/main.libsonnet'),
  cluster:: (import '_gen/cluster/main.libsonnet'),
  ipam:: (import '_gen/ipam/main.libsonnet'),
  runtime:: (import '_gen/runtime/main.libsonnet'),
}
