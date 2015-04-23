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
  pid=`ps -ef | grep "$USER_NAME" | grep "./zebra -config=../config/zebra_config.xml" | grep -v grep | grep -v $0 | awk '{print $2}'`
  if [ "$pid" != "" ]
  then
    printf "${_GREEN}ZebraDB RUNNING${_NORMAL}\n"
    svr_pids["zebra"]=$pid
  else
    printf "${_RED}ZebraDB STOP${_NORMAL}\n"
  fi
}

kill_service()
{
  check_service_exist
  if [ ${#svr_pids["zebra"]} -gt 0 ]
  then
    printf "\t${_YELLOW}killing ZebraDB server ${_CYAN}${_YELLOW}with pid=${svr_pids["zebra"]}...${_NORMAL}\t\t"
    kill ${svr_pids["zebra"]} || exit 0
    printf "\t[${_RED} KILLED${_NORMAL}]\n"
  fi

  sleep 1
}

if [ $# -lt 1 ];
then
  # check the un-started server and start it
  kill_service
  check_service_exist
  while [ ${#svr_pids["zebra"]} -gt 0 ] 
  do
    sleep 1
    check_service_exist 
  done
  exit 0
fi
