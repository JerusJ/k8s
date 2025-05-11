{
  local d = (import 'doc-util/main.libsonnet'),
  '#':: d.pkg(name='v1alpha3', url='', help=''),
  cluster: (import 'cluster.libsonnet'),
  machine: (import 'machine.libsonnet'),
  machineDeployment: (import 'machineDeployment.libsonnet'),
  machineHealthCheck: (import 'machineHealthCheck.libsonnet'),
  machinePool: (import 'machinePool.libsonnet'),
  machineSet: (import 'machineSet.libsonnet'),
}
