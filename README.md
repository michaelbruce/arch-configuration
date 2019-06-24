# sysconf

A utility for installing and configuring Arch Linux

(uses and better README to come as tool is written, currently translating a few bash scripts into Go.)

## Operations

### Boot
Preparing a live USB Arch system, including an X environment and sysconf.  
Simply invoke the command:

`sysconf -live /dev/sdx` 

Where /dev/sdx is an unmounted USB volume, after the setup time (about 10 minutes) you should have a bootable live version of Arch.

### Install
Once you have booted into your live Arch setup, you can install a more permanent Arch Linux on one of your host machine's storage volumes with:

`sysconf -install /dev/sdy`

Where /dev/sdy is a storage volume (e.g an SSD) that you want to reformat and install Arch Linux on.

### Update
- Keeping the new OS synced and up to date

# Requirements

- [ ] Install a new arch distribution (practice on my laptop)
  - [ ] on efi/mbr
  - [ ] on ssd/nvm
- [ ] Update install to latest arch live distro
- [x] Potentially have your own arch live distro?
- [x] Update packages (and remove/warn other packages)
- [x] Store config (and warn/clean unwanted config)
- [ ] Include different configurations (e.g desktop uses nvidia)
- [x] Allow all steps to be performed individually
  - [x] include a post install collection of steps that can be rerun to keep
    packages up to date/install clean etc

# Design Decisions

- UEFI is assumed as chipsets with BIOS are being phased out.

# TODO

- set localtime to London (causes all kinds of bad connection issues) - ln -s /usr/share/zoneinfo/Europe/London /etc/localtime
- sync time - sudo systemctl enable ntpd
- installing AUR packages can be done with git/makepkg/pacman -U (updates with pull potentially)
- gimp settings are not saved currently (I don't do any configuration)
- thunar settings not saved
- sudo cp $(which vim) /usr/local/bin/vi
- lspci | grep NVIDIA - install nvidia (nvidia currently default)
- replacements for:
    - evince (has lots of deps)
    - thunar
    - slock
- install from AUR:
    - yay (installs others potentially)
    - mongo
    - minecraft
    - lierolibre
    - zoom
    - ttf-ms-fonts
- with thunar/pavucontrol etc should you just use a window manager?
- open pdfs in firefox (save to user.js)
- steam?
- efibootmgr/efivar for UEFI machines
- font-manager
- wkhtmltopdf-static
- iptables
- full disk encryption (LUKS)
- -update won't work if pactree isn't present (install pacman-contrib in this
  case)
- Bluetooth requires that AutoEnable=true be set under [Policy] in
  /etc/bluetooth/main.conf to power on at boot.

## OLD README CONTENT (come back to this)

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

### TODO

ssh-agent will fail to start if ssh socket already exists (or is possibly bound
to a port that gets taken elsewhere) - either way removing the file before
starting the service would be great.. maybe we should just /tmp it?

### Bluetooth

sudo systemctl enable bluetooth # include bluetooth at startup
> proceed to pair devices with bluetoothctl - pair <ID>
> proceed to trust devices with bluetoothctl - trust <ID>

Ensure AutoEnable is set to true in `/etc/bluetooth/main.conf`
[Policy]
AutoEnable=true

> Current issue - we can connect to the bluetooth devices using the connect
> command but autoconnect on startup is not happening.

There was a problem with bluetooth agents (unsure of this term)

When running bluetoothctl I was given a: `Agent registered` prompt. Upon typing
default-agent.. I was able to register with my mouse/keyboard as normally
expected.

### Notes

GIMP is still on GTK2 and so will use the config files
