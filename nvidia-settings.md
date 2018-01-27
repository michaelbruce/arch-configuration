# Nvidia display settings

Ensure force composition pipeline is set to avoid video tearing
This is best done by using the Save to X Configuration File button in the
nvidia-settings GUI interface to ~, then copying it to /etc/X11/xorg.conf

Use xcompmgr (compton is slow as hell.)
