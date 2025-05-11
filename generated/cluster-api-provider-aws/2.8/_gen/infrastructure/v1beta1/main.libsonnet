{
  local d = (import 'doc-util/main.libsonnet'),
  '#':: d.pkg(name='v1beta1', url='', help=''),
  awsCluster: (import 'awsCluster.libsonnet'),
  awsClusterControllerIdentity: (import 'awsClusterControllerIdentity.libsonnet'),
  awsClusterRoleIdentity: (import 'awsClusterRoleIdentity.libsonnet'),
  awsClusterStaticIdentity: (import 'awsClusterStaticIdentity.libsonnet'),
  awsClusterTemplate: (import 'awsClusterTemplate.libsonnet'),
  awsFargateProfile: (import 'awsFargateProfile.libsonnet'),
  awsMachine: (import 'awsMachine.libsonnet'),
  awsMachinePool: (import 'awsMachinePool.libsonnet'),
  awsMachineTemplate: (import 'awsMachineTemplate.libsonnet'),
  awsManagedMachinePool: (import 'awsManagedMachinePool.libsonnet'),
}
