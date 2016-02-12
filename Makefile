# Makefile for preparing and running gonit binaries
KERNEL=/boot/vmlinuz-4.4.0-2-generic
INITRD=/boot/initrd.img-4.4.0-2-generic
APPEND="console=ttyS0 root=/dev/sda1 rw init=/init/gonit quiet"
DISK=~/gonit-disk.img
DISK_MOUNT=/tmp/gonit-disk
BUILD_OUT=/tmp/gonit-build
VM_MEMORY=1G

# Colors
REDCOLOR="\033[31m"
GREENCOLOR="\033[32m"
ORANGECOLOR="\033[33m"
BLUECOLOR="\033[34m"
BINCOLOR="\033[35m"
CYANCOLOR="\033[36m"
ENDCOLOR="\033[0m"



prepare_output_folders:
	mkdir -p ${BUILD_OUT}/init
	mkdir -p ${BUILD_OUT}/services

init: prepare_output_folders
	@printf '%b %b\n' $(REDCOLOR)COMPILE$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR)
	go build -o ${BUILD_OUT}/init/gonit github.com/mustafaakin/gonit/cmd/init

services: prepare_output_folders
	@printf '%b %b\n' $(REDCOLOR)COMPILE$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR)
	CGO_ENABLED=0 go build -o ${BUILD_OUT}/services/networking github.com/mustafaakin/gonit/cmd/networking
	go build -o ${BUILD_OUT}/services/terminals github.com/mustafaakin/gonit/cmd/terminals
	
create_disk:
	@printf '%b %b\n' $(GREENCOLOR)DISK$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR)
	modprobe nbd
	qemu-img create -f qcow2 ${DISK} 1G
	qemu-nbd -c /dev/nbd0 ${DISK}
	@echo "Please format the disk image with one partition"
	fdisk /dev/nbd0 
	mkfs -t ext4 /dev/nbd0p1
	mkdir -p ${DISK_MOUNT}
	mount /dev/nbd0p1 ${DISK_MOUNT}    
	# Prepare special folders
	mkdir -p ${DISK_MOUNT}/dev
	mkdir -p ${DISK_MOUNT}/proc
	mkdir -p ${DISK_MOUNT}/sys
	mkdir -p ${DISK_MOUNT}/run


mount_disk:
	@printf '%b %b\n' $(GREENCOLOR)DISK$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR)
	modprobe nbd
	qemu-nbd -c /dev/nbd0 ${DISK}
	mkdir -p ${DISK_MOUNT}	
	mount /dev/nbd0p1 ${DISK_MOUNT}   
	
umount_disk:
	@printf '%b %b\n' $(REDCOLOR)COMPILE$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR)
	qemu-nbd -d /dev/nbd0	
	umount ${DISK_MOUNT}  
	
package: compile
	@printf '%b %b\n' $(BLUECOLOR)PACKAGE$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR)
	# Folders, for better readability
	mkdir -p ${DISK_MOUNT}/init
	mkdir -p ${DISK_MOUNT}/init/services
	
	# Copy required bins
	cp ${BUILD_OUT}/init/gonit ${DISK_MOUNT}/init/gonit
	cp ${BUILD_OUT}/services/networking ${DISK_MOUNT}/init/services/networking
	cp ${BUILD_OUT}/services/terminals  ${DISK_MOUNT}/init/services/terminals
	
	# Copy configs
	cp config.yml ${DISK_MOUNT}/config.yml

compile: init services
		
clean:
	rm -Rf ${BUILD_OUT} 

sync:
	sync
	blockdev --flushbufs /dev/nbd0p1
	blockdev --flushbufs /dev/nbd0
	sync
	sleep 2 
					 
start_vm: clean package sync
	@printf '%b %b\n' $(ORANGECOLOR)RUN$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR)
	# TODO: The following sync and flushbufs needs to be changed, we will just disable cache on our image but it seems problematic on rbd?
	kvm -m ${VM_MEMORY} -nographic -kernel ${KERNEL} -initrd ${INITRD} -append ${APPEND} -hda ${DISK}
	
kill_vm:
	# Beware that kills all VMs, TODO: just my running vm.
	killall -9 qemu-system-x86_64
