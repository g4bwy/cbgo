~~~sh
go install ./cmd/cbmgr/
go install ./cmd/cbcustom/
~~~

Then add the following at the top of ~/.config/cagebreak/config:

```
exec ~/go/bin/cbmgr
```

You may also add key bindings:
```
bind w exec ~/go/bin/cbcustom list_views
bind W exec ~/go/bin/cbcustom list_ws_views
bind C-t exec ~/go/bin/cbcustom other_view
```

Then run cagebreak with IPC enabled:
~~~sh
cagebreak -e
~~~

To get complete window titles with ```list_views``` and ```list_ws_views``` commands, cagebreak
must also be run with ```--bs``` option.
