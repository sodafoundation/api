#!/bin/bash -xe

# Local version to install bindep packages
# Suitable for use for development

function is_fedora {
    [ -f /usr/bin/yum ] && cat /etc/*release | grep -q -e "Fedora"
}

PACKAGES=""
if ! which virtualenv; then
    PACKAGES="$PACKAGES virtualenv"
fi
if ! which make; then
    PACKAGES="$PACKAGES make"
fi
if [[ -n $PACKAGES ]]; then
    sudo apt-get -q --assume-yes install virtualenv
fi

# Check for bindep
if ! which bindep; then
    make bindep
fi

PACKAGES=$(make bindep || true)

# inspired from project-config install-distro-packages.sh
if apt-get -v >/dev/null 2>&1 ; then
    sudo apt-get -qq update
    sudo PATH=/usr/sbin:/sbin:$PATH DEBIAN_FRONTEND=noninteractive \
        apt-get -q --option "Dpkg::Options::=--force-confold" \
        --assume-yes install $PACKAGES
elif emerge --version >/dev/null 2>&1 ; then
    sudo emerge -uDNq --jobs=4 @world
    sudo PATH=/usr/sbin:/sbin:$PATH emerge -q --jobs=4 $PACKAGES
else
    is_fedora && YUM=dnf || YUM=yum
    sudo PATH=/usr/sbin:/sbin:$PATH $YUM install -y $PACKAGES
fi
