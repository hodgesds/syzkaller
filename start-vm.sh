#!/bin/sh
KERNEL=/usr/src/linux-4.19-rc2
qemu-system-x86_64 \
	-kernel $KERNEL/arch/x86_64/boot/bzImage \
	-append "console=ttyS0 root=/dev/sda debug earlyprintk=serial slub_debug=QUZ" \
	-net user,hostfwd=tcp::10021-:22 -net nic \
	-drive format=raw,file=./stretch.img \
	-enable-kvm \
	-nographic \
	-m 2G \
	-smp 2 \
	#-pidfile vm.pid \2>&1 | tee vm.log
	#\ #-hda $IMAGE/stretch.img \
