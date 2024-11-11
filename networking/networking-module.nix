{ config,pkgs, ...}:

{
    networking = {
        interfaces = {
            wlo0 = {
                ipv4 = {
                    addresses = [{
                        address = "192.168.50.1";
                        prefixLength = 24;
                    }];
                    gateway = "192.168.50.1";
                    useDHCP = false;
                };
            };
            eno1 = {
                    ipv4 = {
                        address = [{
                            address = "192.168.50.2";
                            prefixLength = 24;
                        }];
                        gateway = "192.168.50.2";
                        useDHCP = false;
                    };
                };
        };
        firewall = {
            extraCommands = ''
                iptables -t nat -A POSTROUTING -o enp0s20u2 -j MASQUERADE
                iptables -A FORWARD -i enp0s20u2 -o wlan0 -m state --state RELATED,ESTABLISHED -j ACCEPT
                iptables -A FORWARD -i wlan0 -o enp0s20u2 -j ACCEPT
                iptables -A FORWARD -i enp0s20u2 -o eth0 -m state --state RELATED,ESTABLISHED -j ACCEPT
                iptables -A FORWARD -i eth0 -o enp0s20u2 -j ACCEPT
            '';
        };
    };
}