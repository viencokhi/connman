#! /bin/bash
# author eco/ecoshub

# get wifi interface
function gwin () {
    ip link show | \
    while IFS= read int; do
        if [ "${int:3:1}" == "w" ]
        then
            search="w"
            prefix=${int%%$search*}
            int=${int:${#prefix}}
            search=":"
            prefix=${int%%$search*}
            int=${int:0:${#prefix}}
            echo $int
        fi
    done
    exit 1
}

function int_up(){
    result=$(ip link show "$1" | grep "UP")
    if [ -z $reuslt]
    then
        sudo ip link set "$1" up
    fi
}

function error_check(){
    success=$(echo "$1" | grep "Error\|error")
    if [ ! -z "$success" ]
    then
        echo "$1"
        exit 4
    fi
}

function add_network(){
    if [ $# -lt 2 ]; then
        echo "not enoght argument to add a hotspot network"
        exit 2
    fi
    # add network
    result=$(nmcli con add type wifi ifname "$1" con-name "$3" autoconnect yes ssid "$3")
    error_check "$result"
    # set comm node
    nmcli con modify "$3" 802-11-wireless.mode ap 802-11-wireless.band bg ipv4.method shared

    if [ $# -eq 4 ]
    then
        # set comm pass
        nmcli con modify "$3" wifi-sec.key-mgmt wpa-psk
        nmcli con modify "$3" wifi-sec.psk "$4"
    fi
}

function network_manage(){
    if [ $# -lt 2 ]; then
        echo "not enoght argument to $2 a network"
        exit 4
    fi

    case $2 in
        "up"|"down"|"delete"|"del")
            result=$(nmcli con $2 $3)
            ;;
        *)
            echo "not a valid command. use [up down delete]"
            ;;
    esac
    error_check "$result"
}

# get wifi interface
int=$(gwin)

# interface up if its not
int_up "$int"

if [ "$#" -lt 2 ]
then
    echo "Not enoght argument"
    echo "  spotligh [command] [networkname] [?password]"
    exit 6
fi

case $1 in
    "add")
        add_network "$int" "$@"
        ;;
    "add-con")
        add_network "$int" "$@"
        network_up "$@"
        ;;
    "up"|"down"|"delete"|"del")
        network_manage "$int" "$@"
        ;;
    *)
        echo "not a valid command."
        echo "  [command] : add, add-con, up, down, delete"
        ;;
esac
