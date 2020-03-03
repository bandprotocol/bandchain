#!/usr/bin/env bash

# We wrap this whole script in a function, so that we won't execute
# until the entire script is downloaded.
# That's good because it prevents our output overlapping with curl's.
# It also means that we can't run a partially downloaded script.
setup() {
  # Display everything on stderr.
  exec 1>&2

  UNAME=$(uname)
  if [ "$UNAME" != "Linux" ] ; then
    echo "Sorry, this OS is not supported yet via this installer."
    echo "For more details on supported platforms, contact support@bandprotocol.com"
    exit 1
  fi

  # Install dependencies
  apt-get update -y
  apt-get upgrade -y
  apt-get install -y build-essential make git wget curl jq

  # Install golang
  wget https://dl.google.com/go/go1.13.7.linux-amd64.tar.gz
  tar -xvf go1.13.7.linux-amd64.tar.gz
  mv go /usr/local/

  # Export go path
  grep -q -F "export GOPATH=\$HOME/go" $HOME/.bashrc || echo "export GOPATH=\$HOME/go" >> $HOME/.bashrc
  grep -q -F "export GOROOT=/usr/local/go" $HOME/.bashrc || echo "export GOROOT=/usr/local/go" >> $HOME/.bashrc
  grep -q -F "export PATH=\$PATH:\$GOROOT/bin" $HOME/.bashrc || echo "export PATH=\$PATH:\$GOROOT/bin" >> $HOME/.bashrc
  grep -q -F "export PATH=\$PATH:\$GOPATH/bin" $HOME/.bashrc || echo "export PATH=\$PATH:\$GOPATH/bin" >> $HOME/.bashrc  
  source $HOME/.bashrc

  # clone repo & install
  TARGET_BRANCH=${BRANCH:-"master"}
  echo "=================== Target branch is $TARGET_BRANCH ==================="
  git clone https://github.com/bandprotocol/bandchain
  cd bandchain/chain
  git checkout $TARGET_BRANCH
  make install
}

setup
