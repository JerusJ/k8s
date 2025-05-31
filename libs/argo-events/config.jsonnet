local config = import 'jsonnet/config.jsonnet';

local versions = [
  '1.9.6',
];

local manifests = [
  'argoproj.io_eventbus.yaml',
  'argoproj.io_eventsources.yaml',
  'argoproj.io_sensors.yaml',
];

config.new(
  name='argo-cd',
  specs=[
    {
      output: version,
      prefix: '^io\\.argoproj\\..*',
      localName: 'argo_cd',
      crds: [
        'https://raw.githubusercontent.com/argoproj/argo-cd/v%s/manifests/base/crds/%s' %
        [version, manifest]
        for manifest in manifests
      ],
    }
    for version in versions
  ]
)
