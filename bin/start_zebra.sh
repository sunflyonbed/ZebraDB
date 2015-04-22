#!/bin/bash

_NORMAL="\033[0m"
_YELLOW="\033[0;33m"
_CYAN="\033[1;36m"
_GREEN="\033[1;32m"
_RED="\033[1;31m"
_PERPLE="\033[0;35m"

USER_NAME=`whoami`

svr_pids=()

check_service_exist()
{
  svr_pids=()
  pid=`ps -ef | grep "$USER_NAME" | grep "./zebra" | grep -v grep | grep -v $0 | awk '{print $2}'`
  if [ "$pid" != "" ]
  then
    printf "${_GREEN}HLOG RUNNING${_NORMAL}\n"
    svr_pids["zebra"]=$pid
  else
    printf "${_RED}HLOG STOP${_NORMAL}\n"
  fi
}

start_service() 
{
  check_service_exist
  if [ ${#svr_pids["zebra"]} -lt 1 ]
  then
    printf "\t${_YELLOW} executing \"./zebra ${_NORMAL}\t\t\n"
    ## start server command
    nohup ./zebra 1>> zebra_err.log 2>&1 & > /dev/null
  fi

  sleep 1

  check_service_exist
  if [ ${#svr_pids["zebra"]} -lt 1 ]
  then
    printf "${_YELLOW} zebra server start fail${_RED}\n"
  fi
}

if [ $# -lt 1 ];
then
  # check the un-started server and start it
  start_service
  exit 0
fi
