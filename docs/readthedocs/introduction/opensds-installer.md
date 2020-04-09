# OpenSDS Installer
For detailed information about this project, please refer to the
[repo](https://github.com/sodafoundation/api-installer).

## Introduction
This project is designed for locating the code for installing all required
components to set up a cluster, including controller and nbp plugins. Currently
we support several install tools for diversity.

### Ansible
[Ansible](https://github.com/ansible/ansible) is a radically simple IT automation
platform that makes your applications and systems easier to deploy. OpenSDS
installer project holds all code related to `opensds-ansible` in ansible folder
for installing and configuring OpenSDS cluster through ansible tool.

### Helm
[Helm](https://github.com/kubernetes/helm) is a popular tool for managing
Kubernetes charts. Charts are packages of pre-configured Kubernetes resources.
OpenSDS installer project also holds all code related to `opensds-charts` in
charts folder for installing and configuring OpenSDS cluster through helm tool.

## Contact

* Mailing list: [opensds-tech-discuss](https://lists.opensds.io/mailman/listinfo/opensds-tech-discuss)
* slack: #[opensds](https://opensds.slack.com)
* Ideas/Bugs: [issues](https://github.com/sodafoundation/api-installer/issues)
