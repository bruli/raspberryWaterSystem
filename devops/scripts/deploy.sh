#!/usr/bin/env bash

make decryptVault
ansible-playbook -i devops/ansible/inventories/production/hosts devops/ansible/deploy.yml
make encryptVault
