#!/bin/bash
# manuel wifi connection script
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


# get wifi interface
int=$(gwin)

# interface up if its not
int_up "$int"

echo ""
temp=$(iw "$int" link | grep SSID)
if [[ ! -z $temp ]]
then
	echo "Current Wifi Network$temp"
fi	

echo ""

echo Devices wifi interface name: "$int"
echo ""
echo Scanning wifi networks ...
sudo iw dev "$int" scan | grep SSID | \
while IFS= read i; do
    echo "	> ${i:7}"
done

echo ""

read -p "network SSID: " network
read -p "network password: " password

echo "Connecting..."
echo ""
echo "This may take a few seconds do not interrupt this operation it can break older conenection configurations."
echo ""
sudo killall wpa_supplicant
sudo ifconfig "$int" up > /dev/null
wpa_passphrase "$network" "$password" | sudo tee /etc/wpa_supplicant.conf > /dev/null
sudo wpa_supplicant -B -c /etc/wpa_supplicant.conf -i "$int" > /dev/null
sudo dhclient -r > /dev/null 
sudo dhclient -x > /dev/null
sudo dhclient "$int" > /dev/null
temp=$(iw "$int" link | grep SSID)
echo ""
echo "Now connected to: ${temp:7}"
echo ""
echo "Exiting..."
