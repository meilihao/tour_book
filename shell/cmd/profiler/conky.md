# conky
用于X窗口系统的系统监控器, 具有灵活配置.

```bash
$ sudo apt install conky-all
$ cp /etc/conky/conky.conf ~/.conkyrc # conky.conf默认配置
$ vim ~/.conkyrc
conky.config = {
	update_interval = 1,
	cpu_avg_samples = 2,
	net_avg_samples = 2,
	out_to_console = false,
	override_utf8_locale = true,
	double_buffer = true,
	no_buffers = true,
	text_buffer_size = 32768,
	imlib_cache_size = 0,
	own_window = true,
	own_window_type = 'normal',
	own_window_argb_visual = true,
	own_window_argb_value = 50,
	own_window_hints = 'undecorated,below,sticky,skip_taskbar,skip_pager',
	border_inner_margin = 5,
	border_outer_margin = 0,
	xinerama_head = 1,
	alignment = 'bottom_right',
	gap_x = 0,
	gap_y = 33,
	draw_shades = false,
	draw_outline = false,
	draw_borders = false,
	draw_graph_borders = false,
	use_xft = true,
	font = 'Ubuntu Mono:size=12',
	xftalpha = 0.8,
	uppercase = false,
	default_color = 'white',
	own_window_colour = '#000000',
	minimum_width = 300, minimum_height = 0,
	alignment = 'top_right',
};
conky.text = [[
${time %H:%M:%S}${alignr}${time %d-%m-%y}
${voffset -16}${font sans-serif:bold:size=18}${alignc}${time %H:%M}${font}
${voffset 4}${alignc}${time %A %B %d, %Y}
${font}${voffset -4}
${font sans-serif:bold:size=10}SYSTEM ${hr 2}
${font sans-serif:normal:size=8}$sysname $kernel $alignr $machine
Host:$alignr$nodename
Uptime:$alignr$uptime
File System: $alignr${fs_type}
Processes: $alignr ${execi 1000 ps aux | wc -l}

${font sans-serif:bold:size=10}CPU ${hr 2}
${font sans-serif:normal:size=8}${execi 1000 grep model /proc/cpuinfo | cut -d : -f2 | tail -1 | sed 's/\s//'}
${font sans-serif:normal:size=8}${cpugraph cpu1}
CPU: ${cpu cpu1}% ${cpubar cpu1}

${font sans-serif:bold:size=10}MEMORY ${hr 2}
${font sans-serif:normal:size=8}RAM $alignc $mem / $memmax $alignr $memperc%
$membar
SWAP $alignc ${swap} / ${swapmax} $alignr ${swapperc}%
${swapbar}

${font sans-serif:bold:size=10}DISK USAGE ${hr 2}
${font sans-serif:normal:size=8}/ $alignc ${fs_used /} / ${fs_size /} $alignr ${fs_used_perc /}%
${fs_bar /}

${font Ubuntu:bold:size=10}NETWORK ${hr 2}
${font sans-serif:normal:size=8}Local IPs:${alignr}External IP:
${execi 1000 ip a | grep inet | grep -vw lo | grep -v inet6 | cut -d \/ -f1 | sed 's/[^0-9\.]*//g'}  ${alignr}${execi 1000  wget -q -O- http://ipecho.net/plain; echo}
${font sans-serif:normal:size=8}Down: ${downspeed enp0s3}  ${alignr}Up: ${upspeed enp0s3} 
${color lightgray}${downspeedgraph enp0s3 80,130 } ${alignr}${upspeedgraph enp0s3 80,130 }$color
${font sans-serif:bold:size=10}TOP PROCESSES ${hr 2}
${font sans-serif:normal:size=8}Name $alignr PID   CPU%   MEM%${font sans-serif:normal:size=8}
${top name 1} $alignr ${top pid 1} ${top cpu 1}% ${top mem 1}%
${top name 2} $alignr ${top pid 2} ${top cpu 2}% ${top mem 2}%
${top name 3} $alignr ${top pid 3} ${top cpu 3}% ${top mem 3}%
${top name 4} $alignr ${top pid 4} ${top cpu 4}% ${top mem 4}%
${top name 5} $alignr ${top pid 5} ${top cpu 5}% ${top mem 5}%
${top name 6} $alignr ${top pid 6} ${top cpu 6}% ${top mem 6}%
${top name 7} $alignr ${top pid 7} ${top cpu 7}% ${top mem 7}%
${top name 8} $alignr ${top pid 8} ${top cpu 8}% ${top mem 8}%
${top name 9} $alignr ${top pid 9} ${top cpu 9}% ${top mem 9}%
${top name 10} $alignr ${top pid 10} ${top cpu 10}% ${top mem 10}%
]];
```