{
  local d = (import 'doc-util/main.libsonnet'),
  '#':: d.pkg(name='v1alpha4', url='', help=''),
  cluster: (import 'cluster.libsonnet'),
  clusterClass: (import 'clusterClass.libsonnet'),
  machine: (import 'machine.libsonnet'),
  machineDeployment: (import 'machineDeployment.libsonnet'),
  machineHealthCheck: (import 'machineHealthCheck.libsonnet'),
  machinePool: (import 'machinePool.libsonnet'),
  machineSet: (import 'machineSet.libsonnet'),
}
