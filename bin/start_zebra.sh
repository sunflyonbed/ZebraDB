#!/bin/bash

_NORMAL="\033[0m"
_YELLOW="\033[0;33m"
_CYAN="\033[1;36m"
_GREEN="\033[1;32m"
_RED="\033[1;31m"
_PERPLE="\033[0;35m"

USER_NAME=`whoami`

DB_DIR="ZebraDB_PATH"

SVR="${DB_DIR}/bin/zebra"
SVR_CFG="-config=${DB_DIR}/config/zebra_config.xml"
SVR_LOG="${DB_DIR}/log/zebra_err.log"

svr_pids=()

check_service_exist()
{
  svr_pids=()
  pid=`ps -ef | grep "$USER_NAME" | grep "${SVR} ${SVR_CFG}" | grep -v grep | grep -v $0 | awk '{print $2}'`
  if [ "$pid" != "" ]
  then
    printf "${_GREEN}ZebraDB RUNNING${_NORMAL}\n"
    svr_pids["zebra"]=$pid
  else
    printf "${_RED}ZebraDB STOP${_NORMAL}\n"
  fi
}

start_service() 
{
  check_service_exist
  if [ ${#svr_pids["zebra"]} -lt 1 ]
  then
    printf "\t${_YELLOW} executing \"${SVR} ${SVR_CFG} ${_NORMAL}\t\t\n"
    ## start server command
    nohup ${SVR} ${SVR_CFG} 1>> ${SVR_LOG} 2>&1 & > /dev/null
  fi

  sleep 1

  check_service_exist
  if [ ${#svr_pids["zebra"]} -lt 1 ]
  then
    printf "${_YELLOW} ZebraDB server start fail${_RED}\n"
  fi
}

if [ $# -lt 1 ];
then
  # check the un-started server and start it
  start_service
  exit 0
fi
