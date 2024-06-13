# This script updates the dependencies in the Helm chart file.
#
# Usage:
#   python update_chart.py <filename>


import subprocess
import sys

def get_latest_tag(repo):
    cmds = [
        'git',
        'ls-remote',
        '--tags',
        '--refs',
        '--sort=v:refname',
        'https://github.com/llm-operator/%s.git' % repo
    ]
    output = subprocess.check_output(cmds).decode('utf-8')
    tags = output.split('\n')
    tags.reverse()
    latest_tagline = ''
    for tag in tags:
        if tag:
            latest_tagline = tag
            break
    latest_tag = latest_tagline.split('\t')[1].split('/')[-1]
    return latest_tag.strip('v')


def update_chart(filename):
    repos = [
        'cluster-manager',
        'file-manager',
        'inference-manager',
        'job-manager',
        'model-manager',
        'user-manager',
        'rbac-manager',
        'session-manager',
        'vector-store-manager',
        ]
    tags = {}
    for repo in repos:
        tag = get_latest_tag(repo)
        tags[repo] = tag

    deps = {
        'cluster-manager-server': tags['cluster-manager'],
        'dex-server': tags['rbac-manager'],
        'file-manager-server': tags['file-manager'],
        'inference-manager-engine': tags['inference-manager'],
        'inference-manager-server': tags['inference-manager'],
        'job-manager-dispatcher': tags['job-manager'],
        'job-manager-server': tags['job-manager'],
        'model-manager-loader': tags['model-manager'],
        'model-manager-server': tags['model-manager'],
        'rbac-server': tags['rbac-manager'],
        'session-manager-agent': tags['session-manager'],
        'session-manager-server': tags['session-manager'],
        'user-manager-server': tags['user-manager'],
        'vector-store-manager-server': tags['vector-store-manager'],
    }

    chart = """apiVersion: v2
name: llm-operator
description: A Helm chart for LLM Operator
type: application
version: 0.1.0
appVersion: 0.1.0
dependencies:
"""
    for dep, tag in deps.items():
        chart += """- name: %(dep)s
  version: %(tag)s
  repository: "oci://public.ecr.aws/v8n3t7y5/llm-operator-charts"
""" % {'dep': dep, 'tag': tag}
    # Write the chart to the file
    with open(filename, 'w') as f:
        f.write(chart)


def main():
    # Get a path name from the commandline arguments
    if len(sys.argv) < 2:
        print("Usage: python %s <filename>" % sys.argv[0])
        sys.exit(1)
    filename = sys.argv[1]
    update_chart(filename)

if __name__ == "__main__":
    main()
