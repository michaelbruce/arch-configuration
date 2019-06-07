# sysconf

A utility for installing and configuring Arch Linux

(uses and better README to come as tool is written, currently translating a few bash scripts into Go.)

# Requirements

- Install a new arch distribution (practice on my laptop)
  - on efi/mbr
  - on ssd/nvm
- Update install to latest arch live distro
- Potentially have your own arch live distro?
- Update packages (and remove/warn other packages)
- Store config (and warn/clean unwanted config)
- Include different configurations (e.g desktop uses nvidia)
- Allow all steps to be performed individually
  - include a post install collection of steps that can be rerun to keep
    packages up to date/install clean etc

# TODO

- gimp settings are not saved currently (I don't do any configuration)
- thunar settings not saved

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
