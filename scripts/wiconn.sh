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

function connect(){
    if [ "$#" -lt 2 ]
    then
        echo "Not enoght argument for connect"
        echo "  wiconn add <network_name> <password>"
        exit 2
    fi
    nmcli d wifi list > /dev/null
    if [ "$#" -eq 2 ]
    then
        result=$(sudo nmcli device wifi connect "$2")
        error_check "$result"
    elif [ "$#" -eq 3 ]
    then
        result=$(sudo nmcli device wifi connect "$2" password "$3")
        error_check "$result"
    fi
}

function disconnect(){
    if [ "$#" -lt 2 ]
    then
        echo "Not enoght argument for disconnect"
        echo "  wiconn dis <network_name>"
        exit 3
    fi
    result=$(sudo nmcli con down id "$2")
    error_check "$result"
}

# get wifi interface
int=$(gwin)

# interface up if its not
int_up "$int"

if [ "$#" -lt 1 ]
then
    echo "Not enoght argument"
    echo "  wiconn [command] | [networkname] [?password]"
    exit 6
fi

case $1 in
    "add")
        add_connection "$@"
        ;;
    "con"|"connect"|"conn")
        connect "$@"
        ;;
    "dis"|"disconnect")
        disconnect "$@"
        ;;
esac
