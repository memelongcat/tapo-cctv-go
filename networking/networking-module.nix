{ config,pkgs, ...}:

{
	boot.kernel.sysctl = { "net.ipv4.conf.all.forwarding" = true; };
	networking = {
		interfaces = {
			wlo1 = {
				useDHCP = false;
				ipv4.addresses = [{
					address = "192.168.50.1";
					prefixLength = 24;
				}];
			};
			eno1 = {
				useDHCP = false;
				ipv4.addresses = [{
					address = "192.168.60.1";
					prefixLength = 24;
				}];
			};
			enp0s20u2 ={
				useDHCP = true;
			};
		};
#		nat = {
#			enable = true;
#			externalInterface = "enp0s20u2";
#			internalInterfaces = [ "wlo1" "eno1" ];
#		};
#		defaultGateway = {
#			address = "192.168.50.1";
#			interface = "enp0s20u2";
#		};
		firewall = {
			enable = true;
			allowedTCPPorts = [ 53 80 139 443 445 ];
			allowedUDPPorts = [ 53 67 68 ];
			allowPing = true;

			extraCommands = ''
				iptables -t nat -A POSTROUTING -o enp0s20u2 -j MASQUERADE
				iptables -A FORWARD -i wlo1 -o eno1 -j ACCEPT
				iptables -A FORWARD -i eno1 -o wlo1 -j ACCEPT
#				iptables -A FORWARD -i eno1 -j LOG --log-prefix "FORWARD eno1: "
#				iptables -A FORWARD -i wlo1 -j LOG --log-prefix "FORWARD wlo1: "
#				iptables -A FORWARD -o enp0s20u2 -j LOG --log-prefix "FORWARD to enp0s20u2: "
			'';
		};
		networkmanager.unmanaged = [ "wlo1" ];
	};
#	security.wrappers = {
#		rootlesskit = {
#		owner = "root";
#		group = "root";
#		capabilities = "cap_net_bind_service+ep";
#		source = "${pkgs.rootlesskit}/bin/rootlesskit";
#		};
#	};
	environment.systemPackages = with pkgs; [
    	libcap
  	];
	system.activationScripts.setDockerCap = ''
    	setcap cap_net_bind_service=+ep ${pkgs.docker}/bin/dockerd
  	'';
}
