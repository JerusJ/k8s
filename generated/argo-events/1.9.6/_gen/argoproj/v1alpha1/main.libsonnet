{
  local d = (import 'doc-util/main.libsonnet'),
  '#':: d.pkg(name='v1alpha1', url='', help=''),
  eventBus: (import 'eventBus.libsonnet'),
  eventSource: (import 'eventSource.libsonnet'),
  sensor: (import 'sensor.libsonnet'),
}
