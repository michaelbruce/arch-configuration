# Arch Configuration

A set of files that install and configure arch linux

### Setup

- Launch the Arch Linux boot iso via USB/CD/DVD
- Connect to the internet
- Download this repository as a tarball `wget https://github.com/mikepjb/arch-configuration/tarball/master -O - | tar xz`
- run the `./install` script

### First Boot

When starting up the freshing installed operating system you should:
- Login as root (no password required yet)
- Set the root password with `passwd`
- Set the default user password (named hades) with `passwd hades`
- Use `startx` to boot into bspwm, a window manager GUI

### Use

pacaur can get packages from AUR

### Arch tips

remove a package and others that depend on them with pacman -Rc

### Apple Magic mouse

echo 'modprobe hid_magicmouse scroll_acceleration=1 scroll_speed=55' > /etc/modprobe.d/hid_magicmouse.conf

### TODO: enable ssh-agent service by default...
