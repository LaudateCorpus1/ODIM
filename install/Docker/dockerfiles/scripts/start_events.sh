#!/bin/bash
# (C) Copyright [2020] Hewlett Packard Enterprise Development LP
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.

declare PID=0
declare OWN_PID=$$

sigterm_handler()
{
        if [[ $PID -ne 0 ]]; then
                sleep 1
                kill -9 $PID
                wait "$PID" 2>/dev/null
        fi
        exit 0
}

# create a signal trap
create_signal_trap()
{
        trap 'echo "[$(date)] -- INFO  -- SIGTERM received for event, initiating shut down"; sigterm_handler' SIGTERM
}

# keep the script running till SIGTERM is received
run_forever()
{
        wait
}

start_event()
{
        registry_address="consul:8500"
        if [[ ${HA_ENABLED,,} == true ]]; then
                registry_address="[consul1:8500,consul2:8500,consul3:8500]"
        fi

	export CONFIG_FILE_PATH=/etc/odimra_config/odimra_config.json
	nohup /bin/svc-events --registry=consul --registry_address=${registry_address} --server_address=event:45103 --client_request_timeout=`expr $(cat $CONFIG_FILE_PATH | grep SouthBoundRequestTimeoutInSecs | cut -d : -f2 | cut -d , -f1 | tr -d " ")`s >> /var/log/odimra_logs/event.log 2>&1 &
	PID=$!
	sleep 3

	nohup /bin/add-hosts -file /tmp/host.append >> /var/log/odimra_logs/event-add-hosts.log 2>&1 &
}

monitor_process()
{
	while true; do
		pid=$(pgrep -fc svc-events 2> /dev/null)
		if [[ $? -ne 0 ]] || [[ $pid -gt 1 ]]; then
			echo "svc-events has exited"
			kill -15 ${OWN_PID}
			exit 1
		fi
		sleep 5
	done &
}

##############################################
###############  MAIN  #######################
##############################################

start_event

create_signal_trap

monitor_process

run_forever

exit 0