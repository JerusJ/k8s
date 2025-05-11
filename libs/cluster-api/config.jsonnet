local config = import 'jsonnet/config.jsonnet';

local versions = [
  'v1.10.1',
  'v1.0.2',
];

config.new(
  name='cluster-api',
  specs=[
    {
      output: v,
      openapi: 'http://localhost:8001/openapi/v2',
      prefix: '^io\\.x-k8s\\.cluster\\..*',
      crds: [
        'https://github.com/kubernetes-sigs/cluster-api/releases/download/%s/core-components.yaml' % v,
      ],
      localName: 'cluster-api',
    }
    for v in versions
  ]
)
