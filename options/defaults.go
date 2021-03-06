package options

// AddDefaultOptions adds internal options that can be set by the user through
// the command-line interface.
func (o *Options) AddDefaultOptions() {
	o.Add(NewBoolOption("center"))
	o.Add(NewStringOption("columns"))
	o.Add(NewStringOption("sort"))
	o.Add(NewStringOption("topbar"))
}

// Defaults is the default, internal configuration file.
const Defaults string = `
# Global options
set nocenter
set columns=artist,track,title,album,year,time
set sort=file,track,disc,album,year,albumartistsort
set topbar="|$shortname $version||;${tag|artist} \\- ${tag|title}||${tag|album}, ${tag|year};$volume $mode $elapsed ${state} $time;|[${list|index}/${list|total}] ${list|title}||;;"

# Song tag styles
style album teal
style artist yellow
style date green
style time darkmagenta
style title white bold
style track green
style year green

# Tracklist styles
style allTagsMissing red
style currentSong black yellow
style cursor black white
style header green bold
style mostTagsMissing red
style selection white blue

# Topbar styles
style elapsed green
style listIndex darkblue
style listTitle blue bold
style listTotal darkblue
style mute red
style shortName bold
style state default
style switches teal
style tagMissing red
style topbar darkgray
style version gray
style volume green

# Other styles
style commandText default
style errorText white red bold
style readout default
style searchText white bold
style sequenceText teal
style statusbar default
style visualText teal

# Keyboard bindings: cursor movement
bind <Up> cursor up
bind k cursor up
bind <Down> cursor down
bind j cursor down
bind <PgUp> cursor pgup
bind <PgDn> cursor pgdn
bind <Home> cursor home
bind gg cursor home
bind <End> cursor end
bind G cursor end
bind gc cursor current
bind R cursor random
bind b cursor prevOf album
bind e cursor nextOf album

# Keyboard bindings: input mode
bind : inputmode input
bind / inputmode search
bind <F3> inputmode search
bind v select visual
bind V select visual

# Keyboard bindings: player and mixer
bind <Enter> play selection
bind <Space> pause
bind s stop
bind h previous
bind l next
bind + volume +2
bind - volume -2
bind <left> seek -5
bind <right> seek +5
bind M volume mute

# Keyboard bindings: other
bind <C-c> quit
bind <C-l> redraw
bind <C-s> sort
bind i print file
bind t list next
bind T list previous
bind <C-d> list duplicate
bind <C-g> list remove
bind <C-j> isolate artist
bind <C-t> isolate albumartist album
bind & select nearby albumartist album
bind m select toggle
bind a add
bind <Delete> cut
bind x cut
bind y yank
bind p paste after
bind P paste before
`
