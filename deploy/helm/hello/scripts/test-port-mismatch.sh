#!/usr/bin/env bash
set -euo pipefail

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

helm template hello deploy/helm/hello > "$TMPDIR/rendered.yaml"

python3 -c "
import yaml, sys

with open('$TMPDIR/rendered.yaml') as f:
    docs = list(yaml.safe_load_all(f))

errors = 0
svc_port = None
svc_target = None
container_port = None
named_port_val = None
live_path = None
live_port = None
ready_path = None
ready_port = None

for doc in docs:
    if not doc:
        continue
    kind = doc.get('kind', '')

    if kind == 'Service':
        for p in doc.get('spec', {}).get('ports', []):
            svc_port = p.get('port')
            svc_target = p.get('targetPort')

    if kind != 'Deployment':
        continue

    containers = doc.get('spec', {}).get('template', {}).get('spec', {}).get('containers', [])
    for c in containers:
        for cp in c.get('ports', []):
            if cp.get('name') == 'http':
                named_port_val = cp.get('containerPort')
                container_port = cp.get('containerPort')
        lp = c.get('livenessProbe', {}).get('httpGet', {})
        live_path = lp.get('path')
        live_port = lp.get('port')
        rp = c.get('readinessProbe', {}).get('httpGet', {})
        ready_path = rp.get('path')
        ready_port = rp.get('port')

checks = [
    (svc_port == 80, f'Service port is {svc_port}, expected 80'),
    (container_port == 8080, f'containerPort is {container_port}, expected 8080'),
    (named_port_val == 8080, f'Named port http resolves to {named_port_val}, expected 8080'),
    (svc_target == 'http', f'Service targetPort is {svc_target}, expected http'),
    (live_port == 'http' and named_port_val == 8080, f'Liveness probe port resolution issue'),
    (ready_port == 'http' and named_port_val == 8080, f'Readiness probe port resolution issue'),
    (live_path == '/healthz', f'Liveness probe path is {live_path}, expected /healthz'),
    (ready_path == '/readyz', f'Readiness probe path is {ready_path}, expected /readyz'),
]

for ok, msg in checks:
    if not ok:
        print(f'FAIL: {msg}')
        errors += 1

if errors > 0:
    print(f'FAIL: {errors} check(s) failed')
    sys.exit(1)

print('PASS: All ports match expected values')
print('  - Service exposes port 80')
print('  - Container listens on port 8080 (named http)')
print('  - Service targetPort http forwards to 8080')
print('  - Probes target /healthz and /readyz on port 8080')
sys.exit(0)
"
